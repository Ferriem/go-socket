package main

import (
	"fmt"
	"net"
)

func main() {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range interfaces {
		fmt.Printf("Interface: %v\n", v.Name)
		byName, err := net.InterfaceByName(v.Name)
		if err != nil {
			fmt.Println(err)
			return
		}
		addresses, err := byName.Addrs()
		for k, v := range addresses {
			fmt.Printf("Interface Address #%v: %v\n", k, v.String())
		}
		fmt.Println()
	}
}
