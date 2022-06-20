package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

func read() {
	defer wg.Done()

	f, err := os.Open("/project/article-manager/data-loader/dblp.xml")
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
