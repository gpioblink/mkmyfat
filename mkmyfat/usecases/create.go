package usecases

import (
	"fmt"

	"gpioblink.com/app/makemyfat/mkmyfat/models"
	"gpioblink.com/app/makemyfat/mkmyfat/tools"
)

func Create(imgPath string, diskSizeBytes int) error {
	f, err := tools.CreateSpecificatedSizeFileWhenNotExsisted(imgPath, uint64(diskSizeBytes))
	if err != nil {
		return err
	}
	defer f.Close()

	img := models.NewFAT32Image(f, uint64(diskSizeBytes))
	err = img.Export()
	if err != nil {
		return err
	}

	fmt.Println(img)

	return nil
}

func CreateWithEmptyFiles(imgPath string, diskSizeBytes int, fileExt string, numOfFiles int, eachFileSize int) error {
	f, err := tools.CreateSpecificatedSizeFileWhenNotExsisted(imgPath, uint64(diskSizeBytes))
	if err != nil {
		return err
	}
	defer f.Close()

	img := models.NewFAT32Image(f, uint64(diskSizeBytes))

	fmt.Println(img)

	for i := 0; i < numOfFiles; i++ {
		// LFNのデバッグ用に長いファイル名を使用した
		err := img.AddEmptyFileToRoot(fmt.Sprintf("%d.%s", i, fileExt), uint32(eachFileSize))
		if err != nil {
			return err
		}
	}

	err = img.Export()
	if err != nil {
		return err
	}

	fmt.Println(img.GetRootFileList())

	return nil
}
