package api

import (
	"context"
	"github.com/spacemeshos/collector/api/pb"
	"github.com/spacemeshos/go-spacemesh/events"
	"github.com/spacemeshos/go-spacemesh/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"strconv"
)

type Database interface {
	GetTransactionsFrom(orig string) ([]*events.NewTx, error)
	GetTransactionsTo(orig string) ([]*events.NewTx, error)
}

// GrpcService is a grpc server providing the collector api
type GrpcService struct {
	Server   *grpc.Server
	dataBase Database
	port     int
}

// NewGrpcService create a new grpc service using config data.
func NewGrpcService(port int, dataBase Database) *GrpcService {
	server := grpc.NewServer()
	return &GrpcService{Server: server, dataBase: dataBase, port: port}
}

// StartService starts the grpc service.
func (s *GrpcService) StartService(status chan bool) {
	go s.startServiceInternal(status)
}

// StopService stops the grpc service.
func (s GrpcService) StopService() {
	log.Debug("Stopping grpc service...")
	s.Server.Stop()
	log.Debug("grpc service stopped...")

}

// This is a blocking method designed to be called using a go routine
func (s *GrpcService) startServiceInternal(status chan bool) {
	port := s.port
	addr := ":" + strconv.Itoa(int(port))

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		//log.Error("failed to listen", err)
		return
	}

	pb.RegisterCollectorServiceServer(s.Server, s)

	// SubscribeOnNewConnections reflection service on gRPC server
	reflection.Register(s.Server)

	log.Debug("grpc API listening on port %d", port)

	if status != nil {
		status <- true
	}

	// start serving - this blocks until err or server is stopped
	if err := s.Server.Serve(lis); err != nil {
		log.Error("grpc stopped serving", err)
	}

	if status != nil {
		status <- true
	}

}

// Echo returns the response for an echo api request
func (s GrpcService) Echo(ctx context.Context, in *pb.SimpleMessage) (*pb.SimpleMessage, error) {
	return &pb.SimpleMessage{Value: in.Value}, nil
}

func (s GrpcService) GetTransactionsByAccount(ctx context.Context, in *pb.SimpleMessage) (*pb.Txs, error) {
	txsFrom, err := s.dataBase.GetTransactionsFrom(in.Value)
	if err != nil {
		return nil, err
	}
	txsTo, err := s.dataBase.GetTransactionsTo(in.Value)
	if err != nil {
		return nil, err
	}
	resp := make([]*pb.Tx, 0, len(txsFrom)+len(txsTo))
	for _, tx := range txsFrom {
		resp = append(resp, &pb.Tx{Id: tx.Id, Amount: tx.Amount, Origin: tx.Origin, Destination: tx.Destination, Gas: tx.Gas})
	}
	for _, tx := range txsTo {
		resp = append(resp, &pb.Tx{Id: tx.Id, Amount: tx.Amount, Origin: tx.Origin, Destination: tx.Destination, Gas: tx.Gas})
	}

	return &pb.Txs{Txs: resp}, nil
}
