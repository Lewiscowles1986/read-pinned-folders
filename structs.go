package main

type destListHeader struct {
	Version               uint32
	NumberOfEntries       uint32
	NumberOfPinnedEntries uint32
	_                     [4]byte // Unknown Counter
	LastEntryNumber       uint32
	_                     [4]byte // Unknown 1
	LastRevisionNumber    uint32
	_                     [4]byte // Unknown 2
}
