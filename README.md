# eureka_tool

# Description
  call eureka api to sign any http service on eureka. i use it to sidecar node or tomcat app to springcloud.

  
# Example:

    eureka_tool sidecar --app hu_user --ip 172.16.9.17 --port 8080 --eureka-addr http://192.168.240.10:30761/eureka [--hostname 172.16.9.17] --ttl 30
    eureka_tool unregister --app hu_user --ip 172.16.9.17 --port 8080 --eureka-addr http://192.168.240.10:30761/eureka
	eureka_tool unregister --app hu_user --instanceid 172.16.9.17:8080 --eureka-addr http://192.168.240.10:30761/eureka

