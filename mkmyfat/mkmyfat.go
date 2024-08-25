package mkmyfat

import (
	"encoding/binary"
	"fmt"
	"os"

	"gpioblink.com/app/makemyfat/mkmyfat/models"
	"gpioblink.com/app/makemyfat/mkmyfat/usecases"
)

func Create(imgPath string, diskSizeBytes int) error {
	err := usecases.Create(imgPath, diskSizeBytes)
	if err != nil {
		return err
	}
	return nil
}

func PrintBPBFromFile(imgPath string) error {
	f, err := os.Open(imgPath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", imgPath, err)
	}
	defer f.Close()

	var fat32BPB models.Fat32BPB
	err = binary.Read(f, binary.LittleEndian, &fat32BPB)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %s", imgPath, err)
	}
	PrintBPB(fat32BPB)
	return nil
}

func PrintBPB(bpb models.Fat32BPB) {
	fmt.Printf("BS_jmpBoot: %v\n", bpb.BS_jmpBoot)
	fmt.Printf("BS_OEMName: %v\n", bpb.BS_OEMName)
	fmt.Printf("BPB_BytsPerSec: %v\n", bpb.BPB_BytsPerSec)
	fmt.Printf("BPB_SecPerClus: %v\n", bpb.BPB_SecPerClus)
	fmt.Printf("BPB_RsvdSecCnt: %v\n", bpb.BPB_RsvdSecCnt)
	fmt.Printf("BPB_NumFATs: %v\n", bpb.BPB_NumFATs)
	fmt.Printf("BPB_RootEntCnt: %v\n", bpb.BPB_RootEntCnt)
	fmt.Printf("BPB_TotSec16: %v\n", bpb.BPB_TotSec16)
	fmt.Printf("BPB_Media: %v\n", bpb.BPB_Media)
	fmt.Printf("BPB_FATSz16: %v\n", bpb.BPB_FATSz16)
	fmt.Printf("BPB_SecPerTrk: %v\n", bpb.BPB_SecPerTrk)
	fmt.Printf("BPB_NumHeads: %v\n", bpb.BPB_NumHeads)
	fmt.Printf("BPB_HiddSec: %v\n", bpb.BPB_HiddSec)
	fmt.Printf("BPB_TotSec32: %v\n", bpb.BPB_TotSec32)
	fmt.Printf("BPB_FATSz32: %v\n", bpb.BPB_FATSz32)
	fmt.Printf("BPB_ExtFlags: %v\n", bpb.BPB_ExtFlags)
	fmt.Printf("BPB_FSVer: %v\n", bpb.BPB_FSVer)
	fmt.Printf("BPB_RootClus: %v\n", bpb.BPB_RootClus)
	fmt.Printf("BPB_FSInfo: %v\n", bpb.BPB_FSInfo)
	fmt.Printf("BPB_BkBootSec: %v\n", bpb.BPB_BkBootSec)
	fmt.Printf("BPB_Reserved: %v\n", bpb.BPB_Reserved)
	fmt.Printf("BS_DrvNum: %v\n", bpb.BS_DrvNum)
	fmt.Printf("BS_Reserved: %v\n", bpb.BS_Reserved)
	fmt.Printf("BS_BootSig: %v\n", bpb.BS_BootSig)
	fmt.Printf("BS_VolID: %v\n", bpb.BS_VolID)
	fmt.Printf("BS_VolLab: %v\n", bpb.BS_VolLab)
	fmt.Printf("BS_FilSysType: %v\n", bpb.BS_FilSysType)
	fmt.Printf("BS_BootCode32: %v\n", bpb.BS_BootCode32)
	fmt.Printf("BS_Sign: %v\n", bpb.BS_Sign)
}
