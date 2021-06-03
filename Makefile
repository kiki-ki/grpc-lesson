run-server:
	go run server/cmd/main.go

run-client:
	go run client/cmd/main.go

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
