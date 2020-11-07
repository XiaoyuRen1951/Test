package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func main() {
	var filename string
	flag.StringVar(&filename, "n", "null", "File name")
	flag.Parse()

	if filename == "null" {
		fmt.Println("Open fail!")
		return
	}

	File, err := os.Open(filename)
	defer File.Close()
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	rd := bufio.NewReader(File)
	timeLayout := "2006-01-02 15:04:05"

	memmp := make(map[int64]int64)

	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		pos := 0
		for i := 0; i < len(line); i++ {
			if line[i] == ' ' {
				pos = i
			} else {
				if line[i] == 'Z' || line[i] == 'T' {
					line = line[:i] + " " + line[i+1:]
				}
			}
		}
		tim := line[:pos-1]
		mem := line[pos+1:]

		times, _ := time.Parse(timeLayout, tim)

		memuti, err := strconv.ParseInt(mem[:len(mem)-3], 10, 64)

		//fmt.Println(times.Unix())
		//fmt.Println(memuti)

		if mem[len(mem)-3] == 'K' {
			memuti = memuti / 1024
		}

		if _, ok := memmp[times.Unix()]; ok {
			memmp[times.Unix()] += memuti
		} else {
			_, ok1 := memmp[times.Unix()-1]

			if ok1 {
				memmp[times.Unix()-1] += memuti
			} else {
				memmp[times.Unix()] = memuti
			}
		}
	}

	outm, err := os.Create("./result-mem.log")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outm.Close()

	for k, v := range memmp {
		outm.WriteString(fmt.Sprintf("%d %d\n", k, v))
	}
}
