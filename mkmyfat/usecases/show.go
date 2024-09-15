package usecases

import (
	"fmt"
	"os"

	"gpioblink.com/app/makemyfat/mkmyfat/models"
)

func ShowImageInfo(imgPath string) error {
	f, err := os.Open(imgPath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", imgPath, err)
	}
	defer f.Close()

	img, err := models.ImportFAT32Image(f)
	if err != nil {
		return fmt.Errorf("failed to import image: %s", err)
	}

	fmt.Println(img)
	fmt.Println(img.GetRootFileList())

	return nil
}
