module github.com/devil-dwj/db/cmd/protoc-gen-go-sql

go 1.19

require (
	github.com/devil-dwj/db v0.0.0-20230718062448-ac7e4054d94f
	google.golang.org/protobuf v1.31.0
)

require github.com/iancoleman/strcase v0.3.0

replace github.com/devil-dwj/db => ../../
