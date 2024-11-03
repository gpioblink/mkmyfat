package tools

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"os"
)

func SaveImageAsPng(img image.Image, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("cannot create file: %s", err)
	}
	defer file.Close()

	png.Encode(file, img)
	return nil
}

func SaveImageAsJpeg(img image.Image, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("cannot create file: %s", err)
	}
	defer file.Close()

	err = jpeg.Encode(file, img, nil)
	if err != nil {
		return fmt.Errorf("cannot convert file: %s", err)
	}
	return nil
}

func MakeFATVisualizeImage(f *io.SectionReader) (image.Image, error) {
	//size := f.Size()
	width := 512
	//height := int(size / int64(width))
	height := 1500

	sec := 512
	tmpSec := make([]byte, sec)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	headLines := 200 // uint32(size/int64(sec))
	userLines := 600
	file2Lines := 600

	// head lines
	for i := uint32(0); i < uint32(headLines); i++ {
		// FATを1セクタ単位で読み込み
		addr := int64(sec) * int64(i)
		sectionReader := io.NewSectionReader(f, addr, int64(sec))
		err := binary.Read(sectionReader, binary.LittleEndian, &tmpSec)
		if err != nil {
			return nil, err
		}
		// i個ずつimgに格納
		// fmt.Println(addr)
		for y := 0; y < sec/width; y++ {
			for x := 0; x < width; x++ {
				b := tmpSec[x+y]
				img.Set(x, y+(int(addr)/width), color.RGBA{b, b, b, 255})
			}
		}
	}

	// user lines
	for i := uint32(0); i < uint32(userLines-200); i++ {
		// FATを1セクタ単位で読み込み
		addr := 0x0000000000204000 + int64(sec)*int64(i)
		sectionReader := io.NewSectionReader(f, addr, int64(sec))
		err := binary.Read(sectionReader, binary.LittleEndian, &tmpSec)
		if err != nil {
			return nil, err
		}
		// i個ずつimgに格納
		// fmt.Println(addr)
		for y := 0; y < sec/width; y++ {
			for x := 0; x < width; x++ {
				b := tmpSec[x+y]
				img.Set(x, y+(int(addr-0x0000000000204000)/width)+headLines, color.RGBA{b, b, b, 255})
			}
		}
	}

	// file2 lines
	for i := uint32(0); i < uint32(file2Lines-200); i++ {
		// FATを1セクタ単位で読み込み
		addr := 0x8204800 + int64(sec)*int64(i)
		sectionReader := io.NewSectionReader(f, addr, int64(sec))
		err := binary.Read(sectionReader, binary.LittleEndian, &tmpSec)
		if err != nil {
			return nil, err
		}
		// i個ずつimgに格納
		// fmt.Println(addr)
		for y := 0; y < sec/width; y++ {
			for x := 0; x < width; x++ {
				b := tmpSec[x+y]
				img.Set(x, y+(int(addr-0x8204800)/width)+headLines+userLines, color.RGBA{b, b, b, 255})
			}
		}
	}

	return img, nil

}
