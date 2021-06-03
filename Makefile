run-unary-server:
	go run server/unary/cmd/main.go

run-unary-client:
	go run client/unary/cmd/main.go

run-server-streaming-server:
	go run server/server_streaming/cmd/main.go

run-server-streaming-client:
	go run client/server_streaming/cmd/main.go

run-client-streaming-server:
	go run server/client_streaming/cmd/main.go

run-client-streaming-client:
	go run client/client_streaming/cmd/main.go

run-bidirectional-streaming-server:
	go run server/bidirectional_streaming/cmd/main.go

run-bidirectional-streaming-client:
	go run client/bidirectional_streaming/cmd/main.go

APP_ROOT=${GOPATH}/src/github.com/kiki-ki/grpc-lesson/

gen-proto:
	protoc --proto_path ./proto --go_out=plugins=grpc:${APP_ROOT} ${FILENAME}
	make gen-doc

gen-doc:
	protoc --doc_out=./proto/doc --doc_opt=html,index.html ./proto/*.proto

lint:
	golangci-lint run

fmt:
	goimports -w ./
	go vet ./...
	go fmt ./...
	make lint
