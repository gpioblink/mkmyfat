package models

import (
	"encoding/binary"
	"fmt"
	"math"
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

func (fat *FAT) AllocateContinuesSectors(clusterFrom uint32, num int) error {
	// 連続したクラスタをFATに確保
	(*fat)[clusterFrom] = uint32(clusterFrom + 1)
	for i := uint32(1); i < uint32(num)-1; i++ {
		if (*fat)[clusterFrom+i] != 0x0 {
			return fmt.Errorf("cluster %d is already used", clusterFrom+i)
		}
		(*fat)[clusterFrom+i] = uint32(i + 1)
	}
	(*fat)[clusterFrom+uint32(num)-1] = 0x0fffffff
	return nil
}

func (fat *FAT) MarkAsUsed(cluster uint32) error {
	if (*fat)[cluster] != 0x0 {
		return fmt.Errorf("cluster %d is already used", cluster)
	}
	(*fat)[cluster] = 0x0fffffff
	return nil
}

func (fat *FAT) GetNextFreeCluster() (uint32, error) {
	var i uint32 = 0
	for i < math.MaxUint32 { // FIXME: 本来はFAT32の最大クラスタ数を取得するべき
		// キーが存在しなければ、それが次の使用可能なインデックス
		if _, exists := (*fat)[i]; !exists {
			return i, nil
		}
		i++
	}
	return 0, fmt.Errorf("no free cluster")
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
