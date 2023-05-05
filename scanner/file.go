package scanner

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Maginaaa/go-scanner/model"
)

// 文件扫描
// 结构体获取
// 函数内容获取
func (s *Scanner) fileContentScanner() {
	for _, path := range s.PathList.KeySet() {
		if filepath.Ext(path) != ".go" {
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
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, fullPath, nil, 0)
		if err != nil {
			log.Println("parse file error : ", fullPath)
			continue
		}
		folder := filepath.Dir(path)

		importList[f.Name.Name] = folder

		fileName := filepath.Base(path)
		for _, stmt := range f.Imports {
			// import处理
			referencePath := strings.Trim(stmt.Path.Value, "\"") // 被调包路径 // 被调用包的使用名称
			referenceName := filepath.Base(referencePath)
			if nil != stmt.Name {
				// 别名
				referenceName = stmt.Name.Name
			}
			if s.FilterDependency {
				if !strings.HasPrefix(referencePath, s.MicroServerPath) {
					continue
				}
			}
			importList[referenceName] = referencePath
		}
		for _, decl := range f.Decls {
			switch stmt := decl.(type) {
			case *ast.FuncDecl:
				var buf bytes.Buffer
				_ = printer.Fprint(&buf, fset, stmt)
				funcName := stmt.Name.Name
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
					structList := extractStruct(path, folder, importList, stmt.Recv.List)
					for _, st := range structList {
						s.LinkCollection.FuncReceiverList.Add(model.FuncReceiverLink{
							Func:   functionNode,
							Struct: st,
						})
					}
				}

				// 函数入参
				if stmt.Type.Params != nil {
					structList := extractStruct(path, folder, importList, stmt.Type.Params.List)
					for _, st := range structList {
						s.LinkCollection.FuncParamList.Add(model.FuncParamLink{
							Func:   functionNode,
							Struct: st,
						})
					}
				}

				// 函数返回
				if stmt.Type.Results != nil {
					structList := extractStruct(path, folder, importList, stmt.Type.Results.List)
					for _, st := range structList {
						s.LinkCollection.FuncReturnList.Add(model.FuncReturnLink{
							Func:   functionNode,
							Struct: st,
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

// 将三种结构体处理逻辑进行封装（接收器、参数、返回）
// 通常情况下，接收器不存在多个结构体的情况，不存在跨包调用的情况
func extractStruct(path, folder string, importList map[string]string, list []*ast.Field) []model.StructNode {
	structList := make([]model.StructNode, 0)
	var err error
	for _, param := range list {
		structName := ""
		structPath := path
		structPath, structName, err = getStructFromTypeExpr(param.Type, path, importList)
		if err != nil {
			continue
		}
		structList = append(structList, model.StructNode{
			Name:   structName,
			File:   structPath,
			Folder: folder,
		})
	}
	return structList
}

func getStructFromTypeExpr(expr ast.Expr, path string, importList map[string]string) (structPath, structName string, err error) {
	structPath = path
	switch stmt := expr.(type) {
	case *ast.Ident:
		structName = stmt.Name
		break
	case *ast.SelectorExpr:
		pkgName := stmt.X.(*ast.Ident).Name
		if value, ok := importList[pkgName]; !ok {
			return "", "", errors.New("continue")
		} else {
			structPath = value
			structName = stmt.Sel.Name
		}
		break
		// 指针型参数
	case *ast.StarExpr:
		structPath, structName, err = getStructFromTypeExpr(stmt.X, path, importList)
		break
	case *ast.ArrayType:
		structPath, structName, err = getStructFromTypeExpr(stmt.Elt, path, importList)
		break
	case *ast.Ellipsis, *ast.MapType:
		// TODO: ...[]struct 与 map[]struct 暂不处理
		return "", "", errors.New("*ast.Ellipsis, *ast.MapType continue")
	case *ast.ChanType, *ast.FuncType, *ast.InterfaceType:
		return "", "", errors.New("*ast.ChanType, *ast.FuncType, *ast.InterfaceType continue")
	case *ast.StructType:
		return "", "", errors.New("*ast.StructType continue")
	default:
		return "", "", errors.New("default continue")
	}
	if structName == "" {
		return "", "", errors.New("continue")
	}
	return
}
