package main

import (
	"strconv"

	"github.com/ArthurHlt/go-eureka-client/eureka"
)

type EurekaInstance eureka.InstanceInfo

//type EurekaInstance eureka.InstanceInfo

func (instance *EurekaInstance) ToInstanceInfo() *eureka.InstanceInfo {
	return (*eureka.InstanceInfo)(instance)
}

func CreateEurekaInstance(instanceinfo *eureka.InstanceInfo) *EurekaInstance {
	return (*EurekaInstance)(instanceinfo)
}

func NewInstance(hostname string, app string, ip string, port int, ttl uint, isSsl bool, homepage_url *string, healthcheck_url *string) *EurekaInstance {
	instance := CreateEurekaInstance(eureka.NewInstanceInfo(hostname, app, ip, port, ttl, isSsl)) //Create a new instance to register
	instance.Metadata = &eureka.MetaData{
		Map: make(map[string]string),
	}
	//instance.Metadata.Map["InstanceID"] = "172.16.13.13:8080" //add metadata for example
	instance_id := ip + ":" + strconv.Itoa(port)
	instance.InstanceID = instance_id
	if homepage_url != nil {
		//instance.HomePageUrl = "http://172.16.13.13:8080/"
		instance.HomePageUrl = *homepage_url
	}
	if healthcheck_url != nil {
		//instance.HealthCheckUrl = "http://172.16.13.13:8080/user/application.wadl"
		instance.HealthCheckUrl = *healthcheck_url
	}
	instance.VipAddress = app
	//instance.SecureVipAddress = "hu-user"
	//eureka_instance := EurekaInstance{InstanceInfo: *instance}
	//return &eureka_instance
	return instance
}
