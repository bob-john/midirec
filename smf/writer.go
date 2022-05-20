package smf

import (
	"encoding/binary"
	"fmt"
	"io"
)

func WriteHeader(w io.Writer, format, ntracks, division int16) (err error) {
	_, err = w.Write([]byte("MThd"))
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, int32(6))
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, format)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ntracks)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, division)
	if err != nil {
		return
	}
	return
}

func WriteTrack(w io.Writer, events []byte) (err error) {
	_, err = w.Write([]byte("MTrk"))
	if err != nil {
		return
	}
	fmt.Println(len(events), events)
	err = binary.Write(w, binary.BigEndian, int32(len(events)))
	if err != nil {
		return
	}
	_, err = w.Write(events)
	return
}

func WriteEvent(w io.Writer, dtime int, message []byte) (err error) {
	fmt.Printf("%x %x\n", Varint(dtime), message)
	_, err = w.Write(Varint(dtime))
	if err != nil {
		return
	}
	_, err = w.Write(message)
	return
}

func SMPTE(format, ticksPerFrame int8) int16 {
	return -int16(format)<<8 + int16(ticksPerFrame)
}

func Varint(i int) (b []byte) {
	var m byte
	for i != 0 {
		b = append([]byte{byte(i)&0x7F | m}, b...)
		i >>= 7
		m = 0x80
	}
	if len(b) == 0 {
		b = []byte{0}
	}
	return
}

var EOT = []byte{0xFF, 0x2F, 0x00}
