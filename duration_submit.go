package main

import (
	"bufio"
	"flag"
	"time"

	"fmt"
	"io"
	"os"
)


func main() {
	var starttime string
	flag.StringVar(&starttime, "s", "2020-10-29 00:00:00", "start time")
	flag.Parse()

	File, err := os.Open("./dg-.log")
	defer File.Close()
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	rd := bufio.NewReader(File)
	cnt := 0
	start := ""
	end := ""
	var startUtc int64
	var endUtc int64
	timeLayout := "2006-01-02 15:04:05"
	podname:=""

	durationmp := make(map[int64] int)
	submitmp := make(map[int64] int)
	initt,_ := time.Parse(timeLayout, starttime)

	outo, err := os.Create("./result-onehour.log")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outo.Close()

	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		cnt++
		pos:=0
		//fmt.Println(cnt)
		for i:=0; i < len(line); i++ {
			if line[i]==' ' {
				podname=line[:i]
				pos=i+1
			} else {
				if line[i]=='Z' || line[i]=='T' {
					line=line[:i]+" "+line[i+1:]
				}
			}
		}
		if cnt&1==1 {
			start=line[pos:len(line)-2]
			times, _ := time.Parse(timeLayout, start)
			startUtc = times.Unix()

			submit := (startUtc-initt.Unix())/3600

			if _,ok := submitmp[submit];ok {
				submitmp[submit]++
			} else {
				submitmp[submit]=1
			}

		} else {
			end=line[pos:len(line)-2]
			times, _ := time.Parse(timeLayout, end)
			endUtc = times.Unix()

			durtaion :=(endUtc - startUtc)/3600

			if durtaion < 1 {
				outo.WriteString(podname+"\n")
			}

			if _,ok := durationmp[durtaion];ok {
				durationmp[durtaion]++
			} else {
				durationmp[durtaion]=1
			}

		}

	}
	outd, err := os.Create("./result-duration.log")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outd.Close()

	for i:=0;i<200;i++ {
		if _,ok := durationmp[int64(i)];ok {
			outd.WriteString(fmt.Sprintf("%d\n",  durationmp[int64(i)]))
		} else {
			outd.WriteString("0\n")
		}
	}

	outs, err := os.Create("./result-submit.log")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outs.Close()

	for i:=0;i<200;i++ {
		if _,ok := submitmp[int64(i)];ok {
			outs.WriteString(fmt.Sprintf("%d\n",  submitmp[int64(i)]))
		} else {
			outs.WriteString("0\n")
		}
	}
}
