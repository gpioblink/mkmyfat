package models

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand/v2"
	"os"

	"gpioblink.com/app/makemyfat/mkmyfat/tools"
)

type Fat32BPB struct {
	BS_jmpBoot     [3]byte // 0xeb 0x58 0x90 // short jmp x86 NOP
	BS_OEMName     [8]byte // "mkmy.fat"
	BPB_BytsPerSec uint16  // 512 // USBメモリの場合512
	BPB_SecPerClus uint8   // [決め打ちでいいかな]2 // アロケーションユニット当たりのセクタ数
	BPB_RsvdSecCnt uint16  // [決め打ちでいいかな]32 // ブートセクタの後に続く予約領域のセクタ数
	BPB_NumFATs    uint8   // 2 // FATの数(多重化)
	BPB_RootEntCnt uint16  // 0 // FAT32では無効値(0)
	BPB_TotSec16   uint16  // 0 // FAT32では無効値(0)
	BPB_Media      uint8   // 0xf8 // メディアタイプ(標準値0xf8)
	BPB_FATSz16    uint16  // 0 // FAT32では無効値(0)
	BPB_SecPerTrk  uint16  // 32 // ジオメトリのないストレージでは関係ない
	BPB_NumHeads   uint16  // 8 // ジオメトリのないストレージでは関係ない
	BPB_HiddSec    uint32  // 0
	BPB_TotSec32   uint32  // [要変更]ボリュームの総セクタ数

	BPB_FATSz32   uint32    // [要変更]FATあたりのセクタ数
	BPB_ExtFlags  uint16    // 0 // FATの冗長化に関するフラグ, ミラーリングとかしないなら0でよさそう
	BPB_FSVer     uint16    // 0 // ファイルシステムのバージョン。最新は0.0
	BPB_RootClus  uint32    // 2 // ルートディレクトリのクラスタ番号。大抵は先頭のクラスタ番号である2
	BPB_FSInfo    uint16    // 1 // FSINFO構造体の置かれるセクタ番号。常にブートセクタの次で1
	BPB_BkBootSec uint16    // 6 // ブートセクタのバックアップがおかれるセクタ番号。通常はブートセクタの次で6
	BPB_Reserved  [12]byte  // 0 // 予約領域。フォーマット時は0で埋める
	BS_DrvNum     uint8     // 0x80 // ブートドライブ番号。USBメモリの場合固定ディスクなので0x80
	BS_Reserved   uint8     // 0 // 予約領域。フォーマット時は0で埋める
	BS_BootSig    uint8     // 0x29 // ブートセクタの拡張シグネチャ。0x29であることが推奨されている
	BS_VolID      uint32    // [要変更]ボリュームのシリアル番号。時刻などから生成する
	BS_VolLab     [11]byte  // "NO NAME    " // ボリュームラベル。名前のない場合は!NO NAME"。終端も0x00
	BS_FilSysType [8]byte   // "FAT32   " // ファイルシステムの種類。FAT12, FAT16, FAT32のいずれか
	BS_BootCode32 [420]byte // 0 // ブートストラップコード。ブートは今回関係ないので0で埋める
	BS_Sign       uint16    // 0xaa55 // ブートセクタの終了シグネチャ。0xaa55であることが推奨されている
}

func (bpb *Fat32BPB) Export(f *os.File) error {
	// FAT32のBPBをファイルに書き込む
	_, err := f.Seek(0, 0)
	if err != nil {
		return err
	}

	err = binary.Write(f, binary.LittleEndian, bpb)
	if err != nil {
		return err
	}

	// バックアップセクタへも書き込み
	_, err = f.Seek(int64(bpb.Sec2Addr(uint32(bpb.BPB_BkBootSec))), 0)
	if err != nil {
		return err
	}

	err = binary.Write(f, binary.LittleEndian, bpb)
	if err != nil {
		return err
	}

	return nil
}

func (bpb *Fat32BPB) String() string {
	return tools.PrettyPrintStruct("BPB", bpb)
}

func ImportFAT32BPB(f *os.File) (*Fat32BPB, error) {
	var bpb Fat32BPB
	// FAT32のBPBを構造体に読み込む
	err := binary.Read(f, binary.LittleEndian, &bpb)
	if err != nil {
		return nil, err
	}

	// 読み込んだ構造体がFAT32ファイルシステムであることを確認する
	// TODO: これだけだと絶対条件不足なので、ちゃんと考える
	if bpb.BS_jmpBoot != [3]byte{0xeb, 0x58, 0x90} {
		return nil, fmt.Errorf("file %s is not FAT32 file system", f.Name())
	}

	// TODO: バックアップセクタを確認し、一致しなければエラーを返す

	return &bpb, nil
}

func NewFat32BPB(diskSize int) *Fat32BPB {
	const BOOT_BACKUP_SECTOR = 6
	const FAT32_ENTRY_SIZE_BYTE = 4

	bytesPerSector := 512
	sectorsPerCluster := 2
	reservedSectors := 32
	numberFATs := 2
	volId := rand.Uint32()

	// TODO: ビット演算で計算する
	totalSectors := diskSize / bytesPerSector
	totalClusters := totalSectors / sectorsPerCluster
	fatSectors := int(math.Ceil(float64(totalClusters*4) / float64(bytesPerSector)))
	dataSectors := totalSectors - (reservedSectors + (numberFATs * fatSectors))
	totalDataClusters := dataSectors / sectorsPerCluster
	fatTableSize := int(math.Ceil(float64(totalDataClusters*FAT32_ENTRY_SIZE_BYTE) / float64(bytesPerSector)))

	return &Fat32BPB{
		BS_jmpBoot:     [3]byte{0xeb, 0x58, 0x90},
		BS_OEMName:     [8]byte{'m', 'k', 'm', 'y', '.', 'f', 'a', 't'},
		BPB_BytsPerSec: uint16(bytesPerSector),
		BPB_SecPerClus: uint8(sectorsPerCluster),
		BPB_RsvdSecCnt: uint16(reservedSectors),
		BPB_NumFATs:    uint8(numberFATs),
		BPB_RootEntCnt: 0,
		BPB_TotSec16:   0,
		BPB_Media:      0xf8,
		BPB_FATSz16:    0,
		BPB_SecPerTrk:  32,
		BPB_NumHeads:   8,
		BPB_HiddSec:    0,
		BPB_TotSec32:   uint32(totalSectors),

		BPB_FATSz32:   uint32(fatTableSize),
		BPB_ExtFlags:  0,
		BPB_FSVer:     0,
		BPB_RootClus:  2,
		BPB_FSInfo:    1,
		BPB_BkBootSec: BOOT_BACKUP_SECTOR,
		BPB_Reserved:  [12]byte{},
		BS_DrvNum:     0x80,
		BS_Reserved:   0,
		BS_BootSig:    0x29,
		BS_VolID:      volId,
		BS_VolLab:     [11]byte{'N', 'O', ' ', 'N', 'A', 'M', 'E', ' ', ' ', ' ', ' '},
		BS_FilSysType: [8]byte{'F', 'A', 'T', '3', '2', ' ', ' ', ' '},
		BS_BootCode32: [420]byte{},
		BS_Sign:       0xaa55,
	}
}

const RESERVED_CLUSTERS = 2

func (bpb *Fat32BPB) RootSec() uint32 {
	return uint32(bpb.FAT2Sec(bpb.BPB_NumFATs))
}

func (bpb *Fat32BPB) UserSec() uint32 {
	// return bpb.RootSec() + uint32(bpb.BPB_RootEntCnt*32/bpb.BPB_BytsPerSec)
	// FAT32ではルートディレクトリとして特別なセクタは存在しない
	return bpb.RootSec()
}

func (bpb *Fat32BPB) Clus2Sec(cluster uint32) uint32 {
	return uint32(bpb.UserSec() + (cluster-RESERVED_CLUSTERS)*uint32(bpb.BPB_SecPerClus))
}

func (bpb *Fat32BPB) Sec2Addr(sector uint32) uint64 {
	return uint64(sector * uint32(bpb.BPB_BytsPerSec))
}

func (bpb *Fat32BPB) FAT2Sec(fatIndex uint8) uint32 {
	return uint32(bpb.BPB_RsvdSecCnt) + uint32(fatIndex)*uint32(bpb.BPB_FATSz32)
}

func (bpb *Fat32BPB) CalcClusterNum(fileSize uint32) uint32 {
	return fileSize / uint32(uint32(bpb.BPB_SecPerClus)*uint32(bpb.BPB_BytsPerSec))
}
