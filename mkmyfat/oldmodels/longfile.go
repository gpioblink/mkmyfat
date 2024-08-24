package models

type LongFileName struct {
	LDIR_Ord       uint8    // エントリの順番
	LDIR_Name1     [10]uint // ファイル名の1-5文字目
	LDIR_Attr      uint8    // ファイル属性フラグ
	LDIR_Type      uint8    // エントリタイプ
	LDIR_Chksum    uint8    // 短いファイル名のチェックサム
	LDIR_Name2     [12]uint // ファイル名の6-11文字目
	LDIR_FstClusLO uint16   // ファーストクラスタの下位16ビット
	LDIR_Name3     [4]uint  // ファイル名の12-13文字目
}
