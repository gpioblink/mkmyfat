package tools

func Clus2Sec(cluster uint32, secPerClus uint8) uint32 {
	return uint32(cluster * uint32(secPerClus))
}

func Sec2Addr(sector uint32, sectorSize uint16) uint64 {
	return uint64(sector * uint32(sectorSize))
}

func FAT2Sec(reservedSectors uint16, fatTableSize uint32, fatIndex uint8) uint32 {
	return uint32(reservedSectors) + uint32(fatIndex)*uint32(fatTableSize)
}

func CalcClusterNum(fileSize uint32, secPerClus uint8, sectorSize uint16) uint32 {
	return fileSize / uint32(uint32(secPerClus)*uint32(sectorSize))
}
