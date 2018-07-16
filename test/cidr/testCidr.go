package main

import (
	"github.com/yl2chen/cidranger"
	"net"
	"fmt"
	"RsDnsTools/util"
)

type Cidr_cdn struct{
	ipnet			net.IPNet
	cdn_id		int
}

func NewCidrCdn(ip net.IPNet, id int)(cidranger.RangerEntry){
	cc := &Cidr_cdn{ipnet: ip, cdn_id: id}
	return cc
}

func (cc *Cidr_cdn)Network()(net.IPNet){
	return cc.ipnet
}


func main(){
	/*
	ranger := cidranger.NewPCTrieRanger()

	_, network1, _ := net.ParseCIDR("192.168.1.0/24")
	_, network2, _ := net.ParseCIDR("128.168.1.0/24")


	ranger.Insert(NewCidrCdn(*network1, 1))
	ranger.Insert(NewCidrCdn(*network2, 2))

	contains, _ := ranger.Contains(net.ParseIP("192.168.1.5"))
	fmt.Printf("result is %t\r\n", contains)

	contains, _ = ranger.Contains(net.ParseIP("192.168.2.5"))
	fmt.Printf("result2 is %t\r\n", contains)


	entrys, err := ranger.ContainingNetworks(net.ParseIP("192.168.2.6"))

	if(err != nil) {
		fmt.Printf("Something error")
	}
	if(len(entrys)==0){
		fmt.Printf("doesn't match\r\n")
	}
	for _, entry := range entrys{
		fmt.Printf("get cdn id is %d", entry.(*Cidr_cdn).cdn_id)

	}
*/

	cl := util.NewDnsCidr()
	cl.Insert("192.168.1.0/24", 1)
	cl.Insert("192.168.2.0/24", 2)
	cl.Insert("192.168.3.0/24", 3)
	cl.Insert("192.168.4.3", 4)
	//cl.Insert("::/0", 5)
	//cl.Insert("0.0.0.0", 6)

	testIPs := []string{"192.168.1.34", "192.168.2.57", "192.168.3.22", "192.168.4.3","2008:fb::1"}
	for _, ip := range testIPs{
		fmt.Printf("ip:%s in cdn %d\r\n", ip, cl.Get(ip))
	}

}
