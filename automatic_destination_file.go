package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"github.com/lewiscowles1986/read-pinned-folders/utils"
	"github.com/richardlehane/mscfb"
)

// Size of DestinationListHeader
const DestinationListHeaderSize = 32

// Offset start to path
const StartOfPathRelativeOffset = 130

// Offset to path entry
const PathEntryOffset = 128

func accountForWChar(pathLen uint16) uint16 {
	return pathLen * 2
}

func parseDestinationListHeader(buf []byte) destListHeader {
	headerBytes := make([]byte, DestinationListHeaderSize)
	copy(headerBytes, buf[0:DestinationListHeaderSize])
	dLHeader := destListHeader{}
	dLHeader.Version = binary.LittleEndian.Uint32(headerBytes[0:4])
	dLHeader.NumberOfEntries = binary.LittleEndian.Uint32(headerBytes[4:8])
	dLHeader.NumberOfPinnedEntries = binary.LittleEndian.Uint32(headerBytes[8:12])
	dLHeader.LastEntryNumber = binary.LittleEndian.Uint32(headerBytes[16:20])
	dLHeader.LastRevisionNumber = binary.LittleEndian.Uint32(headerBytes[24:28])

	return dLHeader
}

func parseAutomaticDestinationFile(doc *mscfb.Reader, file *os.File) ([]string, error) {
	for entry, err := doc.Next(); err == nil; entry, err = doc.Next() {
		if strings.ToLower(entry.Name) == "destlist" {
			buf := make([]byte, entry.Size)
			len, err := doc.Read(buf)
			if err != nil {
				return nil, err
			}

			parseDestinationListHeader(buf[0:DestinationListHeaderSize])

			paths := []string{}

			for dirIndex := DestinationListHeaderSize; dirIndex < len; {
				extentsOfDir := dirIndex + StartOfPathRelativeOffset
				pathSize := binary.LittleEndian.Uint16(
					buf[dirIndex+PathEntryOffset : extentsOfDir])
				entrySize := int(
					PathEntryOffset + 2 + accountForWChar(pathSize) + 4)

				var dBytes = make([]byte, entrySize)
				copy(dBytes[:], buf[dirIndex:dirIndex+entrySize])

				pathLen := accountForWChar(
					binary.LittleEndian.Uint16(
						dBytes[PathEntryOffset:StartOfPathRelativeOffset]))
				extentsOfPath := StartOfPathRelativeOffset + pathLen
				path, _ := utils.DecodeUTF16(
					dBytes[StartOfPathRelativeOffset:extentsOfPath])

				// Currently skip built-in folders, like Desktop, Documents, etc...
				if !strings.HasPrefix(strings.ToLower(path), "knownfolder") {
					paths = append(paths, path)
				}
				dirIndex += int(entrySize)
			}
			return paths, nil
		}
	}
	return nil, fmt.Errorf("no destlist found")
}
