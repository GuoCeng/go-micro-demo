package main

import (
	"context"
	"encoding/json"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/config/cmd"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	pb "go-micro-demo/proto/consignment"
	"io/ioutil"
	"log"
	"time"
)

const (
	addr            = "localhost:50021"
	defaultFileName = "consignment-cli/consignment.json"
)

func parseFile(fileName string) (*pb.Consignment, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("fail to read file: %v", err)
		return nil, err
	}
	var consignment *pb.Consignment
	//json解析传入指针，需要用指针的指针类型，才不会报错（因为直接传入指针的变量，未初始化默认值为nil,解析时，就算将数据保存下来，也没法传递给外面的函数，所以传入指针的指针，这样就可以传递出来了）
	if err := json.Unmarshal(bytes, &consignment); err != nil {
		log.Fatalf("fail to convert file to json: %v", err)
		return nil, err
	}
	return consignment, nil
}

func main() {
	//grpc客户端实现代码
	/*conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect server: %v", err)
		return
	}
	defer conn.Close()
	client := pb.NewShippingServiceClient(conn)
	if consignment, err := parseFile(defaultFileName); err != nil {
		return
	} else {
		_, err := client.CreateConsignment(context.Background(), consignment)
		resp, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
		if err != nil {
			log.Fatalf("grpc fail: %v", err)
			return
		}
		log.Printf("response: %v", resp)
	}*/

	//go-micro实现客户端代码
	cmd.Init()
	reg := consul.NewRegistry(
		func(op *registry.Options) {
			op.Addrs = []string{
				"127.0.0.1:8500",
			}
		},
	)
	s := micro.NewService(
		micro.Registry(reg),
		micro.Name("client"),
		micro.Version("latest"),
	)
	go func() {
		time.Sleep(1 * time.Second)
		client := pb.NewShippingServiceClient("", s.Client())
		if consignment, err := parseFile(defaultFileName); err != nil {
			return
		} else {
			_, err := client.CreateConsignment(context.Background(), consignment)
			resp, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
			if err != nil {
				log.Fatalf("grpc fail: %v", err)
				return
			}
			log.Printf("response: %v", resp)
		}
	}()
	if err := s.Run(); err != nil {
		log.Fatalf("micro server start fail: %v", err)
	}
}
