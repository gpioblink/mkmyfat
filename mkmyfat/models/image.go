package models

import (
	"fmt"
	"os"
	"time"
)

type FAT32Image struct {
	file     *os.File
	fat32BPB *Fat32BPB
	fsInfo   *FSInfo
	fat      *FAT
	rootClus *EntryCluster
}

func (f *FAT32Image) GetRootFileList() string {
	res := "***** Root File List *****\n"
	for _, v := range f.rootClus.cluster {
		if !v.IsLongName() {
			size := v.(*DirectoryEntry).DIR_FileSize
			clus := uint32(v.(*DirectoryEntry).DIR_FstClusLO) + uint32(v.(*DirectoryEntry).DIR_FstClusHI)<<16
			sec := f.fat32BPB.Clus2Sec(clus)
			addrS := f.fat32BPB.Sec2Addr(sec)
			addrE := uint32(addrS) + size
			res += fmt.Sprintf("%s[%dbytes]: %#x-%#x clus=%d\n", v.(*DirectoryEntry).DIR_Name, size, addrS, addrE, clus)
		}
	}
	return res
}

func (f *FAT32Image) AddEmptyFileToRoot(fileName string, fileSizeByte uint32) error {
	return f.rootClus.AddFileEntry(fileName, fileSizeByte, time.Now())
}

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

func ImportFAT32Image(f *os.File) (*FAT32Image, error) {
	bpb, err := ImportFAT32BPB(f)
	if err != nil {
		return nil, err
	}

	fsInfo, err := ImportFSInfo(bpb, f)
	if err != nil {
		return nil, err
	}

	fat, err := ImportFAT(bpb, f)
	if err != nil {
		return nil, err
	}

	rootClus, err := ImportRoot(bpb, fat, f)
	if err != nil {
		return nil, err
	}

	return &FAT32Image{
		file:     f,
		fat32BPB: bpb,
		fsInfo:   fsInfo,
		fat:      fat,
		rootClus: rootClus,
	}, nil
}

func (img *FAT32Image) String() string {
	return img.fat32BPB.String() + img.fsInfo.String() /*+ img.fat.String() + img.rootClus.String()*/
}

func NewFAT32Image(file *os.File, size uint64) *FAT32Image {
	bpb := NewFat32BPB(int(size))
	fsInfo := NewFSInfo()
	fat := NewFAT()
	rootClus := NewEntryCluster(bpb, fat)

	return &FAT32Image{
		file:     file,
		fat32BPB: bpb,
		fsInfo:   fsInfo,
		fat:      fat,
		rootClus: rootClus,
	}
}
