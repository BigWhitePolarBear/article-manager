package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	read2match = make(chan []byte, 16)
	readOK     = make(chan struct{})
)

func main() {
	go read()

	exist := map[string]struct{}{}
	builder := strings.Builder{}
	inTag := 0
	for {
		select {
		case <-readOK:
			os.Exit(0)
		case line := <-read2match:
			for i := range line {
				if line[i] == '<' {
					inTag++
					builder.WriteByte(line[i])
				} else if inTag > 0 {
					builder.WriteByte(line[i])
					if line[i] == '>' {
						inTag--
					}
				}
				if inTag == 0 && builder.Len() != 0 {
					str := builder.String()
					if len(str) > 1000 {
						builder.Reset()
						continue
					}
					_, ok := exist[str]
					if !ok {
						exist[str] = struct{}{}
						fmt.Println(str)
					}
					builder.Reset()
				}
			}
		}
	}
}

func read() {
	f, err := os.Open("../dblp.xml")
	if err != nil {
		panic(err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	reader := bufio.NewReader(f)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			readOK <- struct{}{}
			if err == io.EOF {
				readOK <- struct{}{}
				return
			}
			panic(err)
		}
		read2match <- line
	}
}
