package models

import (
	"encoding/binary"
	"fmt"
	"os"

	"gpioblink.com/app/makemyfat/mkmyfat/tools"
)

type FAT map[uint32]uint32

func ImportFAT(f *os.File) (*FAT, error) {
	var bpb Fat32BPB
	fat := make(FAT)

	// FAT32のBPBを構造体に読み込む
	err := binary.Read(f, binary.LittleEndian, &bpb)
	if err != nil {
		return nil, err
	}

	// FATの配置箇所を特定する
	fatStart := tools.FAT2Sec(bpb.BPB_RsvdSecCnt, bpb.BPB_FATSz32, 0)
	_, err = f.Seek(int64(tools.Sec2Addr(fatStart, bpb.BPB_BytsPerSec)), 0)
	if err != nil {
		return nil, err
	}

	// FATを読み込む。0以外になっているデータを読み込む
	for i := uint32(0); i < bpb.BPB_FATSz32; i++ {
		var value uint32
		err = binary.Read(f, binary.LittleEndian, &value)
		if err != nil {
			return nil, err
		}

		if value != 0 {
			fat[i] = value
		}
	}

	// TODO:読み込んだ構造体がFATの構造体であることを確認する
	// TODO: これだけだと絶対条件不足なので、ちゃんと考える
	if fat[0] != 0x0ffffff8 || fat[1] != 0x0fffffff || fat[2] != 0x0ffffff8 {
		return nil, fmt.Errorf("file %s is not FAT", f.Name())
	}

	// TODO: 別のfatを確認し、一致しなければエラーを返す

	return &fat, nil
}

func (fat *FAT) Export(f *os.File) error {
	// BPBを読み込み、配置箇所を特定する
	bpb, err := ImportFAT32BPB(f)
	if err != nil {
		return err
	}

	for i := uint8(0); i < bpb.BPB_NumFATs; i++ {
		// FATの配置箇所を特定する
		fatStart := tools.FAT2Sec(bpb.BPB_RsvdSecCnt, bpb.BPB_FATSz32, i)
		_, err = f.Seek(int64(tools.Sec2Addr(fatStart, bpb.BPB_BytsPerSec)), 0)
		if err != nil {
			return err
		}

		// テーブルをすべて0で埋める
		_, err = f.Write(make([]byte, tools.Sec2Addr(bpb.BPB_FATSz32, bpb.BPB_BytsPerSec)))
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

func NewFAT() FAT {
	return FAT{
		0: 0x0ffffff8,
		1: 0x0fffffff,
		2: 0x0ffffff8,
	}
}
