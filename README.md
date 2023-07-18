# protoc-gen-go-procedure

## 作用

通过定义proto文件生成存储过程执行函数

## 安装

```go
go get -u github.com/devil-dwj/wms-bee/db/cmd/protoc-gen-go-procedure
```

## 使用

```go
protoc --proto_path=. --go-procedure_out=paths=source_relative:. *_pr.proto
```
