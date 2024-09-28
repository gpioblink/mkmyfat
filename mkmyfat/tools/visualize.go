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
	size := f.Size()
	width := 512
	height := int(size / int64(width))

	sec := 512
	tmpSec := make([]byte, sec)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for i := uint32(0); i < uint32(size/int64(sec)); i++ {
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
				img.Set(x, y+(int(addr)/width), color.RGBA{b, b, b, b})
			}
		}
	}
	return img, nil

}
