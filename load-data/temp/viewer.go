package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	read2match = make(chan []byte, 16)
	readOK     = make(chan struct{})
)

func main() {
	go read()

	cnt := 0
	for {
		select {
		case line := <-read2match:
			cnt++
			if cnt > 35000000 {
				fmt.Println(string(line))
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
