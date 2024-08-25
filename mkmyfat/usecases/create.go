package usecases

import (
	"gpioblink.com/app/makemyfat/mkmyfat/models"
	"gpioblink.com/app/makemyfat/mkmyfat/tools"
)

func Create(imgPath string, diskSizeBytes int) error {
	// 空のイメージファイルを新規作成
	f, err := tools.CreateSpecificatedSizeFileWhenNotExsisted(imgPath, uint64(diskSizeBytes))
	if err != nil {
		return err
	}

	img := models.NewFAT32Image(f, uint64(diskSizeBytes))
	err = img.Export()
	if err != nil {
		return err
	}

	return nil
}
