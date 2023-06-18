#!/bash/bash

# 安装必要的工具
function init() {
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/favadi/protoc-go-inject-tag@v1.4.0
}

# 根据 protobuf 生成 pb.go 已经 tag
function build() {
    # 获取 errors protos 文件
    ERROR_PROTO_FILES=$(find errors -name *.proto)
    protoc --go_out  ./  --go_opt=paths=source_relative $ERROR_PROTO_FILES
    # 获取 errors 目录下的pb.go 文件
    API_PROTO_FILES_GO_TAGS=$(find errors -name *.pb.go)
    # 处理 pb.go
    for file in $API_PROTO_FILES_GO_TAGS
    do
      protoc-go-inject-tag -input=$file
    done
}

init
build
