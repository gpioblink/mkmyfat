package models

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"

	"gpioblink.com/app/makemyfat/mkmyfat/tools"
)

type Entry interface {
	IsLongName() bool
}

type EntryCluster struct {
	cluster []Entry
	bpb     *Fat32BPB // FIXME: いったん定数値の取得用に追加。もっと適切な方法がありそう
	fat     *FAT
}

func (ec *EntryCluster) isShortNameDuplicated(shortName [11]byte) bool {
	for _, v := range ec.cluster {
		if !v.IsLongName() && v.(*DirectoryEntry).DIR_Name == shortName {
			return true
		}
	}
	return false
}

func (ec *EntryCluster) AddFileEntry(fileName string, fileSize uint32, lastModifiedDateTime time.Time) error {
	tmpEntries := []Entry{}

	if !tools.CheckAsciiString(fileName) {
		return fmt.Errorf("only ascii characters are allowed in the file name")
	}

	// shortNameが生成可能なファイル名か確認
	shortName, err := tools.GetShortNameFromLongName(fileName)
	if err != nil {
		return err
	}

	// shortNameが重複していないか確認
	if ec.isShortNameDuplicated(shortName) {
		return fmt.Errorf("short name is duplicated")
	}

	// clusterFromのためにFATから空きクラスタを取得
	clusterFrom, err := ec.fat.GetNextFreeCluster()
	if err != nil {
		return err
	}

	dirEntry := NewDirectoryEntry(shortName, ATTR_ARCHIVE, lastModifiedDateTime, clusterFrom, fileSize)
	tmpEntries = append(tmpEntries, dirEntry)

	// lfn用にファイル名を分割
	splitedFileName, err := tools.SplitLongFileNamePerEntry(fileName)
	if err != nil {
		return err
	}

	for i, v := range splitedFileName {
		// 13文字ごとにlfnの構築
		lfn := NewLongFileName(uint8(i), v, tools.GetShortNameCheckSum(shortName), clusterFrom)
		tmpEntries = append(tmpEntries, lfn)
	}

	// 作成したエントリの順番を変更 (TODO: mkfs.fatの仕様に合わせたが本当に必要かは不明)
	for i, j := 0, len(tmpEntries)-1; i < j; i, j = i+1, j-1 {
		tmpEntries[i], tmpEntries[j] = tmpEntries[j], tmpEntries[i]
	}

	// -- ここから破壊的な変更あり --

	// FATから取得した空きクラスタを使用済みにする
	err = ec.fat.MarkAsUsed(clusterFrom)
	if err != nil {
		return err
	}

	// ファイルサイズに必要なクラスタ分だけ連続したクラスタをFATに確保
	clusterNum := ec.bpb.CalcClusterNum(fileSize)
	err = ec.fat.AllocateContinuesSectors(clusterFrom, int(clusterNum))
	if err != nil {
		return err
	}

	ec.cluster = append(ec.cluster, tmpEntries...)
	return nil
}

func ImportRoot(bpb *Fat32BPB, fat *FAT, f *os.File) (*EntryCluster, error) {
	entries := []Entry{}
	entSize := 32
	tmpEnt := [32]byte{}

	// FIXME: とりあえず最初のクラスタしか見てないので、他のクラスタも見るようにする

	for i := uint32(0); i < uint32(bpb.BPB_BytsPerSec)*uint32(bpb.BPB_SecPerClus)/uint32(entSize); i++ {
		sectionReader := io.NewSectionReader(f, int64(bpb.Sec2Addr(bpb.UserSec()))+int64(i*uint32(entSize)), int64(entSize))
		err := binary.Read(sectionReader, binary.LittleEndian, &tmpEnt)
		if err != nil {
			return nil, err
		}
		reader := bytes.NewReader(tmpEnt[:])
		// tmpEntのタイプがlfnかどうかで処理を分ける
		if tmpEnt[11] == 0x0f {
			var lfn LongFileName
			err := binary.Read(reader, binary.LittleEndian, &lfn)
			if err != nil {
				return nil, err
			}
			entries = append(entries, &lfn)
		} else if tmpEnt[0] != 0x00 {
			var de DirectoryEntry
			err := binary.Read(reader, binary.LittleEndian, &de)
			if err != nil {
				return nil, err
			}
			entries = append(entries, &de)
		}
	}

	return &EntryCluster{entries, bpb, fat}, nil
}

func (ec *EntryCluster) ExportRoot(bpb *Fat32BPB, f *os.File) error {
	rootClusterAddr := ec.bpb.Sec2Addr(ec.bpb.Clus2Sec(bpb.BPB_RootClus))

	_, err := f.Seek(int64(rootClusterAddr), 0)
	if err != nil {
		return err
	}

	for _, v := range ec.cluster {
		err = binary.Write(f, binary.LittleEndian, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dc *EntryCluster) String() string {
	return tools.PrettyPrintStruct("DIR", dc.cluster)
}

func NewEntryCluster(bpb *Fat32BPB, fat *FAT) *EntryCluster {
	return &EntryCluster{
		cluster: make([]Entry, 0),
		bpb:     bpb,
		fat:     fat,
	}
}
