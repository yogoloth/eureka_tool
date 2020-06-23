package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ArthurHlt/go-eureka-client/eureka"
)

type EurekaClient struct {
	eureka.Client
	retry_max int
	//client    *eureka.Client
	instance *EurekaInstance
}

func NewClient(eureka_url string, instance *EurekaInstance) (*EurekaClient, error) {
	eureka_servers := strings.Split(eureka_url, ",")
	//c := new(Client)
	eureka_client := eureka.NewClient(eureka_servers)
	if eureka_client == nil {
		err_string := fmt.Sprintf("eureka_servers format error: %s \n", eureka_url)
		return nil, errors.New(err_string)
	}
	//c := new(EurekaClient{Client: *eureka_client, retry_max: 3, instance: instance})
	c := EurekaClient{Client: *eureka_client, retry_max: 3, instance: instance}
	return &c, nil
}

func (c *EurekaClient) UnRegister() error {
	return c.UnregisterInstance(c.instance.App, c.instance.InstanceID)
}

func (c *EurekaClient) Register() error {
	exit_flag := false
	instance := c.instance

	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for ss := range s {
			switch ss {
			case syscall.SIGINT:
				//exit_flag = true
				fallthrough
			case syscall.SIGTERM:
				exit_flag = true
			}
		}
	}()

	for {
		if exit_flag {
			fmt.Printf("(debug) register got SIGTERM.\n")
			c.UnRegister()
			return errors.New("term by signal")
		}
		//client.GetApplication(app)                     // retrieve the application "test"
		for {
			retry_count := 0
			if _, err := c.GetInstance(instance.App, instance.InstanceID); err != nil {
				fmt.Printf("(debug) get instance error: %v \n", err)
				//if err := c.RegisterInstance(instance.App, &instance.InstanceInfo); err != nil {
				if err := c.RegisterInstance(instance.App, instance.ToInstanceInfo()); err != nil {
					fmt.Printf("(debug) register instance error: %v\n", err)
				} else {
					fmt.Printf("(debug) register instance success.\n", err)
					break
				}
			} else {
				break
			}
			retry_count++
			if retry_count > c.retry_max {
				err_string := fmt.Sprintf("(error) register instance error %d > %d times", retry_count, c.retry_max)
				return errors.New(err_string)
			}
		}
		fmt.Printf("(debug) Got instance: %v\n", instance)
		if err := c.SendHeartbeat(instance.App, instance.InstanceID); err != nil {
			fmt.Printf("(debug) sendheartbeat fail: %v\n", err)
		} else {
			fmt.Printf("(debug) sendheartbeat success\n")
		}
		time.Sleep(5 * time.Second)
	}
	return nil
}
