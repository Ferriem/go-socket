package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number")
		return
	}
	PORT := ":" + arguments[1]
	s, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	c, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()
	buffer := make([]byte, 1024)
	rand.Seed(time.Now().Unix())

	for {
		n, addr, err := c.ReadFromUDP(buffer)
		fmt.Print("-> ", string(buffer[0:n-1]))
		if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
			fmt.Println("Exiting UDP server!")
			return
		}
		data := []byte(strconv.Itoa(random(1, 1001)))
		fmt.Printf("data: %s\n", string(data))
		_, err = c.WriteToUDP(data, addr)
		if err != nil {
			fmt.Println(err)
			return
		}

	}

}

func random(minn, maxn int) int {
	return rand.Intn(maxn-minn) + minn
}
