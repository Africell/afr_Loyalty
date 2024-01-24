package Lendme

import (
	"log"
	"net"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

//var IsTarget_Up float64

type destinationHost struct {
	HostName string
	HostIP   string
	HostPort string
}

var destinationHots []destinationHost

func Init_DestinationHosts() {
	host := destinationHost{
		HostName: "MongoDB",
		HostIP:   Configuration.MongoDB.HostIP_1,
		HostPort: Configuration.MongoDB.HostPort_1,
	}
	destinationHots = append(destinationHots, host)

	host = destinationHost{
		HostName: "Recharge Feed",
		HostIP:   "10.10.231.50", //container name
		HostPort: "9658",
	}
	destinationHots = append(destinationHots, host)

	host = destinationHost{
		HostName: "SMSC",
		HostIP:   "10.10.215.52",
		HostPort: "15403",
	}
	destinationHots = append(destinationHots, host)

	host = destinationHost{
		HostName: "WhatsApp",
		HostIP:   "10.10.231.52",
		HostPort: "9010",
	}
	destinationHots = append(destinationHots, host)
}

func PortlinkInquiry() {
	log.Println("link inquiry started")
	Init_DestinationHosts()
	for range time.Tick(time.Second * 15) {
		for _, host := range destinationHots {
			hoststat := linkInquiry(host.HostIP, host.HostPort)
			PortStatus.With(prometheus.Labels{"DestinationHost": host.HostName}).Set(hoststat)
			// if host.HostIP == "SubscriberProfiler" {
			// 	if IsTarget_Up != hoststat {
			// 		IsTarget_Up = hoststat
			// 	}
			// }
		}
	}
}

func linkInquiry(host, port string) (linkStatus float64) {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		//fmt.Println("Connecting error:", err)
		linkStatus = 2 //down
		return
	}
	if conn != nil {
		defer conn.Close()
		//fmt.Println("Opened", net.JoinHostPort(host, port))
		linkStatus = 1 //up
		return
	} else {
		linkStatus = 2 //down
		return
	}
}
