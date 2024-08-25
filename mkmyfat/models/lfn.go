package models

type LongFileName struct {
	LDIR_Ord       uint8     // エントリの順番
	LDIR_Name1     [5]uint16 // ファイル名の1-5文字目
	LDIR_Attr      uint8     // ファイル属性フラグ
	LDIR_Type      uint8     // エントリタイプ
	LDIR_Chksum    uint8     // 短いファイル名のチェックサム
	LDIR_Name2     [6]uint16 // ファイル名の6-11文字目
	LDIR_FstClusLO uint16    // ファーストクラスタの下位16ビット
	LDIR_Name3     [2]uint16 // ファイル名の12-13文字目
}

func (de *LongFileName) IsLongName() bool {
	return (de.LDIR_Attr & ATTR_LONG_NAME) != 0
}

func NewLongFileName(order uint8, name [13]uint16, chksum uint8, clusterFrom uint32) *LongFileName {
	return &LongFileName{
		LDIR_Ord:       order,
		LDIR_Name1:     [5]uint16{name[0], name[1], name[2], name[3], name[4]},
		LDIR_Attr:      ATTR_LONG_NAME,
		LDIR_Type:      0,
		LDIR_Chksum:    chksum,
		LDIR_Name2:     [6]uint16{name[5], name[6], name[7], name[8], name[9], name[10]},
		LDIR_FstClusLO: uint16(clusterFrom & 0xFFFF),
		LDIR_Name3:     [2]uint16{name[11], name[12]},
	}
}
