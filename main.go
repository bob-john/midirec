package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"time"
)

func main() {
	var err error
	var i = flag.Int("i", 0, "midi input port (see `midicat ins`)")
	flag.Parse()
	cmd := exec.Command("midicat", "in", "-i="+strconv.Itoa(*i))
	out, err := cmd.StdoutPipe()
	check(err)
	check(cmd.Start())
	go func() {
		t1 := time.Now()
		scanner := bufio.NewScanner(out)
		for scanner.Scan() {
			t2 := time.Now()
			b, err := hex.DecodeString(scanner.Text())
			check(err)
			if len(b) == 0 || b[0]&0xF0 == 0xF0 {
				continue
			}
			d := t2.Sub(t1)
			t1 = t2
			fmt.Println(int64(d / time.Millisecond))
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	check(cmd.Process.Kill())
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
