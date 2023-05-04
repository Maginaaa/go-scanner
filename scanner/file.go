package scanner

import (
	"bytes"
	"github.com/Maginaaa/go-scanner/model"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"log"
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
			log.Println("read file error : ", fullPath)
			continue
		}
		importList := make(map[string]string)
		// 当前文件相对路径
		//fileRelativePath := strings.TrimPrefix(path, s.RootPath+"/")
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, fullPath, nil, 0)
		if err != nil {
			log.Println("parse file error : ", fullPath)
			continue
		}
		sl := strings.Split(path, "/")
		pathList := sl[:len(sl)-1]
		folder := strings.Join(pathList, "/")
		importList[f.Name.Name] = folder

		fileName := sl[len(sl)-1]
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
				funcName := stmt.Name.Name
				if funcName == "ExtractApi" {
					log.Println(buf.String())
				}
				// 函数名
				body := strconv.Quote(string(fileContent[int(stmt.Pos())-fset.Position(stmt.Pos()).Column : stmt.End()-1]))
				functionNode := model.FunctionNode{
					Name:      funcName,
					File:      path,
					Folder:    folder,
					Content:   body,
					StartLine: fset.Position(stmt.Pos()).Line,
					EndLine:   fset.Position(stmt.End()).Line,
				}
				s.NodeCollection.FuncList.Add(functionNode)

				// 函数接收器
				if stmt.Recv != nil {
					for _, recv := range stmt.Recv.List {
						structName := ""
						switch recObj := recv.Type.(type) {
						case *ast.StarExpr:
							// 指针接收器
							structName = recObj.X.(*ast.Ident).Name
						case *ast.Ident:
							// 值接收器
							structName = recObj.Name
						}
						s.LinkCollection.FuncReceiverList.Add(model.FuncReceiverLink{
							Func: functionNode,
							Struct: model.StructNode{
								Name:   structName,
								File:   path,
								Folder: folder,
							},
						})
					}
				}

				// 函数入参
				if stmt.Type.Params != nil {
					for _, param := range stmt.Type.Params.List {
						structName := ""
						switch recObj := param.Type.(type) {
						case *ast.StarExpr:
							// 指针接收器
							structName = recObj.X.(*ast.Ident).Name
						case *ast.Ident:
							// 值接收器
							structName = recObj.Name
						}
						s.LinkCollection.FuncParamList.Add(model.FuncParamLink{
							Func: functionNode,
							Struct: model.StructNode{
								Name:   structName,
								File:   path,
								Folder: folder,
							},
						})
					}
				}

				// 函数返回
				if stmt.Type.Results != nil {
					for _, result := range stmt.Type.Results.List {
						structName := ""
						switch recObj := result.Type.(type) {
						case *ast.StarExpr:
							// 指针接收器
							structName = recObj.X.(*ast.Ident).Name
						case *ast.Ident:
							// 值接收器
							structName = recObj.Name
						}
						s.LinkCollection.FuncReturnList.Add(model.FuncReturnLink{
							Func: functionNode,
							Struct: model.StructNode{
								Name:   structName,
								File:   path,
								Folder: folder,
							},
						})
					}
				}

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
							//	//log.Println("Field Name:", field.Names[0].Name)
							//	//log.Println("Field Type:", field.Type)
							//}
						}
					}
				}
			}
		}
	}
}
