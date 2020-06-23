package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const root_usage = `Usage:

%s [command]

Commands:
    register
    sidecar
    heartbeat
    unregister
`
const usage1 = `Usage:
  
  %s %s --app [app] --ip [ip] --port [port] --eureka-addr [eureka_addr] [--hostname hostname] --ttl [ttl]
  
  Example:
    eureka_tool sidecar --app hu-user --ip 172.16.9.17 --port 8080 --eureka-addr http://192.168.240.10:30761/eureka [--hostname 172.16.9.17] --ttl 30

`
const usage2 = `Usage:
  
  %s %s --app app [--ip ip --port port | --instanceid instanceid] --eureka-addr eureka_addr
  
  Example:
    eureka_tool unregister --app hu-user --ip 172.16.9.17 --port 8080 --eureka-addr http://192.168.240.10:30761/eureka
	or
	eureka_tool unregister --app hu-user --instanceid 172.16.9.17:8080 --eureka-addr http://192.168.240.10:30761/eureka
	
`

type ActionConfig struct {
	action      string
	app         string
	ip          string
	hostname    string
	instanceid  string
	port        int
	ttl         int
	eureka_addr string
}

func (config *ActionConfig) format_eureka_addr(eureka_addr_string string) string {
	b := strings.Split(eureka_addr_string, ",")
	new_addr := ""
	for i, d := range b {
		if i == 0 {
			new_addr = strings.TrimRight(d, "/")
		} else {
			new_addr = new_addr + "," + strings.TrimRight(d, "/")
		}
	}
	return new_addr

}

func (config *ActionConfig) SetUpFlagSet(action string) *flag.FlagSet {
	switch action {
	case "root":
		flagset := flag.NewFlagSet(action, flag.ExitOnError)
		flagset.Usage = func() {
			fmt.Printf(root_usage, os.Args[0])
			flagset.PrintDefaults()
		}
		return flagset
	case "register", "sidecar":
		flagset := flag.NewFlagSet(action, flag.ExitOnError)
		flagset.Usage = func() {
			fmt.Printf(usage1, os.Args[0], action)
			flagset.PrintDefaults()
		}
		config.action = action
		flagset.StringVar(&config.app, "app", "", "app is appname, it must have")
		flagset.StringVar(&config.ip, "ip", "", "ip is app's ip, it must have")
		flagset.StringVar(&config.hostname, "hostname", "", "hostname default to ip")
		flagset.IntVar(&config.port, "port", 8080, "port is int, it must have")
		flagset.IntVar(&config.ttl, "ttl", 30, "ttl default 30")
		flagset.StringVar(&config.eureka_addr, "eureka-addr", "", "eureka_addr is string, it must have")
		return flagset
	case "unregister", "heartbeat":
		flagset := flag.NewFlagSet(action, flag.ExitOnError)
		flagset.Usage = func() {
			fmt.Printf(usage2, os.Args[0], action)
			flagset.PrintDefaults()
		}
		flagset.StringVar(&config.app, "app", "", "app is appname, it must have")
		flagset.StringVar(&config.ip, "ip", "", "ip is app's ip, ip:port is instanceid")
		flagset.IntVar(&config.port, "port", 8080, "port is int, ip:port is instanceid")
		flagset.StringVar(&config.eureka_addr, "eureka-addr", "", "eureka_addr is string, it must have")
		flagset.StringVar(&config.instanceid, "instanceid", "", "instanceid is combined by ip:port, it‘s priorityo higher than ip:port")
		return flagset
	}
	return nil
}

func parseArgs() *ActionConfig {
	config := new(ActionConfig)
	menu := config.SetUpFlagSet("root")
	if len(os.Args) > 1 {
		action := strings.ToLower(os.Args[1])
		switch action {
		case "register", "sidecar":
			f := config.SetUpFlagSet(action)
			if len(os.Args) <= 2 {
				f.Usage()
				return nil
			}
			if len(os.Args) > 2 {
				if err := f.Parse(os.Args[2:]); err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing params %s, error: %v", os.Args[2:], err)
					return nil
				}
				if config.app == "" && config.eureka_addr == "" && config.ip == "" {
					f.Usage()
					return nil
				}
				if config.hostname == "" {
					config.hostname = config.ip
				}
				config.instanceid = config.ip + ":" + strconv.Itoa(config.port)

			}
		case "unregister", "heartbeat":
			f := config.SetUpFlagSet(action)
			if len(os.Args) <= 2 {
				f.Usage()
				return nil
			}
			if len(os.Args) > 2 {
				if err := f.Parse(os.Args[2:]); err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing params %s, error: %v", os.Args[2:], err)
					return nil
				}
				if config.app == "" && config.eureka_addr == "" {
					f.Usage()
					return nil
				}
				if config.instanceid == "" {
					config.instanceid = config.ip + ":" + strconv.Itoa(config.port)
				}
			}
		default:
			fmt.Println("Invalid command")
			menu.Usage()
			return nil
		}
	} else {
		menu.Usage()
		return nil
	}
	return config
}
