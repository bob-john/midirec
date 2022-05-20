package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"flag"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"time"

	"github.com/bob-john/midirec/smf"
)

func main() {
	var err error
	var (
		i = flag.Int("i", 0, "midi input port (see `midicat ins`)")
		// b = flag.Int("b", 0, "bpm")
	)
	flag.Parse()
	cmd := exec.Command("midicat", "in", "-i="+strconv.Itoa(*i))
	out, err := cmd.StdoutPipe()
	check(err)
	check(cmd.Start())
	var events bytes.Buffer
	go func() {
		t1 := time.Now()
		scanner := bufio.NewScanner(out)
		for scanner.Scan() {
			t2 := time.Now()
			b, err := hex.DecodeString(scanner.Text())
			check(err)
			// if len(b) == 1 && b[0] == 0xF8 {
			// 	tick++
			// }
			if len(b) == 0 || b[0]&0xF0 == 0xF0 {
				continue
			}
			d := t2.Sub(t1)
			t1 = t2
			// smf.WriteEvent(&events, int(d/time.Millisecond), b)
			smf.WriteEvent(&events, int(d*24/120/time.Minute), b)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	check(cmd.Process.Kill())

	smf.WriteEvent(&events, 0, smf.EOT)
	f, err := os.Create("out.mid")
	check(err)
	defer f.Close()
	check(smf.WriteHeader(f, 0, 1, smf.SMPTE(25, 40)))
	// check(smf.WriteHeader(f, 0, 1, 24))
	check(smf.WriteTrack(f, events.Bytes()))
	check(f.Close())
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
