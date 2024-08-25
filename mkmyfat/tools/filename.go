package tools

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf16"
)

const RANDOM_CHAR = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GetShortNameCheckSum(shortName [11]byte) uint8 {
	sum := uint8(0)
	for i := 0; i < 11; i++ {
		sum = (sum >> 1) + (sum << 7) + shortName[i]
	}
	return sum
}

func SplitLongFileNamePerEntry(filename string) ([][13]uint16, error) {
	var longName [][13]uint16

	if !CheckAsciiString(filename) {
		return nil, fmt.Errorf("filename must be ascii string")
	}

	words := utf16.Encode([]rune(filename))

	// 13文字ずつに分割
	for i := 0; i < len(words); i += 13 {
		var entry [13]uint16
		for j := 0; j < 13; j++ {
			if i+j < len(words) {
				entry[j] = words[i+j]
			} else {
				entry[j] = 0xFFFF
			}
		}
		longName = append(longName, entry)
	}

	return longName, nil
}

func GetRandomShortName(filename string) [11]byte {
	var shortName [11]byte
	for i := range shortName {
		shortName[i] = RANDOM_CHAR[rand.Intn(len(RANDOM_CHAR))]
	}
	return shortName
}

func GetShortNameFromLongName(longName string) ([11]byte, error) {
	shortName := [11]byte{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}

	if !CheckAsciiString(longName) {
		return shortName, fmt.Errorf("longName must be ascii string")
	}

	tname := strings.ToUpper(longName)

	ext := filepath.Ext(tname)
	if len(ext) < 3 {
		return shortName, fmt.Errorf("longName must have extension and more than 2 characters")
	}

	// 8.3形式のファイル名を取得

	for i := 0; i < 8; i++ {
		if i < len(tname) {
			if tname[i] == '.' {
				break
			}
			shortName[i] = tname[i]
		}
	}

	for i := 8; i < 11; i++ {
		if i-8 < len(ext) {
			shortName[i] = ext[i-7]
		}
	}

	return shortName, nil
}

func CheckAsciiString(s string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9\.\s]+$`).MatchString(s)
}
