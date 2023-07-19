package main

import (
	"fmt"
	"strings"

	"github.com/devil-dwj/db/sql"
	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

const (
	dbPackage      = protogen.GoImportPath("github.com/devil-dwj/db")
	contextPackage = protogen.GoImportPath("context")
)

func generateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {

	filename := file.GeneratedFilenamePrefix + ".sql.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-go-sql. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()

	generateFileContent(gen, file, g)

	return g
}

func generateFileContent(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile) {
	if len(file.Messages) == 0 {
		return
	}

	origFileName := file.Proto.GetName()
	lend := strings.Split(origFileName, "/")
	trimName := strings.TrimSuffix(lend[len(lend)-1], ".proto")
	goTypeName := strcase.ToCamel(trimName)
	goStructName := goTypeName + "Sql"
	g.P("type ", goStructName, " interface {")
	for _, message := range file.Messages {
		if !proto.HasExtension(message.Desc.Options(), sql.E_Sql) {
			continue
		}
		g.P(message.Comments.Leading,
			interfaceSignature(gen, file, g, message),
		)
	}
	g.P("}")
	g.P()

	// struct impl
	goStructImplName := goTypeName + "_Sql"
	g.P("type ", goStructImplName, " struct {")
	g.P("db *", dbPackage.Ident("DB"))
	g.P("}")

	// new
	g.P("func New", goStructName, "(db *", dbPackage.Ident("DB"), ") ", goStructName, "{")
	g.P("return &", goStructImplName, "{db: db}")
	g.P("}")
	g.P()

	// GetRawDB
	// g.P("func ", getEnclosureIdent(goStructImplName), "GetRawDB() *", gormPackage.Ident("DB"), "{")
	// g.P("return r.db")
	// g.P("}")
	// g.P()

	// impl methods
	for _, message := range file.Messages {
		if !proto.HasExtension(message.Desc.Options(), sql.E_Sql) {
			continue
		}

		implFuncSignature(gen, file, g, message, goStructImplName)
	}
}

func implFuncSignature(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, message *protogen.Message, goStructImplName string) {
	var req = ""
	var rsp = ""
	var result = ""
	var total = ""
	var totalPage = ""
	var pageTotal = ""
	for _, field := range message.Fields {

		if field.GoName == "Req" && field.Message != nil {
			req = "req *" + g.QualifiedGoIdent(field.Message.GoIdent)
		}
		if field.GoName == "Rsp" && field.Message != nil {
			if field.Desc.IsList() {
				rsp = "*[]*" + g.QualifiedGoIdent(field.Message.GoIdent) + ", "
			} else {
				rsp = "*" + g.QualifiedGoIdent(field.Message.GoIdent) + ", "
			}
		}
		if field.GoName == "Result" {
			result = "int" + ", "
		}
		if field.GoName == "TotalCount" {
			total = "int, "
		}
		if field.GoName == "TotalPage" {
			totalPage = "int, "
		}
		if field.GoName == "PageTotal" {
			pageTotal = "int, "
		}
	}

	g.P("func ", getEnclosureIdent(goStructImplName), message.GoIdent.GoName, "(ctx context.Context, ", req, ") (", rsp, result, total, totalPage, pageTotal, "error) {")
	// call prps
	callSqlSignature(gen, file, g, message)
	g.P("}")
	g.P()
}

func callSqlSignature(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, message *protogen.Message) {
	resIdent := ""
	ivTempl := ""
	sql := getSqlRule(message)
	ivTempl = "\"" + sql.Raw + "\""
	for _, field := range message.Fields {
		if field.GoName == "Rsp" && field.Message != nil {
			if field.Desc.IsList() {
				resIdent = "res := &[]*%s{}"
			} else {
				resIdent = "res := &%s{}"
			}
			resIdent = fmt.Sprintf(resIdent, g.QualifiedGoIdent(field.Message.GoIdent))
		}

		if field.GoName == "Req" {
			for _, field1 := range field.Message.Fields {
				ivTempl += ", req." + strcase.ToCamel(field1.GoName)
			}
		}
	}

	g.P(resIdent)
	g.P("err := r.db.WithContext(ctx).")
	g.P("Raw(", ivTempl, ").")
	g.P("Scan(res).")
	g.P("Error")
	g.P("return res, err")
}

func getSqlRule(message *protogen.Message) *sql.Sql {
	r := proto.GetExtension(message.Desc.Options(), sql.E_Sql)
	return r.(*sql.Sql)
}

func getEnclosureIdent(name string) string {
	return "(r *" + name + ")"
}

func interfaceSignature(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, message *protogen.Message) string {
	var req = ""
	var rsp = ""
	var result = ""
	var total = ""
	for _, field := range message.Fields {
		if field.GoName == "Req" && field.Message != nil {
			req = "*" + g.QualifiedGoIdent(field.Message.GoIdent)
		}
		if field.GoName == "Rsp" && field.Message != nil {
			if field.Desc.IsList() {
				rsp = "*[]*" + g.QualifiedGoIdent(field.Message.GoIdent) + ", "
			} else {
				rsp = "*" + g.QualifiedGoIdent(field.Message.GoIdent) + ", "
			}
		}
		if field.GoName == "Result" {
			result = "int" + ", "
		}
		if field.GoName == "TotalCount" {
			total = "int, "
		}
	}

	return message.GoIdent.GoName + "(" + g.QualifiedGoIdent(contextPackage.Ident("Context")) + ", " + req + ") (" + rsp + result + total + "error)"
}