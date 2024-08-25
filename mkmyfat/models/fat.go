package models

import (
	"encoding/binary"
	"os"

	"gpioblink.com/app/makemyfat/mkmyfat/tools"
)

type FAT map[uint32]uint32

func (fat *FAT) Export(bpb Fat32BPB, f *os.File) error {
	// FATの数だけ繰り返し
	for i := uint8(0); i < bpb.BPB_NumFATs; i++ {
		// FATの配置箇所を特定する
		fatStart := tools.FAT2Sec(bpb.BPB_RsvdSecCnt, bpb.BPB_FATSz32, i)
		_, err := f.Seek(int64(tools.Sec2Addr(fatStart, bpb.BPB_BytsPerSec)), 0)
		if err != nil {
			return err
		}

		// 値のある場所に書き込む
		for key, value := range *fat {
			_, err = f.Seek(int64(tools.Sec2Addr(fatStart, bpb.BPB_BytsPerSec))+int64(4*key), 0)
			if err != nil {
				return err
			}

			err = binary.Write(f, binary.LittleEndian, value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (fat *FAT) String() string {
	return tools.PrettyPrintStruct("FAT", fat)
}

func NewFAT() *FAT {
	return &FAT{
		0: 0x0ffffff8,
		1: 0x0fffffff,
		2: 0x0ffffff8,
	}
}
