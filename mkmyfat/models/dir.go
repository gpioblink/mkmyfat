package models

import (
	"encoding/binary"
	"os"
	"time"

	"gpioblink.com/app/makemyfat/mkmyfat/tools"
)

const ATTR_READ_ONLY = 0x01
const ATTR_HIDDEN = 0x02
const ATTR_SYSTEM = 0x04
const ATTR_VOLUME_ID = 0x08
const ATTR_DIRECTORY = 0x10
const ATTR_ARCHIVE = 0x20
const ATTR_LONG_NAME = 0x0f

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

func ImportDirectryEntry(f *os.File) (*DirectoryEntry, error) {
	bpb := &Fat32BPB{}
	// TODO:読み込んだ構造体がDirEntryの構造体であることを確認する

	// TODO: 別のfatを確認し、一致しなければエラーを返す

	return &bpb, nil
}

func (dir *DirectoryEntry) Export(f *os.File) error {
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

func NewDirectoryEntry(shortName [11]byte, attrFlag uint8, writeDateTime time.Time, clusterFrom uint32, fileSize uint32) *DirectoryEntry {
	date, time, _ := tools.GetDateTimeForFAT(writeDateTime)
	return &DirectoryEntry{
		DIR_Name:         shortName,
		DIR_Attr:         attrFlag,
		DIR_NTRes:        0,
		DIR_CrtTimeTenth: 0, // 作成日時を使わなければ常に0
		DIR_CrtTime:      0, // 作成日時を使わなければ常に0
		DIR_CrtDate:      0, // 作成日時を使わなければ常に0
		DIR_LstAccDate:   0, // オープン日時を使わなければ常に0
		DIR_FstClusHI:    uint16(clusterFrom >> 16),
		DIR_WrtTime:      time,
		DIR_WrtDate:      date,
		DIR_FstClusLO:    uint16(clusterFrom & 0xFFFF),
		DIR_FileSize:     fileSize,
	}
}
