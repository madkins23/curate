package hash8

import (
	"github.com/sigurn/crc8"
)

// Hash encapsulates the github.com/sigurn/crc8 package to make it simpler to use.
// Adding a Write() method makes it an io.Writer which allows io.Copy to be used.
// This is a simplification of the hash.Hash interface in the standard library.
type Hash struct {
	crc   uint8
	table *crc8.Table
}

// The specific CRC8 algorithm is specified by this table.
// CRC8_MAXIM was chosen more or less randomly from the various options.
var crc8Table = crc8.MakeTable(crc8.CRC8_MAXIM)

// New returns a new, initialized Hash object.
func New() *Hash {
	return &Hash{
		crc:   crc8.Init(crc8Table),
		table: crc8Table,
	}
}

// Sum8 returns the CRC8 hash for the data previously written to the Hash object.
func (h *Hash) Sum8() uint8 {
	return crc8.Complete(h.crc, h.table)
}

// Write a block of data to the Hash object.
func (h *Hash) Write(p []byte) (n int, err error) {
	h.crc = crc8.Update(h.crc, p, h.table)
	return len(p), nil
}
