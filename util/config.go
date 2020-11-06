package util

import (
	"RsDnsTools/pool"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
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
	ServiceAddress	string `json:"ServerAddress"`
	IsAdaptiveEcs	bool	`json:"isAdaptiveEcs"`
	IpConfigFile	string	`json:"IpConfigFile"`
	DomainConfigFile	string	`json:DomainConfigFile`
	FwdConfigFile	string	`json:FwdConfigFile`
	EcsMapConfigFile	string	`json:EcsMapConfigFile`
	AdaptiveFwdConfigFile 	string 	`json:"AdaptiveFwdConfigFile"`
	//acl name ==> acl id
	AclNameIdMap	map[string]int
	//acl id ==> acl name
	AclIdNameMap	map[int]string
	//domain list name ==> domain list id
	DLNameIdMap		map[string]int
	//domain list id ==> domain list name
	DLIdNameMap		map[int]string
	//acl id ==> ecs ip
	AclEcsMap	map[int]string

	AclMap		*dns_cidr
	DnameList	*dnamelist
	Forwarders  *ForwarderArray
	AdaptiveForwarders  *AdaptiveForwarderArray
}

type Fwder struct{
	Name 	string	`json:Name`
	Address	string	`json:Address`
	Timeout	int
	Ecs 	bool
	Acl 	[]string
	Domains	[]string
	Default bool	`json:Default`

	FwdPool pool.Pool

}

type HitPoint struct{
	AclId	int
	NameId	int
}

type ForwarderArray struct{
	Forwarder	[]*Fwder `json:Forwarder`
	FwdMap		map[HitPoint]*Fwder

}

type AdapiveFwder struct{
	Name 	string	`json:Name`
	Address	string	`json:Address`
	Timeout	int
	Ecs 	bool
	Nemask 	int		`json:"nemask"`
	Domains	[]string
	Default bool	`json:Default`
	SetRespEcsIP bool `json:"SetRespEcsIp"`
	FwdPool pool.Pool
}
type AdaptiveForwarderArray struct{
	Forwarder	[]*AdapiveFwder `json:Forwarder`
	DnameTree		*dnametree
}

func NewConfig(configFile string) *Config{
	config := parseConfigJson(configFile)
	if config.IsAdaptiveEcs == true{
		getAdaptiveConfig(config)
	}else{
		getConfig(config)
	}
	return config
}

func getAdaptiveConfig(config *Config) *Config{
	//domain name部分
	dnmap := ParseConfigMapStringSlice(config.DomainConfigFile)
	config.DnameList = NewDnameList()
	cdnNames := make(map[string][]string)
	for dgname, dn := range dnmap{
		cdnNames[dgname] = append(cdnNames[dgname], "1")
		for _, dname := range dn{
			if !(dname[len(dname)-1:] == ".") {
				dname += "."
				cdnNames[dgname] = append(cdnNames[dgname], dname)
			}
			log.Info(dname)
		}
		cdnNames[dgname] = cdnNames[dgname][1:]
	}

	//fwd部分
	tfa := ParseAdaptiveFwdArray(config.AdaptiveFwdConfigFile)
	tfa.DnameTree = NewDnameTree()
	for i:=0; i<len(tfa.Forwarder);i++ {
		log.Debugf("%+v\n", tfa.Forwarder[i])
		for _, dnSetName := range tfa.Forwarder[i].Domains {
			for _, dname := range cdnNames[dnSetName] {
				tfa.DnameTree.Insert(dname, tfa.Forwarder[i])
			}
		}
	}
	config.AdaptiveForwarders = tfa
	return config
}

func getConfig(config *Config) *Config {
	//acl部分
	ips := ParseConfigMapStringSlice(config.IpConfigFile)
	id :=0

	config.AclMap = NewDnsCidr()
	config.AclNameIdMap = make(map[string]int)
	config.AclIdNameMap = make(map[int]string)

	for name, cidrs := range ips{
		log.Info("view name is ", name)
		config.AclNameIdMap[name] = id
		config.AclIdNameMap[id] = name

		for _, cidr := range cidrs{
			log.Info(cidr)
			config.AclMap.Insert(cidr, id)
		}
		id++
	}

	//domain name部分

	dnmap := ParseConfigMapStringSlice(config.DomainConfigFile)
	id = 0

	config.DnameList = NewDnameList()
	config.DLNameIdMap = make(map[string]int)
	config.DLIdNameMap = make(map[int]string)

	for dgname, dn := range dnmap{
		log.Info("domain group name is ", dgname)
		config.DLNameIdMap[dgname] = id
		config.DLIdNameMap[id] = dgname
		for _, dname := range dn{
			if !(dname[len(dname)-1:] == ".") {
				dname += "."
			}
			log.Info(dname)

			config.DnameList.Insert(dname, id)
		}
		id++
	}

	//ecs部分
	ecs := ParseConfigMapString(config.EcsMapConfigFile)

	config.AclEcsMap = make(map[int]string)

	for aclname, aclip := range ecs{
		log.Infof("aclname is: %s, ecs ip is: %s", aclname, aclip)
		config.AclEcsMap[config.AclNameIdMap[aclname]] = aclip
	}

	//fwd部分
	tfa := ParseFwdArray(config.FwdConfigFile)
//	for _, fa := range tfa.Forwarder{
//		fmt.Printf("%+v\n", fa)
//	}
	tfa.FwdMap = make(map[HitPoint]*Fwder)
	for i:=0; i<len(tfa.Forwarder);i++{
		log.Debugf("%+v\n", tfa.Forwarder[i])
		for _, aclname :=  range tfa.Forwarder[i].Acl{
			aclid, ok := config.AclNameIdMap[aclname]
			if !ok{
				log.Fatal("Fwd config has a missing Acl: ", aclname)
			}
			for _, dname := range tfa.Forwarder[i].Domains{
				dnameid, ok := config.DLNameIdMap[dname]
				if !ok{
					log.Fatal("Fwd config has a missing dname: ", dname)
				}
				tfa.FwdMap[HitPoint{aclid, dnameid}] = tfa.Forwarder[i]

			}
		}
	}
	config.Forwarders = tfa

	return config
}



func parseConfigJson(path string) *Config {

	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Open config file failed: ", err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("Read config file failed: ", err)
	}

	j := new(Config)
	err = json.Unmarshal(b, j)
	if err != nil {
		log.Fatal("Json syntex error: ", err)
	}

	return j
}
//读取dns adaptive forward配置文件，得到原始adaptive fwd列表
func ParseAdaptiveFwdArray(path string) *AdaptiveForwarderArray{
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Open config file failed: ", err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("Read config file failed: ", err)
	}

	afa := new(AdaptiveForwarderArray)
	err = json.Unmarshal(b, afa)
	if err != nil{
		log.Fatal("Json syntex error: ", err)
	}

	return afa
}

//读取dns forward配置文件，得到原始fwd列表
func ParseFwdArray(path string) *ForwarderArray{
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Open config file failed: ", err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("Read config file failed: ", err)
	}

	fa := new(ForwarderArray)
	err = json.Unmarshal(b, fa)
	if err != nil{
		log.Fatal("Json syntex error: ", err)
	}

	return fa
}
//从一个名值对中获取配置, 保存为map
//形如：(设置ecs ip)
//{
//  "beijing_acl" : "8.8.8.8",
//  "tianjin_acl" : "114.114.114.114"
//}
//
func ParseConfigMapString(configfile string) map[string]string{
	f, err := os.Open(configfile)
	if err != nil {
		log.Fatalf("File not found: %s\n", configfile)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil{
		log.Fatal("Read File failed: ", err)
	}

	var result map[string]string

	if err := json.Unmarshal(b, &result); err != nil {
		log.Fatal("json file format error", err)
	}

	return result
}

//从一个名称对应多值得配置文件中获取配置，保存为map，键值为string，值为string slice
//配置形如：
//{
//  "beijing_acl" : ["1.4.4.0/24", "1.2.2.0/24",
//    "123.112.0.0/12", "192.168.1.1/32"],
//  "tianjin_acl" : ["103.1.20.0/22", "192.168.2.1/32"]
//}
//
func ParseConfigMapStringSlice(configfile string) map[string][]string{
	f, err := os.Open(configfile)
	if err != nil {
		log.Fatalf("File not found: %s\n", configfile)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil{
		log.Fatal("Read File Failed: ", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Fatal("json file format error", err)
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
