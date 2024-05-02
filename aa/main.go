package main

import (
	"fmt"
	"net"
)

func main() {
	Ips()
	// for k, v := range ips {
	// 	fmt.Println(k, v)
	// }
}

func Ips() (map[string]string, error) {

	ips := make(map[string]string)

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range interfaces {

		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Name:%v, HardwareAddr:%v, MTU:%v\n", byName.Name, byName.HardwareAddr, byName.MTU)
		byName.
		addresses, err := byName.Addrs()
		for _, v := range addresses {
			ips[byName.Name] = v.String()
		}
	}
	return ips, nil
}
