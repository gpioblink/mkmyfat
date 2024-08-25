package models

import (
	"encoding/binary"
	"os"

	"gpioblink.com/app/makemyfat/mkmyfat/tools"
)

type FSInfo struct {
	FSI_LeadSig    uint32 // 0x41615252 // FSINFOのシグネチャ
	FSI_Reserved1  [480]byte
	FSI_StrucSig   uint32   // 0x61417272 // FSINFOのシグネチャ
	FSI_Free_Count uint32   // 空きクラスタ数 // 0xFFFFFFFFで無効
	FSI_Nxt_Free   uint32   // 次の空きクラスタ番号 // 0xFFFFFFFFで無効
	FSI_Reserved2  [12]byte // 予約領域。フォーマット時は0で埋める
	FSI_TrailSig   uint32   // 0xaa550000 // FSINFOの終了シグネチャ
}

func (fi *FSInfo) Export(bpb *Fat32BPB, f *os.File) error {
	// FSINFOの配置箇所を特定する
	_, err := f.Seek(int64(bpb.Sec2Addr(uint32(bpb.BPB_FSInfo))), 0)
	if err != nil {
		return err
	}

	err = binary.Write(f, binary.LittleEndian, fi)
	if err != nil {
		return err
	}

	// バックアップセクタへも書き込み
	_, err = f.Seek(int64(bpb.Sec2Addr(uint32(bpb.BPB_BkBootSec+bpb.BPB_FSInfo))), 0)
	if err != nil {
		return err
	}

	err = binary.Write(f, binary.LittleEndian, fi)
	if err != nil {
		return err
	}

	return nil
}

func (fi *FSInfo) String() string {
	return tools.PrettyPrintStruct("FSInfo", fi)
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
