package models

import (
	"os"
)

type FAT32Image struct {
	file  *os.File
	cache *FAT32ImageImportantSectorsCache
}

type FAT32ImageImportantSectorsCache struct {
	Fat32BPB *Fat32BPB
	FSInfo   *FSInfo
}

func (f *FAT32Image) AddFile(filePath string) error {
}

func NewFAT32Image(file *os.File, size uint64) (*FAT32Image, error) {
	tmp := &FAT32ImageImportantSectorsCache{
		Fat32BPB: NewFat32BPB(int(size)),
		FSInfo:   NewFSInfo(),
	}
}
