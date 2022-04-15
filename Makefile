build:
	@go install

run:
	@btracer.exe

generate-proto:
	@protoc -I=btrace-idl/proto --go_out=proto-gen .\btrace-idl\proto\agent.proto
	@protoc -I=btrace-idl/proto --go-grpc_out=proto-gen .\btrace-idl\proto\agent.proto
