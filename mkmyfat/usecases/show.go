package usecases

import (
	"fmt"
	"os"

	"gpioblink.com/app/makemyfat/mkmyfat/models"
	"gpioblink.com/app/makemyfat/mkmyfat/tools"
)

func ShowImageInfo(imgPath string) error {
	f, err := os.Open(imgPath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", imgPath, err)
	}
	defer f.Close()

	isMBR := tools.IsMBR(f)

	var op tools.FAT32Operator

	if isMBR {
		op = tools.NewMBRFAT32Manager(f)
		mbr, err := models.ImportMBR(op.(tools.MBROperator).GetMBRSectionReader())
		if err != nil {
			return fmt.Errorf("failed to import MBR: %s", err)
		}
		fmt.Println(mbr)
	} else {
		op = tools.NewSimpleFAT32Manager(f)
	}

	img, err := models.ImportFAT32Image(op.GetFAT32SectionReader())
	if err != nil {
		return fmt.Errorf("failed to import image: %s", err)
	}

	fmt.Println(img)
	fmt.Println(img.PrintRootFileList())

	if isMBR {
		fmt.Println(img.GetRootFileListWithMBR(op.(tools.MBROperator).GetMBRSectionReader().Size()))
	}

	return nil
}
