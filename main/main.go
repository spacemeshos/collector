package main

import (
	"flag"
	"fmt"
	"github.com/spacemeshos/collector"
	"github.com/spacemeshos/collector/api"
	"github.com/spacemeshos/smutil/log"
)

const (
	defaultGRPCServerPort = 9081
	defaultJSONServerPort = 9080
)

func main() {
	/*url := "tcp://localhost:56565"
	username := "postgres"
	pass := "mysecretpassword"
	grpcPort := defaultGRPCServerPort
	httpPort := defaultJSONServerPort*/

	url := flag.String("url", "tcp://localhost:56565", "url from which events will be received")
	postgresUrl := flag.String("pg_url", "localhost:5432", "postgres url")
	username := flag.String("pg_uname", "postgres", "postgres username")
	pass := flag.String("pg_passwd", "mysecretpassword", "postgres password")
	grpcPort := flag.Int("grpc-port", defaultGRPCServerPort, "start grpc on this port")
	httpPort := flag.Int("http-port", defaultJSONServerPort, "start http server on this port")

	flag.Parse()

	db := collector.NewDb(*username, *pass, *postgresUrl)
	err := db.Start()
	if err != nil {
		log.Error("cannot create DB %v ", err)
		return
	}
	defer db.Close()
	c := collector.NewCollector(db, *url)

	grpcService := api.NewGrpcService(*grpcPort, db)
	jsonService := api.NewJSONHTTPServer(*httpPort)
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
