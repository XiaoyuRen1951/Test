package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
)
type MetricInfo struct {
	_Name string `json:"__name__"`
	Akii string `json:"app_kubernetes_io_instance"`
	Akimb string `json:"app_kubernetes_io_managed_by"`
	Akin string `json:"app_kubernetes_io_name"`
	Akipo string `json:"app_kubernetes_io_part_of"`
	Akiv string `json:"app_kubernetes_io_version"`
	ContainerName string `json:"container_name"`
	Gpu string `json:"gpu"`
	Hsc string `json:"helm_sh_chart"`
	Instance string `json:"instance"`
	Job string `json:"job"`
	KName string `json:"kubernetes_name"`
	KNamespace string `json:"kubernetes_namespace"`
	Name string `json:"name"`
	PName string `json:"pod_name"`
	PNamespace string `json:"pod_namespace"`
	Uuid string `json:"uuid"`
	Container string `json:"container"`
	Namespace string `json:"namespace"`
	Node string `json:"node"`
	Pod string `json:"pod"`
	Resource string `json:"resource"`
	Unit string `json:"unit"`
}

type ResultInfo struct {
	Metric MetricInfo `json:"metric"`
	RValue []interface{} `json:"values"`
}
type DataInfo struct {
	ResultType string `json:"resultType"`
	Result []ResultInfo `json:"result"`
}
type PrometheusInfo struct {
	Status string `json:"status"`
	Data DataInfo `json:"data"`
}

func main() {

	dcgmFile, err := os.Open("./0925-dcgm.log")
	defer dcgmFile.Close()
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	dcgmdec := json.NewDecoder(dcgmFile)
	//t,err := dec.Token()
	dcgm := new(PrometheusInfo)
	err = dcgmdec.Decode(&dcgm)
	if err != nil {
		fmt.Println(err)
		return
	}
	gpuFile, err := os.Open("./0925-pr.log")
	defer gpuFile.Close()
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	prdec := json.NewDecoder(gpuFile)

	pr := new(PrometheusInfo)
	err = prdec.Decode(&pr)
	if err != nil {
		fmt.Println(err)
		return
	}

	outf, err := os.Create("test.log")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer outf.Close()
	mp := make(map[string]int64)
	time := make([]string,0)
	//res := make(map[string]map[string]string)
	for _,v := range pr.Data.Result {
		if v.Metric.Resource  != "nvidia_com_gpu" {
			continue
		}
		//fmt.Println(v.RValue)
		//podinfo := res[v.Metric.Pod]
		value := reflect.ValueOf(v.RValue)
		for i:= 0; i < value.Len(); i++ {
			tmp := value.Index(i).Elem()
			timestamp := tmp.Index(0).Elem().Float()
			cntstr := tmp.Index(1).Elem().String()
			cnt, err := strconv.ParseInt(cntstr,10,64)
			if err != nil {
				fmt.Println(err)
				return
			}
			tmps := strconv.FormatInt(int64(timestamp), 10)
			if _,ok := mp[tmps]; ok {
				mp[tmps]+=cnt
			} else {
				mp[tmps]=cnt
				time = append(time,tmps)
			}
		}
	}
	sort.Strings(time)
	for _,idx := range time {
		//outf.WriteString(idx+" "+strconv.FormatInt(mp[idx],10)+"\n")
		outf.WriteString(strconv.FormatInt(mp[idx],10)+"\n")
	}
}

/*outf.WriteString(v.Metric.PName+" "+v.Metric.Gpu+" ")

		timeLayout := "2006-01-02 15:04:05"
		l := s.Index(0).Elem()
		ll := l.Index(0).Elem().Float()
		r := s.Index(s.Len()-1).Elem()
		rr := r.Index(0).Elem().Float()
		prel := int64(ll)
		prer := int64(rr)
		last := prer-prel
		datetime := time.Unix(prel, 0).Format(timeLayout)
		outf.WriteString(datetime+"\n")
		h := strconv.FormatInt(last/3600, 10)
		minute := strconv.FormatInt((last%3600)/60, 10)
		//sec := strconv.FormatInt(last%60,10)
		outf.WriteString( h +"h"+minute+"m\n")*/
