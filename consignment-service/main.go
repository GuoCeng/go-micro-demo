package main

import (
	"context"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	pb "go-micro-demo/proto/consignment"
	"log"
)

type IRepository interface {
	Create(consignment *pb.Consignment) (*pb.Consignment, error)
	GetAll() ([]*pb.Consignment, error)
}

type Repository struct {
	consignments []*pb.Consignment
}

func (r *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	r.consignments = append(r.consignments, consignment)
	return consignment, nil
}

func (r *Repository) GetAll() ([]*pb.Consignment, error) {
	return r.consignments, nil
}

type service struct {
	repo IRepository
}

func (s *service) CreateConsignment(ctx context.Context, consignment *pb.Consignment, resp *pb.Response) error {
	if result, err := s.repo.Create(consignment); err != nil {
		resp = &pb.Response{
			Created: false,
		}
		return err
	} else {
		resp = &pb.Response{
			Created:     true,
			Consignment: result,
		}
		return nil
	}
}

func (s *service) GetConsignments(ctx context.Context, getRequest *pb.GetRequest, resp *pb.Response) error {
	consignments, _ := s.repo.GetAll()
	resp = &pb.Response{
		Created:      true,
		Consignment:  nil,
		Consignments: consignments,
	}
	return nil
}

func main() {
	//直接使用grpc版本，注意自动生成的proto文件与使用go-micro的不一样
	/*repo := new(Repository)
	lis, err := net.Listen("tcp", ":50021")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterShippingServiceServer(s, &service{repo})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}*/

	reg := consul.NewRegistry(
		func(op *registry.Options) {
			op.Addrs = []string{"127.0.0.1:8500"}
		},
	)

	//使用go-micro
	s := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
		micro.Registry(reg),
	)
	s.Init()
	pb.RegisterShippingServiceHandler(s.Server(), &service{new(Repository)})
	if err := s.Run(); err != nil {
		log.Fatalf("micro server start fail: %v", err)
	}
}
