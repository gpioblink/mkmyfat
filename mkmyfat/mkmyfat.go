package mkmyfat

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"time"

	. "gpioblink.com/app/makemyfat/mkmyfat/models"
)

func PrintBPBFromFile(imgPath string) error {
	f, err := os.Open(imgPath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", imgPath, err)
	}
	defer f.Close()

	var fat32BPB Fat32BPB
	err = binary.Read(f, binary.LittleEndian, &fat32BPB)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %s", imgPath, err)
	}
	PrintBPB(fat32BPB)
	return nil
}

func Add(imgPath string, fileList []string) error {
	// imgPathのファイルを開く
	f, err := os.OpenFile(imgPath, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", imgPath, err)
	}
	defer f.Close()

	// fileListのファイルが存在することを確認する
	for _, file := range fileList {
		if !Exists(file) {
			return fmt.Errorf("file %s is not found", file)
		}
	}

	// FAT32のBPBを読み込む
	var bpb Fat32BPB
	err = binary.Read(f, binary.LittleEndian, &bpb)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %s", imgPath, err)
	}

	// imgPathのファイルがFAT32ファイルシステムであることを確認する
	if bpb.BS_jmpBoot != [3]byte{0xeb, 0x58, 0x90} {
		// TODO: これだけだと絶対条件不足なので、ちゃんと考える
		return fmt.Errorf("file %s is not FAT32 file system", imgPath)
	}

	// imgPathのファイルにfileListのファイルを書き込む
	for _, fileName := range fileList {

		// ファイルを開く
		src, err := os.Open(fileName)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %s", fileName, err)
		}
		defer src.Close()

		// ファイルのサイズを取得
		fi, err := src.Stat()
		if err != nil {
			return fmt.Errorf("failed to get file info %s: %s", fileName, err)
		}
		fileSize := int(fi.Size())

		// ファイルのクラスタ数を計算
		clusterSize := int(bpb.BPB_SecPerClus) * int(bpb.BPB_BytsPerSec)
		clusterCount := int(math.Ceil(float64(fileSize) / float64(clusterSize)))

		dirEntry := NewDirectoryEntry(GetShortName(fileName), ATTR_ARCHIVE, time.Now(), uint32(clusterCount), uint32(fileSize))

		// ルートディレクトリに書き込む
		_, err = f.Seek(int64(int(bpb.BPB_BytsPerSec)*int(bpb.BPB_RsvdSecCnt)), 0)
		if err != nil {
			return fmt.Errorf("failed to seek file %s: %s", imgPath, err)
		}

		err = binary.Write(f, binary.LittleEndian, dirEntry)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %s", imgPath, err)
		}

		// ファイルのデータを書き込む
		_, err = f.Seek(int64(int(bpb.BPB_BytsPerSec)*int(int(bpb.BPB_RsvdSecCnt)+int(bpb.BPB_NumFATs)*int(bpb.BPB_FATSz32))), 0)
		if err != nil {
			return fmt.Errorf("failed to seek file %s: %s", imgPath, err)
		}

		buf := make([]byte, clusterSize)
		for i := 0; i < clusterCount; i++ {
			_, err = src.Read(buf)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %s", fileName, err)
			}
			_, err = f.Write(buf)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %s", imgPath, err)
			}
		}

	}

	return nil
}
