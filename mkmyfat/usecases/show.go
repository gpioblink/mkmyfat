package usecases

import (
	"fmt"
	"io"
	"os"

	"gpioblink.com/app/makemyfat/mkmyfat/models"
)

func ShowImageInfo(imgPath string) error {
	f, err := os.Open(imgPath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", imgPath, err)
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %s", err)
	}

	img, err := models.ImportFAT32Image(io.NewSectionReader(f, 0, fileInfo.Size()))
	if err != nil {
		return fmt.Errorf("failed to import image: %s", err)
	}

	fmt.Println(img)
	fmt.Println(img.GetRootFileList())

	return nil
}
