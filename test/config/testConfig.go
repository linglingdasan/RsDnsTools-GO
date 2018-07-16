package main

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"reflect"
)

func main() {
	parseIps("configs/cdn.ips.json")
}


func parseIps(configfile string){
	f, err := os.Open(configfile)
	if err != nil {
		os.Exit(1)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil{
		os.Exit(1)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		fmt.Println("json file format error")
	}

	for name, value := range result{
		fmt.Printf("cdn name is %s\r\n", name)
		for _, cidr := range ToSlice(value){
			fmt.Println(cidr)
		}
	}

}

func ToSlice(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		panic("toslice arr not slice")
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}
