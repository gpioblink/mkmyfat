package models

import (
	"os"
)

type FAT32Image struct {
	file     *os.File
	fat32BPB *Fat32BPB
	fsInfo   *FSInfo
	fat      *FAT
	rootClus *DirectoryCluster
}

// func (f *FAT32Image) GetRootFileList() error {
// }

// func (f *FAT32Image) AddFile(filePath string) error {
// }

func (img *FAT32Image) Export() error {
	err := img.fat32BPB.Export(img.file)
	if err != nil {
		return err
	}

	err = img.fsInfo.Export(img.fat32BPB, img.file)
	if err != nil {
		return err
	}

	err = img.fat.Export(*img.fat32BPB, img.file)
	if err != nil {
		return err
	}

	err = img.rootClus.ExportRoot(img.fat32BPB, img.file)
	if err != nil {
		return err
	}

	return nil
}

func (img *FAT32Image) String() string {
	return img.fat32BPB.String() + img.fsInfo.String() + img.fat.String() + img.rootClus.String()
}

func NewFAT32Image(file *os.File, size uint64) *FAT32Image {
	bpb := NewFat32BPB(int(size))
	fsInfo := NewFSInfo()
	fat := NewFAT()
	rootClus := NewDirectoryCluster()

	return &FAT32Image{
		file:     file,
		fat32BPB: bpb,
		fsInfo:   fsInfo,
		fat:      fat,
		rootClus: rootClus,
	}
}
