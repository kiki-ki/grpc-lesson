run:
	go run server/cmd/main.go

gen-proto:
	protoc --proto_path ./proto --go_out=${GOPATH}/src/github.com/kiki-ki/grpc-lesson/ ${FILENAME}
