package usecases

import (
	"fmt"
	"io"
	"os"

	"gpioblink.com/app/makemyfat/mkmyfat/tools"
)

func SaveVisualizeBinary(imgPath string, outPath string) error {
	f, err := os.Open(imgPath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", imgPath, err)
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %s", err)
	}

	fmt.Println("Start making visualize image...")
	img, err := tools.MakeFATVisualizeImage(io.NewSectionReader(f, 0, fileInfo.Size()))
	if err != nil {
		return fmt.Errorf("failed to make visualize image: %s", err)
	}

	fmt.Println("Saving visualize image...")
	err = tools.SaveImageAsJpeg(img, outPath)
	if err != nil {
		return fmt.Errorf("failed to save image: %s", err)
	}

	fmt.Println("Done!")

	return nil
}
