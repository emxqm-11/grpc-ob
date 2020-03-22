proto_generate: 
	protoc -I/usr/local/include -I. \
	-I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.13.0/third_party/googleapis \
	--go_out=plugins=grpc:. \
  	./proto/ob-service.proto

custom_generate: 
	protoc -I/usr/local/include -I. \
	-I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.13.0/third_party/googleapis \
	--go_out=plugins=grpc:. \
  	./aggregator/customob.proto

swagger_to_proto:
	openapi2proto -spec cds_full.json -annotate -out ob-service.proto

server.start:
	go run main.go

client.start:
	go run client/client.go
