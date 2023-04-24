package scanner

import (
	"bytes"
	"fmt"
	"github.com/Maginaaa/go-scanner/model"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

// 文件扫描
// 结构体获取
// 函数内容获取
func (s *Scanner) fileContentScanner() {
	for _, path := range s.PathList.KeySet() {
		if !strings.HasSuffix(path, ".go") {
			return
		}
		fullPath := filepath.Join(s.RootPath, path)
		fileContent, err := ioutil.ReadFile(fullPath)
		if err != nil {
			fmt.Println("read file error : ", fullPath)
			continue
		}
		importList := make(map[string]string)
		// 当前文件相对路径
		//fileRelativePath := strings.TrimPrefix(path, s.RootPath+"/")
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, fullPath, nil, 0)
		if err != nil {
			fmt.Println("parse file error : ", fullPath)
			continue
		}
		sl := strings.Split(path, "/")
		pathL := sl[:len(sl)-1]
		folder := strings.Join(pathL, "/")
		fileName := sl[len(sl)-1]
		importList[f.Name.Name] = folder
		for _, stmt := range f.Imports {
			// import处理
			referencePath := strings.Trim(stmt.Path.Value, "\"") // 被调包路径
			pathList := strings.Split(referencePath, "/")
			referenceName := pathList[len(pathList)-1] // 被调用包的使用名称
			if nil != stmt.Name {
				// 别名
				referenceName = stmt.Name.Name
			}
			importList[referenceName] = referencePath
		}
		for _, decl := range f.Decls {
			switch stmt := decl.(type) {
			case *ast.FuncDecl:
				// 函数入参
				//stmt.Type.Params.List
				// 函数出参
				//stmt.Type.Results.List
				var buf bytes.Buffer
				_ = printer.Fprint(&buf, fset, stmt)
				// 函数名
				body := strconv.Quote(string(fileContent[int(stmt.Pos())-fset.Position(stmt.Pos()).Column : stmt.End()-1]))
				s.NodeCollection.FuncList.Add(model.FunctionNode{
					Name:      stmt.Name.Name,
					File:      path,
					Folder:    folder,
					Content:   body,
					StartLine: fset.Position(stmt.Pos()).Line,
					EndLine:   fset.Position(stmt.End()).Line,
				})
			case *ast.GenDecl:
				// 处理结构体
				if stmt.Tok != token.TYPE {
					continue
				}
				for _, spec := range stmt.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if _, ok = typeSpec.Type.(*ast.StructType); ok {
							structNode := model.StructNode{
								Name:    typeSpec.Name.Name,
								File:    path,
								Folder:  folder,
								Content: strconv.Quote(string(fileContent[int(typeSpec.Pos())-fset.Position(typeSpec.Pos()).Column : typeSpec.End()-1])),
							}
							s.NodeCollection.StructList.Add(structNode)
							s.LinkCollection.HasStructLinkList.Add(model.FileToStructLink{
								File: model.FileNode{
									Name: fileName,
									Path: path,
								},
								Struct: structNode,
							})
							// 结构体的字段处理
							//for _, field := range structType.Fields.List {
							//	if len(field.Names) == 0 {
							//		continue
							//	}
							//	//fmt.Println("Field Name:", field.Names[0].Name)
							//	//fmt.Println("Field Type:", field.Type)
							//}
						}
					}
				}
			}
		}
	}
}

func (s *Scanner) serverToPkgInit() {
	for _, pkg := range s.NodeCollection.PackageList.List() {
		if strings.HasPrefix(pkg.Path, s.NodeCollection.MicroServer.Path) {
			s.LinkCollection.HasPkgLinkList.Add(model.ServerToPkgLink{
				Server: s.NodeCollection.MicroServer,
				Pkg:    pkg,
			})
		}
	}
}
