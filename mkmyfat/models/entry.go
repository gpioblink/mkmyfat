package models

import (
	"encoding/binary"
	"os"

	"gpioblink.com/app/makemyfat/mkmyfat/tools"
)

type Entry interface {
	IsLongName() bool
}

type EntryCluster struct {
	cluster []Entry
}

func (dc *EntryCluster) ExportRoot(bpb *Fat32BPB, f *os.File) error {
	rootClusterAddr := tools.Sec2Addr(tools.Clus2Sec(bpb.BPB_RootClus, bpb.BPB_SecPerClus), bpb.BPB_BytsPerSec)

	// rootClusterの配置箇所に移動する
	_, err := f.Seek(int64(rootClusterAddr), 0)
	if err != nil {
		return err
	}

	for _, v := range dc.cluster {
		err = binary.Write(f, binary.LittleEndian, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dc *EntryCluster) String() string {
	return tools.PrettyPrintStruct("DIR", dc)
}

func NewEntryCluster() *EntryCluster {
	return &EntryCluster{
		cluster: make([]Entry, 0),
	}
}
