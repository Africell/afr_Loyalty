package Lendme

import (
	"fmt"
	"log"
	"net"
	"strings"
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
	// if Configuration.LoyaltyMongoDB.HostIP_1 != "" {
	// 	host := destinationHost{
	// 		HostName: "Loyalty MongoDB 01",
	// 		HostIP:   Configuration.LoyaltyMongoDB.HostIP_1,
	// 		HostPort: Configuration.LoyaltyMongoDB.HostPort_1,
	// 	}
	// 	destinationHots = append(destinationHots, host)
	// }
	// if Configuration.LoyaltyMongoDB.HostIP_2 != "" {
	// 	host := destinationHost{
	// 		HostName: "Loyalty MongoDB 02",
	// 		HostIP:   Configuration.LoyaltyMongoDB.HostIP_2,
	// 		HostPort: Configuration.LoyaltyMongoDB.HostPort_2,
	// 	}
	// 	destinationHots = append(destinationHots, host)
	// }
	// if Configuration.LoyaltyMongoDB.HostIP_3 != "" {
	// 	host := destinationHost{
	// 		HostName: "Loyalty MongoDB 03",
	// 		HostIP:   Configuration.LoyaltyMongoDB.HostIP_3,
	// 		HostPort: Configuration.LoyaltyMongoDB.HostPort_3,
	// 	}
	// 	destinationHots = append(destinationHots, host)
	// }
	uri := Configuration.Mongo.URI
	uri = strings.TrimPrefix(uri, "mongodb+srv://")
	uri = strings.TrimPrefix(uri, "mongodb://")
	if i := strings.Index(uri, "/"); i >= 0 {
		uri = uri[:i]
	}
	if i := strings.Index(uri, "?"); i >= 0 {
		uri = uri[:i]
	}
	for idx, hp := range strings.Split(uri, ",") {
		hp = strings.TrimSpace(hp)
		if hp == "" {
			continue
		}
		host, port, err := net.SplitHostPort(hp)
		if err != nil {
			continue
		}
		destinationHots = append(destinationHots, destinationHost{
			HostName: fmt.Sprintf("Loyalty MongoDB %02d", idx+1),
			HostIP:   host,
			HostPort: port,
		})
	}

	if Configuration.App_AUC.Hostname != "" {
		host := destinationHost{
			HostName: "AUC",
			HostIP:   Configuration.App_AUC.Hostname,
			HostPort: Configuration.App_AUC.Port,
		}
		destinationHots = append(destinationHots, host)
	}

	if Configuration.CGW_AUC.Hostname != "" {
		host := destinationHost{
			HostName: "UCGW AUC",
			HostIP:   Configuration.CGW_AUC.Hostname,
			HostPort: Configuration.CGW_AUC.Port,
		}
		destinationHots = append(destinationHots, host)
	}

	if Configuration.CGW.Hostname != "" {
		host := destinationHost{
			HostName: "UCGW",
			HostIP:   Configuration.CGW.Hostname,
			HostPort: Configuration.CGW.Port,
		}
		destinationHots = append(destinationHots, host)
	}

	if Configuration.OKAPI_AUC.Hostname != "" {
		host := destinationHost{
			HostName: "OKAPI AUC",
			HostIP:   Configuration.OKAPI_AUC.Hostname,
			HostPort: Configuration.OKAPI_AUC.Port,
		}
		destinationHots = append(destinationHots, host)
	}

	if Configuration.Propylaea.Hostname != "" {
		host := destinationHost{
			HostName: "Propylaea",
			HostIP:   Configuration.Propylaea.Hostname,
			HostPort: Configuration.Propylaea.Port,
		}
		destinationHots = append(destinationHots, host)
	}

	if Configuration.IN.IP != "" {
		host := destinationHost{
			HostName: "IN Web Services",
			HostIP:   Configuration.IN.IP, //container name
			HostPort: Configuration.IN.Port,
		}
		destinationHots = append(destinationHots, host)
	}

	if Configuration.SMPP.IP != "" {
		host := destinationHost{
			HostName: "SMSC",
			HostIP:   Configuration.SMPP.IP,
			HostPort: Configuration.SMPP.Port,
		}
		destinationHots = append(destinationHots, host)
	}
	if Configuration.APGW.Hostname != "" {
		host := destinationHost{
			HostName: Configuration.APGW.Description,
			HostIP:   Configuration.APGW.Hostname,
			HostPort: Configuration.APGW.Port,
		}
		destinationHots = append(destinationHots, host)
	}

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
