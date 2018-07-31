package util

import (
	"encoding/json"
	"log"
	"os"
	"io/ioutil"
	"fmt"
	"reflect"
)

//一个view的名称，对应一组cidr设置
//匹配的时候，根据输入的一个IP，从配置中获取对应的view_id，如果没有匹配成功，则id为-1
//同理，根据输入的一个域名name，可以从配置中获取对应的name_id，
//最终的fwd是根据name_id，view_id联合确定的，因为fwd的数量不会很大，我们可以在这里用简单遍历的方式对fwd进行匹配
//（如果很大，那么数据结构上应该再进行调整）使用map[[2]int],int的方式来进行快速匹配

type cdn_ips struct {
	name2id	map[string]int
	cidr []dns_cidr
}

type Config struct{
	ServiceAddress	string
}

type Fwder struct{
	Name 	string
	Address	string
	Timeout	int
	Ecs 	string
	Acl 	[]string
	Domains	[]string

}
type ForwarderArray struct{
	Forwarder	[]*Fwder
}

func NewConfig(configFile string) *Config{
	config := parseConfigJson(configFile)

	return config
}

func parseConfigJson(path string) *Config {

	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Open config file failed: ", err)
		os.Exit(1)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("Read config file failed: ", err)
		os.Exit(1)
	}

	j := new(Config)
	err = json.Unmarshal(b, j)
	if err != nil {
		log.Fatal("Json syntex error: ", err)
		os.Exit(1)
	}

	return j
}

func ParseFwdArray(path string) *ForwarderArray{
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Open config file failed: ", err)
		os.Exit(1)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("Read config file failed: ", err)
		os.Exit(1)
	}

	fa := new(ForwarderArray)
	err = json.Unmarshal(b, fa)
	if err != nil{
		log.Fatal("Json syntex error: ", err)
		os.Exit(1)
	}

	return fa
}

func ParseConfigMapString(configfile string) map[string]string{
	f, err := os.Open(configfile)
	if err != nil {
		os.Exit(1)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil{
		os.Exit(1)
	}

	var result map[string]string

	if err := json.Unmarshal(b, &result); err != nil {
		fmt.Println("json file format error")
	}

	return result
}


func ParseConfigMapStringSlice(configfile string) map[string][]string{
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

	ret := make(map[string][]string)

	for name, value := range result{
		for _, ele := range toSlice(value){
			ret[name] = append(ret[name], fmt.Sprintf("%s", ele))
		}
	}

	return ret
}

func toSlice(arr interface{}) []interface{} {
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