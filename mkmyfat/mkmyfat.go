package mkmyfat

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"
)

const BOOT_BACKUP_SECTOR = 6

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

type FSInfo struct {
	FSI_LeadSig    uint32 // 0x41615252 // FSINFOのシグネチャ
	FSI_Reserved1  [480]byte
	FSI_StrucSig   uint32   // 0x61417272 // FSINFOのシグネチャ
	FSI_Free_Count uint32   // 空きクラスタ数 // 0xFFFFFFFFで無効
	FSI_Nxt_Free   uint32   // 次の空きクラスタ番号 // 0xFFFFFFFFで無効
	FSI_Reserved2  [12]byte // 予約領域。フォーマット時は0で埋める
	FSI_TrailSig   uint32   // 0xaa550000 // FSINFOの終了シグネチャ
}

type DirectoryEntry struct {
	DIR_Name         [11]byte // 短いファイル名
	DIR_Attr         uint8    // ファイル属性フラグ
	DIR_NTRes        uint8    // 小文字情報を記録するフラグ、不要なら0
	DIR_CrtTimeTenth uint8    // ファイル作成時の10ミリ秒単位
	DIR_CrtTime      uint16   // ファイル作成時刻(0-199)、2秒を200に分けたもの
	DIR_CrtDate      uint16   // ファイル作成日付、分解能は2秒
	DIR_LstAccDate   uint16   // 最終アクセス日付
	DIR_FstClusHI    uint16   // ファーストクラスタの上位16ビット
	DIR_WrtTime      uint16   // ファイル更新時刻（必須）
	DIR_WrtDate      uint16   // ファイル更新日付（必須）
	DIR_FstClusLO    uint16   // ファーストクラスタの下位16ビット。ファイルサイズがゼロの時はクラスタは割り当てられず常に0
	DIR_FileSize     uint32   // バイト単位のファイルサイズ。ディレクトリの場合は常に0
}

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

func PrintBPB(bpb Fat32BPB) {
	fmt.Printf("BS_jmpBoot: %v\n", bpb.BS_jmpBoot)
	fmt.Printf("BS_OEMName: %v\n", bpb.BS_OEMName)
	fmt.Printf("BPB_BytsPerSec: %v\n", bpb.BPB_BytsPerSec)
	fmt.Printf("BPB_SecPerClus: %v\n", bpb.BPB_SecPerClus)
	fmt.Printf("BPB_RsvdSecCnt: %v\n", bpb.BPB_RsvdSecCnt)
	fmt.Printf("BPB_NumFATs: %v\n", bpb.BPB_NumFATs)
	fmt.Printf("BPB_RootEntCnt: %v\n", bpb.BPB_RootEntCnt)
	fmt.Printf("BPB_TotSec16: %v\n", bpb.BPB_TotSec16)
	fmt.Printf("BPB_Media: %v\n", bpb.BPB_Media)
	fmt.Printf("BPB_FATSz16: %v\n", bpb.BPB_FATSz16)
	fmt.Printf("BPB_SecPerTrk: %v\n", bpb.BPB_SecPerTrk)
	fmt.Printf("BPB_NumHeads: %v\n", bpb.BPB_NumHeads)
	fmt.Printf("BPB_HiddSec: %v\n", bpb.BPB_HiddSec)
	fmt.Printf("BPB_TotSec32: %v\n", bpb.BPB_TotSec32)
	fmt.Printf("BPB_FATSz32: %v\n", bpb.BPB_FATSz32)
	fmt.Printf("BPB_ExtFlags: %v\n", bpb.BPB_ExtFlags)
	fmt.Printf("BPB_FSVer: %v\n", bpb.BPB_FSVer)
	fmt.Printf("BPB_RootClus: %v\n", bpb.BPB_RootClus)
	fmt.Printf("BPB_FSInfo: %v\n", bpb.BPB_FSInfo)
	fmt.Printf("BPB_BkBootSec: %v\n", bpb.BPB_BkBootSec)
	fmt.Printf("BPB_Reserved: %v\n", bpb.BPB_Reserved)
	fmt.Printf("BS_DrvNum: %v\n", bpb.BS_DrvNum)
	fmt.Printf("BS_Reserved: %v\n", bpb.BS_Reserved)
	fmt.Printf("BS_BootSig: %v\n", bpb.BS_BootSig)
	fmt.Printf("BS_VolID: %v\n", bpb.BS_VolID)
	fmt.Printf("BS_VolLab: %v\n", bpb.BS_VolLab)
	fmt.Printf("BS_FilSysType: %v\n", bpb.BS_FilSysType)
	fmt.Printf("BS_BootCode32: %v\n", bpb.BS_BootCode32)
	fmt.Printf("BS_Sign: %v\n", bpb.BS_Sign)
}

func Create(imgPath string, diskSizeBytes int) error {
	// imgPathのファイルが存在しないことを確認する
	if Exists(imgPath) {
		return fmt.Errorf("imgPath %s is exsisted. Please remove first", imgPath)
	}

	f, err := os.Create(imgPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %s", imgPath, err)
	}
	defer f.Close()

	// BOOTセクタの作成
	bytesPerSector := 512
	totalSectors := diskSizeBytes / 512
	sectorsPerCluster := 2
	resarvedSectors := 32
	numberFATs := 2

	totalClusters := totalSectors / sectorsPerCluster
	fatSectors := int(math.Ceil(float64(totalClusters*4) / float64(bytesPerSector)))
	dataSectors := totalSectors - (resarvedSectors + (numberFATs * fatSectors))
	totalDataClusters := dataSectors / sectorsPerCluster
	fatTableSize := int(math.Ceil(float64(totalDataClusters*4) / float64(bytesPerSector))) // FAT32のFATエントリサイズは4バイト(32ビット)

	volId := rand.Uint32()

	bpb := NewBPBSector(bytesPerSector, sectorsPerCluster, resarvedSectors, numberFATs, totalSectors, totalDataClusters, fatTableSize, volId)
	fsi := NewFSInfo()

	// 容量の確保
	err = os.Truncate(imgPath, int64(diskSizeBytes))
	if err != nil {
		return fmt.Errorf("failed to truncate file %s: %s", imgPath, err)
	}

	// BOOTセクタ
	err = binary.Write(f, binary.LittleEndian, bpb)
	if err != nil {
		return fmt.Errorf("failed to write BOOT sector: %s", err)
	}

	// fsinfoセクタ
	err = binary.Write(f, binary.LittleEndian, fsi)
	if err != nil {
		return fmt.Errorf("failed to write BOOT sector: %s", err)
	}

	// バックアップブートセクタとして最初の3セクタをセクタ6にコピー
	_, err = f.Seek(int64(bytesPerSector*BOOT_BACKUP_SECTOR), 0)
	if err != nil {
		return fmt.Errorf("failed to seek file %s: %s", imgPath, err)
	}

	// BOOTセクタ
	err = binary.Write(f, binary.LittleEndian, bpb)
	if err != nil {
		return fmt.Errorf("failed to write BOOT sector: %s", err)
	}

	// fsinfoセクタ
	err = binary.Write(f, binary.LittleEndian, fsi)
	if err != nil {
		return fmt.Errorf("failed to write BOOT sector: %s", err)
	}

	/* f8 ff ff 0f ff ff ff 0f  f8 ff ff 0f を各FATの最初に入れる */
	// FAT1
	_, err = f.Seek(int64(bytesPerSector*resarvedSectors), 0)
	if err != nil {
		return fmt.Errorf("failed to seek file %s: %s", imgPath, err)
	}
	err = binary.Write(f, binary.LittleEndian, []byte{0xf8, 0xff, 0xff, 0x0f, 0xff, 0xff, 0xff, 0x0f, 0xf8, 0xff, 0xff, 0x0f})
	if err != nil {
		return fmt.Errorf("failed to write FAT1: %s", err)
	}

	// FAT2
	_, err = f.Seek(int64(bytesPerSector*(resarvedSectors+fatSectors)), 0)
	if err != nil {
		return fmt.Errorf("failed to seek file %s: %s", imgPath, err)
	}
	err = binary.Write(f, binary.LittleEndian, []byte{0xf8, 0xff, 0xff, 0x0f, 0xff, 0xff, 0xff, 0x0f, 0xf8, 0xff, 0xff, 0x0f})
	if err != nil {
		return fmt.Errorf("failed to write FAT2: %s", err)
	}

	return nil
}

func NewDirectoryEntry() *DirectoryEntry {
	return &DirectoryEntry{
		DIR_Name:         [11]byte{},
		DIR_Attr:         0,
		DIR_NTRes:        0,
		DIR_CrtTimeTenth: 0,
		DIR_CrtTime:      0,
		DIR_CrtDate:      0,
		DIR_LstAccDate:   0,
		DIR_FstClusHI:    0,
		DIR_WrtTime:      0,
		DIR_WrtDate:      0,
		DIR_FstClusLO:    0,
		DIR_FileSize:     0,
	}
}

func NewFSInfo() *FSInfo {
	return &FSInfo{
		FSI_LeadSig:    0x41615252,
		FSI_Reserved1:  [480]byte{},
		FSI_StrucSig:   0x61417272,
		FSI_Free_Count: 0xffffffff,
		FSI_Nxt_Free:   0xffffffff,
		FSI_Reserved2:  [12]byte{},
		FSI_TrailSig:   0xaa550000,
	}
}

func NewBPBSector(bytesPerSector int, sectorsPerCluster int, reservedSectors int, numberFATs int, totalSectors int, totalClusters int, fatTableSize int, volId uint32) *Fat32BPB {
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

func Add(imgPath string, fileList []string) error {
	// imgPathのファイルが存在することを確認する
	if !Exists(imgPath) {
		return fmt.Errorf("imgPath %s is not exsisted", imgPath)
	}

	// fileListのファイルが存在することを確認する
	for _, file := range fileList {
		if !Exists(file) {
			return fmt.Errorf("file %s is not found", file)
		}
	}

	// imgPathのファイルがFAT32ファイルシステムであることを確認する

	// imgPathのファイルが書き込み可能であることを確認する
	// imgPathのファイルにfileListのファイルを書き込む

	return nil
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
