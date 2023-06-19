#/bin/bash

GOHOSTOS=(go env GOHOSTOS)
GOPATH=(go env GOPATH)
VERSION=(git describe --tags --always)
# 寻找 gotags 注释的文件


if [ $GOHOSTOS == "windows" ]; then
    	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
    	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
else
  	API_PROTO_FILES=$(find api -name *.proto)
fi

echo $API_PROTO_FILES

function init() {
    	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    	go install github.com/llmuz/ijk/cmd/protoc-gen-go-gin@latest
    	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
    	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
    	go install github.com/favadi/protoc-go-inject-tag@v1.4.0
}


function api() {
    	protoc --proto_path=./api \
    	       -I third_party \
     	       --go_out=paths=source_relative:./api \
    	       --openapi_out=fq_schema_naming=true,default_response=false:. \
    	       $API_PROTO_FILES

}

function gotags() {
    API_PROTO_FILES_GO_TAGS=$(find api -name *.pb.go)
    for file in $API_PROTO_FILES_GO_TAGS
    do
      echo $file
      protoc-go-inject-tag -input=$file
    done
}

function ginServer() {
    protoc -I ./api  --proto_path=./third_party  --go-gin_out ./api --go-gin_opt=paths=source_relative   $API_PROTO_FILES

}


init
api
gotags
ginServer
