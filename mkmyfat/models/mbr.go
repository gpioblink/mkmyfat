package models

import (
	"encoding/binary"
	"os"

	"gpioblink.com/app/makemyfat/mkmyfat/tools"
)

type Partition struct {
	PT_BootID    uint8  // ブートフラグ。0x00=ブート不可
	PT_StartHd   uint8  // 開始ヘッド LBAの場合は0
	PT_StartCySc uint16 // 開始セクタ LBAの場合は0
	PT_System    uint8  // ファイルシステムID 0x0c=FAT32(LBA)
	PT_EndHd     uint8  // 終了ヘッド LBAの場合は0
	PT_EndCySc   uint16 // 終了セクタ LBAの場合は0
	PT_LbaOfs    uint32 // [要変更]開始セクタのLBA
	PT_LbaSize   uint32 // [要変更]LBAでのサイズ
}

type MBR struct {
	MBR_BootCode   [446]byte // 未使用時は0x00
	MBR_Partition1 Partition
	MBR_Partition2 Partition
	MBR_Partition3 Partition
	MBR_Partition4 Partition
	MBR_Signature  uint16 // 0xaa55
}

func NewFAT32Partition(lbaOffset uint32, lbaSize uint32) *Partition {
	const PT_SYSTEM_FAT32_LBA = 0x0c
	return &Partition{
		PT_BootID:    0x00,
		PT_StartHd:   0,
		PT_StartCySc: 0,
		PT_System:    PT_SYSTEM_FAT32_LBA,
		PT_EndHd:     0,
		PT_EndCySc:   0,
		PT_LbaOfs:    lbaOffset,
		PT_LbaSize:   lbaSize,
	}
}

func (mbr *MBR) Export(f *os.File) error {
	_, err := f.Seek(0, 0)
	if err != nil {
		return err
	}

	err = binary.Write(f, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}

	return nil
}

func ImportMBR(f *os.File) (*MBR, error) {
	mbr := &MBR{}
	_, err := f.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	err = binary.Read(f, binary.LittleEndian, mbr)
	if err != nil {
		return nil, err
	}

	return mbr, nil
}

func (mbr *MBR) String() string {
	return tools.PrettyPrintStruct("MBR", mbr)
}

func NewFAT32MBR(storageSize uint64) *MBR {
	const MBR_SIGNATURE = 0xaa55
	const MBR_BYTES_PER_SEC = 512

	// とりあえずオフセットは決め打ち
	lbaOffset := uint32(0x000000800)

	lbaSize := uint32(storageSize / MBR_BYTES_PER_SEC)

	return &MBR{
		MBR_BootCode:   [446]byte{},
		MBR_Partition1: *NewFAT32Partition(lbaOffset, lbaSize),
		MBR_Partition2: Partition{},
		MBR_Partition3: Partition{},
		MBR_Partition4: Partition{},
		MBR_Signature:  MBR_SIGNATURE,
	}
}
