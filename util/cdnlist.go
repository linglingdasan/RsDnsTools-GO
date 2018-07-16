package util

import (
	"github.com/yl2chen/cidranger"
	"net"
	"strings"
)

type dns_cidr struct{
	ranger cidranger.Ranger
}

type dns_cidr_entry struct{
	ipnet			net.IPNet
	cdn_id		int
}

func NewDnsCidrEntry(ip net.IPNet, id int)(cidranger.RangerEntry){
	cc := &dns_cidr_entry{ipnet: ip, cdn_id: id}
	return cc
}

func (cc *dns_cidr_entry)Network()(net.IPNet){
	return cc.ipnet
}

func NewDnsCidr() *dns_cidr {
	cl := &dns_cidr{}
	cl.ranger = cidranger.NewPCTrieRanger()
	return cl
}

func (cl *dns_cidr)Insert(cidrStr string, cdn_id int)bool{

	i := strings.IndexByte(cidrStr, '/')
	if(i<0){
		if(strings.IndexByte(cidrStr, ':')<0){//IPv4
			if(cidrStr == "0.0.0.0"){
				cidrStr = cidrStr+"/0"
			}else{
				cidrStr = cidrStr+"/32"
			}
		}else{//IPv6
			if(cidrStr == "::0"){
				cidrStr = cidrStr+"/0"
			}else {
				cidrStr = cidrStr + "/128"
			}
		}
	}

	_, ipnet, err:= net.ParseCIDR(cidrStr)
	if(err != nil){
		return false
	}

	err = cl.ranger.Insert(NewDnsCidrEntry(*ipnet, cdn_id))
	if(err != nil){
		return false
	}
	return true
}

func (cl *dns_cidr)Get(ipstr string)(cdn_id int){
	ipnet := net.ParseIP(ipstr)
	entrys, err:=cl.ranger.ContainingNetworks(ipnet)
	if(err!=nil || len(entrys)==0){
		return -1
	}
	return entrys[0].(*dns_cidr_entry).cdn_id

}