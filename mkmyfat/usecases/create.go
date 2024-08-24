package usecases

import (
	"encoding/binary"
	"fmt"

	"gpioblink.com/app/makemyfat/mkmyfat/tools"
)

func Create(imgPath string, diskSizeBytes int) error {
	// 空のイメージファイルを新規作成
	f, err := tools.CreateSpecificatedSizeFileWhenNotExsisted(imgPath, uint64(diskSizeBytes))
	if err != nil {
		return err
	}

	fsi := NewFSInfo()

	// BOOTセクタ
	err = binary.Write(f, binary.LittleEndian, bpb)
	if err != nil {
		return fmt.Errorf("failed to write BOOT sector: %s", err)
	}

	// fsinfoセクタ
	err = binary.Write(f, binary.LittleEndian, fsi)
	if err != nil {
		return fmt.Errorf("failed to write BOOT sector: %s", err)
	}

	// バックアップブートセクタとして最初の3セクタをセクタ6にコピー
	_, err = f.Seek(int64(bytesPerSector*BOOT_BACKUP_SECTOR), 0)
	if err != nil {
		return fmt.Errorf("failed to seek file %s: %s", imgPath, err)
	}

	// BOOTセクタ
	err = binary.Write(f, binary.LittleEndian, bpb)
	if err != nil {
		return fmt.Errorf("failed to write BOOT sector: %s", err)
	}

	// fsinfoセクタ
	err = binary.Write(f, binary.LittleEndian, fsi)
	if err != nil {
		return fmt.Errorf("failed to write BOOT sector: %s", err)
	}

	/* f8 ff ff 0f ff ff ff 0f  f8 ff ff 0f を各FATの最初に入れる */
	// FAT1
	_, err = f.Seek(int64(bytesPerSector*resarvedSectors), 0)
	if err != nil {
		return fmt.Errorf("failed to seek file %s: %s", imgPath, err)
	}
	err = binary.Write(f, binary.LittleEndian, []byte{0xf8, 0xff, 0xff, 0x0f, 0xff, 0xff, 0xff, 0x0f, 0xf8, 0xff, 0xff, 0x0f})
	if err != nil {
		return fmt.Errorf("failed to write FAT1: %s", err)
	}

	// FAT2
	_, err = f.Seek(int64(bytesPerSector*(resarvedSectors+fatTableSize)), 0)
	if err != nil {
		return fmt.Errorf("failed to seek file %s: %s", imgPath, err)
	}
	err = binary.Write(f, binary.LittleEndian, []byte{0xf8, 0xff, 0xff, 0x0f, 0xff, 0xff, 0xff, 0x0f, 0xf8, 0xff, 0xff, 0x0f})
	if err != nil {
		return fmt.Errorf("failed to write FAT2: %s", err)
	}

	return nil
}
