package main

import (
	"fmt"
	"github.com/spacemeshos/collector"
	"github.com/spacemeshos/collector/api"
	"github.com/spacemeshos/smutil/log"
	"os"
)

func main() {
	url := "tcp://localhost:56565"
	username := "postgres"
	pass := "mysecretpassword"
	if len(os.Args) > 1 {
		url = os.Args[1]
	}
	if len(os.Args) > 3 {
		username = os.Args[2]
		pass = os.Args[3]
	}
	db := collector.NewDb(username, pass)
	err := db.Start()
	if err != nil {
		log.Error("cannot create DB %v ", err)
		return
	}
	defer db.Close()
	c := collector.NewCollector(db, url)

	grpcService := api.NewGrpcService(db)
	jsonService := api.NewJSONHTTPServer()
	waitGrpc := make(chan bool)
	waitHttp := make(chan bool)

	grpcService.StartService(waitGrpc)
	jsonService.StartService(waitHttp)

	<-waitGrpc
	defer grpcService.StopService()
	<-waitHttp
	defer jsonService.StopService()

	fmt.Print("Start collecting events")
	c.Start(true)
}
