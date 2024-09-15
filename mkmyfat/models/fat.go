package models

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"

	"gpioblink.com/app/makemyfat/mkmyfat/tools"
)

type FAT map[uint32]uint32

func (fat *FAT) Export(bpb Fat32BPB, f *os.File) error {

	// TODO: こんな1個ずつ読むのは遅いので、もっと効率的な方法を考える

	// FATの数だけ繰り返し
	for i := uint8(0); i < bpb.BPB_NumFATs; i++ {
		// FATの配置箇所を特定する
		fatStart := bpb.FAT2Sec(i)
		_, err := f.Seek(int64(bpb.Sec2Addr(fatStart)), 0)
		if err != nil {
			return err
		}

		// 値のある場所に書き込む
		for key, value := range *fat {
			_, err = f.Seek(int64(bpb.Sec2Addr(fatStart))+int64(4*key), 0)
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

func ImportFAT(bpb *Fat32BPB, f *os.File) (*FAT, error) {
	fat := make(FAT)
	fatPerSec := bpb.BPB_BytsPerSec / 4
	tmpSec := make([]uint32, fatPerSec)

	for i := uint32(0); i < bpb.BPB_FATSz32; i++ {
		// FATを1セクタ単位で読み込み
		sectionReader := io.NewSectionReader(f, int64(bpb.Sec2Addr(bpb.FAT2Sec(0)+i)), int64(bpb.BPB_BytsPerSec))
		err := binary.Read(sectionReader, binary.LittleEndian, &tmpSec)
		if err != nil {
			return nil, err
		}
		// i個ずつfatに格納
		for j := 0; j < len(tmpSec); j++ {
			if tmpSec[j] != 0 {
				fat[uint32(i*uint32(fatPerSec)+uint32(j))] = tmpSec[j]
			}
		}
	}

	// 別のFATがある場合、それぞれが一致しているか確認する
	for i := uint8(1); i < bpb.BPB_NumFATs; i++ {
		for j := uint32(0); j < bpb.BPB_FATSz32; j++ {
			// FATを1セクタ単位で読み込み
			sectionReader := io.NewSectionReader(f, int64(bpb.Sec2Addr(bpb.FAT2Sec(i)+j)), int64(bpb.BPB_BytsPerSec))
			err := binary.Read(sectionReader, binary.LittleEndian, &tmpSec)
			if err != nil {
				return nil, err
			}
			// i個ずつfatに格納
			for k := 0; k < len(tmpSec); k++ {
				if tmpSec[k] != 0 {
					if fat[uint32(j*uint32(fatPerSec)+uint32(k))] != tmpSec[k] {
						return nil, fmt.Errorf("FAT%d is not matched with FAT0", i)
					}
				}
			}
		}
	}

	return &fat, nil
}

func (fat *FAT) AllocateContinuesSectors(clusterFrom uint32, num int) error {
	// 連続したクラスタをFATに確保
	(*fat)[clusterFrom] = uint32(clusterFrom + 1)
	for i := uint32(1); i < uint32(num)-1; i++ {
		if (*fat)[clusterFrom+i] != 0x0 {
			return fmt.Errorf("cluster %d is already used", clusterFrom+i)
		}
		(*fat)[clusterFrom+i] = uint32(clusterFrom + i + 1)
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
