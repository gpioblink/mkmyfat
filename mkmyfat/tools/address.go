package tools

func Sec2Addr(sector uint32, sectorSize uint16) uint64 {
	return uint64(sector * uint32(sectorSize))
}

func FAT2Sec(reservedSectors uint16, fatTableSize uint32, fatIndex uint8) uint32 {
	return uint32(reservedSectors) + uint32(fatIndex)*uint32(fatTableSize)
}
