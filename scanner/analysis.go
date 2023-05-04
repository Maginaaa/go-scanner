package scanner

import (
	"errors"
	"fmt"
	"github.com/Maginaaa/go-scanner/model"
	"golang.org/x/tools/go/callgraph"
	"path/filepath"
	"strings"
)

// BuildMap
// cypherList 保存所有函数节点， callList函数节点的调用关系
func (s *Scanner) BuildMap(edge *callgraph.Edge) error {

	// 过滤非源代码
	if edge.Caller.Func.Pkg == nil || edge.Callee.Func.Synthetic != "" {
		return nil
	}

	caller := edge.Caller
	callee := edge.Callee
	callerFileName := filepath.Base(edge.Caller.Func.Prog.Fset.Position(caller.Func.Pos()).Filename)

	// . 排除默认 init
	if s.FilterInit {
		if callerFileName == "." {
			return nil
		}
	}

	// 调用方
	callerNode, err := s.makeSet(*caller)
	if err != nil {
		return nil
	}

	// 被调用方
	calleeNode, err := s.makeSet(*callee)
	if err != nil {
		return nil
	}

	// 创建 函数-call->函数关系
	s.LinkCollection.CallLinkList.Add(model.FuncCallFuncLink{Caller: callerNode, Callee: calleeNode})

	return nil

}

func (s *Scanner) makeSet(node callgraph.Node) (funcNode model.FunctionNode, err error) {

	// 包路径
	callerPkgPath := node.Func.Pkg.Pkg.Path()
	// 过滤非源代码
	if s.FilterDependency {
		if !strings.HasPrefix(callerPkgPath, s.ProjectName) {
			return funcNode, errors.New("delete node")
		}
	}
	prog := node.Func.Prog
	pkgName := node.Func.Pkg.Pkg.Name()
	// 文件名
	fileName := filepath.Base(prog.Fset.Position(node.Func.Pos()).Filename)
	if s.ProjectName != s.ProjectPath {
		callerPkgPath = strings.Replace(callerPkgPath, s.ProjectName, s.ProjectPath, 1)
	}
	filePath := fmt.Sprintf("%s/%s", callerPkgPath, fileName)
	// 函数名
	funcName := node.Func.Name()
	if funcName == "ExtractApi" {
		fmt.Println(false)
	}
	// 匿名函数处理
	if strings.Contains(funcName, "$") {
		ss := strings.Split(funcName, "$")
		funcName = ss[0]
	}
	// 创建包节点
	packageNode := model.NewPackageNode(pkgName, callerPkgPath)
	s.NodeCollection.PackageList.Add(packageNode)
	// 创建文件节点
	fileNode := model.NewFileNode(fileName, filePath)
	s.NodeCollection.FileList.Add(fileNode)

	// 保存后续文件扫描需要的文件列表
	s.PathList.Add(filePath)

	funcNode = model.FunctionNode{
		Name:      funcName,
		File:      filePath,
		StartLine: prog.Fset.Position(node.Func.Syntax().Pos()).Line,
		EndLine:   prog.Fset.Position(node.Func.Syntax().End()).Line,
	}

	// 创建 包-contains->文件关系
	s.LinkCollection.HasFileLinkList.Add(model.PkgToFileLink{Pkg: packageNode, File: fileNode})

	// 创建 文件-contains-> 函数关系
	s.LinkCollection.HasFunctionLinkList.Add(model.FileToFuncLink{File: fileNode, Func: funcNode})

	return funcNode, nil

}
