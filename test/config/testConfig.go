package main

import (
	"RsDnsTools/util"
	"fmt"
)

func main() {
	ips := util.ParseConfigMapStringSlice("configs/cdn.ips.json")

	for name, cidrs := range ips{
		println("view name is ", name)
		for _, cidr := range cidrs{
			println(cidr)
		}
	}

	dnmap := util.ParseConfigMapStringSlice("configs/cdn.names.json")
	for name, dn := range dnmap{
		println("domain names is ", name)
		for _, names := range dn{
			println(names)
		}
	}

	ecsmap := util.ParseConfigMapString("configs/ecs.map.json")
	for name, ecs := range ecsmap{
		println(name, ecs)
	}

	tfa := util.ParseFwdArray("configs/fwd.json")
	fmt.Printf("%+v\n", tfa.Forwarder[0])
	//println(tfa.Forwarder)

}



