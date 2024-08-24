package tools

import "math/rand"

const RANDOM_CHAR = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GetShortName(filename string) [11]byte {
	var shortName [11]byte
	for i := range shortName {
		shortName[i] = RANDOM_CHAR[rand.Intn(len(RANDOM_CHAR))]
	}
	return shortName
}
