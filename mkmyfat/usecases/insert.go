package usecases

import (
	"fmt"
	"os"

	"gpioblink.com/app/makemyfat/mkmyfat/models"
	"gpioblink.com/app/makemyfat/mkmyfat/tools"
)

func Insert(imgPath string, filePath string, entryNum int) error {
	f, err := os.OpenFile(imgPath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", imgPath, err)
	}
	defer f.Close()

	ff, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", imgPath, err)
	}
	defer ff.Close()

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

	// 置き換えるファイルの場所を特定
	addr, size, err := img.GetRootFileInfo(entryNum)
	if err != nil {
		return err
	}

	fmt.Printf("inserting data into entryNum: %d, addr: %d, size: %d\n", entryNum, addr, size)

	fileInfo, err := ff.Stat()
	if err != nil {
		return err
	}

	if fileInfo.Size() > int64(size) {
		return fmt.Errorf("file size is too large")
	}

	// クラスタごとに置き換え
	// FIXME: クラスタサイズを4096バイト(8セクタ)で固定している
	w := op.GetFAT32OffsetWriter()
	_, err = w.Seek(int64(addr), 0)
	if err != nil {
		return err
	}

	buf := make([]byte, 4096)
	for {
		n, err := ff.Read(buf)
		if err != nil {
			break
		}

		if n < 4096 {
			// 最終クラスタの場合は最後まで0で埋める
			for i := n; i < 4096; i++ {
				buf[i] = 0
			}
		}

		_, err = w.Write(buf)
		if err != nil {
			return err
		}
	}

	fmt.Println("finished!")

	return nil
}
