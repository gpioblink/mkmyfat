package tools

import (
	"io"
	"os"
)

type FAT32Operator interface {
	GetFAT32SectionReader() *io.SectionReader
	GetFAT32OffsetWriter() *io.OffsetWriter
}

type MBROperator interface {
	GetMBRSectionReader() *io.SectionReader
	GetMBROffsetWriter() *io.OffsetWriter
}

type SimpleFAT32Manager struct {
	f *os.File
	FAT32Operator
}

func (m *SimpleFAT32Manager) GetFAT32SectionReader() *io.SectionReader {
	fileInfo, _ := m.f.Stat()
	return io.NewSectionReader(m.f, 0, fileInfo.Size())
}

func (m *SimpleFAT32Manager) GetFAT32OffsetWriter() *io.OffsetWriter {
	return io.NewOffsetWriter(m.f, 0)
}

func NewSimpleFAT32Manager(f *os.File) *SimpleFAT32Manager {
	return &SimpleFAT32Manager{f: f}
}

type MBRFAT32Manager struct {
	f *os.File
	MBROperator
	FAT32Operator
}

func (m *MBRFAT32Manager) GetMBRSectionReader() *io.SectionReader {
	// FIXME: セクタあたりバイト数(512)とセクタ数(2048)をハードコーディングしている
	return io.NewSectionReader(m.f, 0, 512*2048)
}

func (m *MBRFAT32Manager) GetMBROffsetWriter() *io.OffsetWriter {
	return io.NewOffsetWriter(m.f, 0)
}

func (m *MBRFAT32Manager) GetFAT32SectionReader() *io.SectionReader {
	fileInfo, _ := m.f.Stat()
	// FIXME: セクタあたりバイト数(512)とセクタ数(2048)をハードコーディングしている
	return io.NewSectionReader(m.f, 512*2048, fileInfo.Size()-512*2048)
}

func (m *MBRFAT32Manager) GetFAT32OffsetWriter() *io.OffsetWriter {
	return io.NewOffsetWriter(m.f, 512*2048)
}

func NewMBRFAT32Manager(f *os.File) *MBRFAT32Manager {
	return &MBRFAT32Manager{f: f}
}

func IsMBR(f *os.File) bool {
	// FIXME: これだけだと明らかにMBRかどうかの判定が不十分
	firstSector := make([]byte, 512)
	_, err := f.Read(firstSector)
	if err != nil {
		return false
	}
	return firstSector[0x1c2] == 0x0c
}
