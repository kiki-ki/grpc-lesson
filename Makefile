run:
	go run server/cmd/main.go

APP_ROOT=${GOPATH}/src/github.com/kiki-ki/grpc-lesson/
gen-proto:
	protoc --proto_path ./proto --go_out=plugins=grpc:${APP_ROOT} ${FILENAME}

gen-doc:
	protoc --doc_out=./doc --doc_opt=html,index.html ./proto/*.proto
