package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
)

// ICO structire
type ICO struct {
	file string
	fh   *os.File
}

// PNG structure
type PNG struct {
	fileHandle *os.File
	height     uint8
	width      uint8
	depth      uint16 // bit/pixel
	size       uint32
	offset     uint32
	buffer     []byte
}

func OpenPng(file string) (*PNG, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer fd.Close()

	/*
		25byte PNG header - BigEndian
		00:	89 50 4e 47 0d 0a 1a 0a // 8byte - magic number
		IHDR chunk
		08:	xx xx xx xx // 4byte - chunk length
		12:	49 48 44 52 // 4byte - chunk type(IHDR)
		16:	xx xx xx xx // 4byte - width
		20:	xx xx xx xx // 4byte - height
		24:	xx          // 1byte - bit depth (bit/pixel)
	*/

	header := make([]byte, 25)
	_, err = fd.Read(header)
	if err != nil {
		return nil, err
	}

	// 8byte header[0:8] - magic number
	magic := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}
	if !bytes.Equal(magic[:], header[:8]) {
		return nil, errors.New("Not PNG file")
	}

	// 4byte header[12:16] - chunk type IHDR
	if !bytes.Equal([]byte("IHDR"), header[12:16]) {
		return nil, errors.New("PNG no IHDR chunk")
	}

	// 4byte header[16:20] - width
	width := binary.BigEndian.Uint32(header[16:20])

	// 4byte header[20:24] - height
	height := binary.BigEndian.Uint32(header[20:24])

	if width <= 256 && height <= 256 {
		// ICO format use 0 for 256px or larger
		if width >= 256 {
			width = 0
		}
		if height >= 256 {
			height = 0
		}
	}

	png := new(PNG)

	png.width = uint8(width)
	png.height = uint8(height)

	// 1byte header[25] - color depth
	png.depth = uint16(uint8(header[24]))

	stat, _ := os.Stat(file)
	png.size = uint32(stat.Size())

	_, err = fd.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	png.buffer = make([]byte, png.size)
	cnt, err := fd.Read(png.buffer)
	if err != nil {
		return nil, err
	}

	if cnt != int(png.size) {
		return nil, errors.New("Load PNG buffer failed")
	}

	return png, nil
}

func ICONDirNumber(num uint16) []byte {
	/*
		6byte ICONDIR - LittleEndian
		00:   00 00 // 2byte, must be 0
		02:   01 00 // 2byte, 1 for ICO
		04:   xx xx // 2byte, img number
	*/
	body := []byte{0, 0, 1, 0, 0, 0}
	binary.LittleEndian.PutUint16(body[4:6], num)
	return body
}

func ICONDirEntry(png *PNG) []byte {
	/*
		16byte ICONDIRENTRY - LittleEndian
		00:   xx    // 1byte, width
		01:   xx    // 1byte, height
		02:   00    // 1byte, color palette number, 0 for PNG
		03:   00    // 1byte, reserved, always 0
		04:   00 00 // 2byte, color planes, 0 for PNG
		06:   xx xx // 2byte, color depth
		08:   xx xx xx xx // 4byte, image size
		12:   xx xx xx xx // 4byte, image offset
	*/
	body := make([]byte, 16)
	copy(body[0:6], []byte{png.width, png.height, 0, 0, 0, 0})
	binary.LittleEndian.PutUint16(body[6:8], png.depth)
	binary.LittleEndian.PutUint32(body[8:12], png.size)
	binary.LittleEndian.PutUint32(body[12:16], png.offset)
	return body
}

func PNGToICON(fileList []string, fileOut string) error {
	pngList := []*PNG{}

	const (
		ICON_DIR_LEGNTH   uint32 = 6
		ICON_ENTRY_LEGNTH uint32 = 16
	)

	var totalSize uint32 = 0
	var entryLength uint32 = ICON_ENTRY_LEGNTH * uint32(len(fileList))

	for _, file := range fileList {
		png, err := OpenPng(file)
		if err != nil {
			return err
		}
		png.offset = ICON_DIR_LEGNTH + entryLength + totalSize
		pngList = append(pngList, png)

		totalSize += png.size
	}

	// remove if output file exist.
	os.Remove(fileOut)

	output, err := os.Create(fileOut)
	if err != nil {
		return err
	}

	defer output.Close()

	_, err = output.Write(ICONDirNumber(uint16(len(fileList))))
	if err != nil {
		return err
	}

	// Copy PNG info
	for _, png := range pngList {
		_, err = output.Write(ICONDirEntry(png))
		if err != nil {
			return err
		}
	}

	for _, png := range pngList {
		_, err = output.Write(png.buffer)
		if err != nil {
			return err
		}
	}

	return nil
}
