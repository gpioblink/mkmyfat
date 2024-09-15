package usecases

import (
	"fmt"

	"gpioblink.com/app/makemyfat/mkmyfat/models"
	"gpioblink.com/app/makemyfat/mkmyfat/tools"
)

func Create(imgPath string, diskSizeBytes int, fileExt string, numOfFiles int, eachFileSize int, withMBR bool) error {
	f, err := tools.CreateSpecificatedSizeFileWhenNotExsisted(imgPath, uint64(diskSizeBytes))
	if err != nil {
		return err
	}
	defer f.Close()

	var op tools.FAT32Operator
	var mbr *models.MBR

	if withMBR {
		op = tools.NewMBRFAT32Manager(f)
		mbr = models.NewFAT32MBR(uint64(diskSizeBytes))
	} else {
		op = tools.NewSimpleFAT32Manager(f)
	}

	img := models.NewFAT32Image(uint64(op.GetFAT32SectionReader().Size()))

	for i := 0; i < numOfFiles; i++ {
		// LFNのデバッグ用に長いファイル名を使用した
		err := img.AddEmptyFileToRoot(fmt.Sprintf("%d.%s", i, fileExt), uint32(eachFileSize))
		if err != nil {
			return err
		}
	}

	if withMBR {
		err = mbr.Export(op.(tools.MBROperator).GetMBROffsetWriter())
		if err != nil {
			return err
		}
		fmt.Println(mbr)
	}

	fmt.Println(img)

	err = img.Export(op.GetFAT32OffsetWriter())
	if err != nil {
		return err
	}

	fmt.Println(img.GetRootFileList())

	return nil
}
