package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"github.com/lewiscowles1986/read-pinned-folders/utils"
	"github.com/richardlehane/mscfb"
)

func openFile(path string) (*mscfb.Reader, *os.File, error) {
	file, _ := os.Open(path)
	doc, err := mscfb.New(file)
	if err != nil {
		return nil, nil, err
	}
	return doc, file, nil
}

func parseAutomaticDestinationFile(path string) ([]string, error) {

	doc, file, err := openFile(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	for entry, err := doc.Next(); err == nil; entry, err = doc.Next() {
		if entry.Name == "DestList" {
			buf := make([]byte, entry.Size)
			len, err := doc.Read(buf)
			if err != nil {
				return nil, err
			}

			headerBytes := make([]byte, 32)
			copy(headerBytes, buf[0:32])
			dLHeader := destListHeader{}
			dLHeader.Version = binary.LittleEndian.Uint32(headerBytes[0:4])
			dLHeader.NumberOfEntries = binary.LittleEndian.Uint32(headerBytes[4:8])
			dLHeader.NumberOfPinnedEntries = binary.LittleEndian.Uint32(headerBytes[8:12])
			dLHeader.LastEntryNumber = binary.LittleEndian.Uint32(headerBytes[16:20])
			dLHeader.LastRevisionNumber = binary.LittleEndian.Uint32(headerBytes[24:28])

			paths := []string{}

			for dirIndex := 32; dirIndex < len; {
				pathSize := binary.LittleEndian.Uint16(buf[dirIndex+128 : dirIndex+130])
				entrySize := int(128 + 2 + pathSize*2 + 4)

				var dBytes = make([]byte, entrySize)
				copy(dBytes[:], buf[dirIndex:dirIndex+entrySize])

				pathLen := binary.LittleEndian.Uint16(dBytes[128:130]) * 2
				path, _ := utils.DecodeUTF16(dBytes[130 : 130+pathLen])

				// Currently skip built-in folders, like Desktop, Documents, etc...
				if !strings.HasPrefix(path, "knownfolder") {
					paths = append(paths, path)
				}
				dirIndex += int(entrySize)
			}
			return paths, nil
		}
	}
	return nil, fmt.Errorf("No DestList Found")
}
