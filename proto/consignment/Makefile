build:

    # 生成go-micro相关的grpc代码  -I指定路径，后面如果有路径可以配置成-I设置的相对路径
	protoc -I=./ --go_out=plugins=micro:. consignment.proto
	#生成grpc相关代码
	#protoc -I=./ --go_out=plugins=grpc:. consignment.proto

test:
	# micro工具查询服务
	micro --registry consul --registry_address 127.0.0.1:8500  get service go.micro.srv.consignment
	#micro --registry consul --registry_address 127.0.0.1:8500  get service client