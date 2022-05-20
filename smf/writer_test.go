package smf

import (
	"bytes"
	"testing"
)

func TestVarint(t *testing.T) {
	var testcases = map[int][]byte{
		0x00000000: {0x0},
		0x00000040: {0x40},
		0x0000007F: {0x7F},
		0x00000080: {0x81, 0x0},
		0x00002000: {0xC0, 0x0},
		0x00003FFF: {0xFF, 0x7F},
		0x00004000: {0x81, 0x80, 0x00},
		0x00100000: {0xC0, 0x80, 0x00},
		0x001FFFFF: {0xFF, 0xFF, 0x7F},
		0x00200000: {0x81, 0x80, 0x80, 0x00},
		0x08000000: {0xC0, 0x80, 0x80, 0x00},
		0x0FFFFFFF: {0xFF, 0xFF, 0xFF, 0x7F},
	}
	for i, o := range testcases {
		if !bytes.Equal(Varint(i), o) {
			t.Errorf("Varint(%x) = %x; wants %x", i, Varint(i), o)
		}
	}
}
