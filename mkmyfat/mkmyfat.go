package mkmyfat

import (
	"fmt"
	"os"
)

type fat32BPB struct {
	BS_jmpBoot     [3]byte // 0xeb 0x58 0x90 // short jmp x86 NOP
	BS_OEMName     [8]byte // "mkmy.fat"
	BPB_BytsPerSec uint16  // 512 // USBメモリの場合512
	BPB_SecPerClus uint8   // 2 // アロケーションユニット当たりのセクタ数
	BPB_RsvdSecCnt uint16  // 32 // ブートセクタの後に続く予約領域のセクタ数
	BPB_NumFATs    uint8   // 2 // FATの数(多重化)
	BPB_RootEntCnt uint16  // 0
	BPB_TotSec16   uint16  // 0
	BPB_Media      uint8   // 0xf8
	BPB_FATSz16    uint16  // 0
	BPB_SecPerTrk  uint16  // 32
	BPB_NumHeads   uint16  // 64
	BPB_HiddSec    uint32  // 0
	BPB_TotSec32   uint32  // 0

	BPB_FATSz32   uint32
	BPB_ExtFlags  uint16
	BPB_FSVer     uint16
	BPB_RootClus  uint32
	BPB_FSInfo    uint16
	BPB_BkBootSec uint16
	BPB_Reserved  [12]byte
	BS_DrvNum     uint8
	BS_Reserved   uint8
	BS_BootSig    uint8
	BS_VolID      uint32
	BS_VolLab     [11]byte
	BS_FilSysType [8]byte
	BS_BootCode32 [420]byte
	BS_BootSig32  uint16
	BS_Sign       uint32
}

func Create(imgPath string, fileList []string) error {
	// imgPathのファイルが存在しないことを確認する
	if Exists(imgPath) {
		return fmt.Errorf("imgPath %s is exsisted. Please remove first", imgPath)
	}

	// fileListのファイルが存在することを確認する
	for _, file := range fileList {
		if !Exists(file) {
			return fmt.Errorf("file %s is not found", file)
		}
	}

	// imgPathのファイルがFAT32ファイルシステムであることを確認する

	// imgPathのファイルが書き込み可能であることを確認する
	// imgPathのファイルにfileListのファイルを書き込む

	return nil
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
