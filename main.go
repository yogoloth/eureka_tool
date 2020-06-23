package main

import (
	"fmt"
	"os"
)

func main() {
	config := parseArgs()
	if config == nil {
		fmt.Println("args error,exit.")
		return
	}

	app := config.app
	ip := config.ip
	port := config.port
	eureka_addr := config.eureka_addr
	hostname := config.hostname
	fmt.Printf("%v\n", config)

	//app := "hu-user"
	//ip := "172.16.9.17"
	//port := 8080
	//eureka_addr := "http://192.168.240.10:30761/eureka"
	//hostname := ip

	instance := NewInstance(hostname, app, ip, port, 30, false)

	var myClient *EurekaClient
	if c, err := NewClient(eureka_addr, instance); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	} else {
		myClient = c
	}

	switch config.action {
	case "register", "sidecar":
		if err := myClient.Register(); err != nil {
			fmt.Printf("regiter %v\n", err)
		} else {
			fmt.Println("register instance success.")
		}
	case "unregister":
		if err := myClient.UnRegister(); err != nil {
			fmt.Printf("regiter %v\n", err)
		} else {
			fmt.Println("register instance success.")
		}
	case "heartbeat":
		//todo
	default:
		fmt.Println("should not happen.")
	}
}
