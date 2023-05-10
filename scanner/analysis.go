package scanner

import (
	"errors"
	"github.com/Maginaaa/go-scanner/model"
	"golang.org/x/tools/go/callgraph"
	"path/filepath"
	"regexp"
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
		if callerFileName == "." || strings.Contains(caller.Func.Name(), "init") {
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

	// 过滤非源代码
	if s.FilterDependency {
		if !strings.HasPrefix(node.Func.Pkg.Pkg.Path(), s.ProjectName) {
			return funcNode, errors.New("FilterDependency delete node")
		}
	}

	prog := node.Func.Prog
	pkgName := node.Func.Pkg.Pkg.Name()
	// 文件名
	fileAbsPath := prog.Fset.Position(node.Func.Pos()).Filename
	fileName := filepath.Base(fileAbsPath)

	fileRelativePath, _ := filepath.Rel(s.RootPath, fileAbsPath)
	callerPkgPath := filepath.Dir(fileRelativePath)

	// 函数名
	funcName := ""
	startLine, endLine := 0, 0
	// 匿名函数处理
	if strings.Contains(node.Func.Name(), "$") || node.Func.Parent() != nil {
		parentNode := node.Func.Parent()
		funcName = parentNode.Name()
		startLine = parentNode.Prog.Fset.Position(parentNode.Syntax().Pos()).Line
		endLine = parentNode.Prog.Fset.Position(parentNode.Syntax().End()).Line
	} else {
		funcName = node.Func.Name()
		startLine = prog.Fset.Position(node.Func.Syntax().Pos()).Line
		endLine = prog.Fset.Position(node.Func.Syntax().End()).Line
	}

	// 自定义过滤
	if len(s.FilterCustomize) > 0 {
		for _, reg := range s.FilterCustomize {
			if regexp.MustCompile(reg).MatchString(fileRelativePath) {
				return funcNode, errors.New("FilterCustomize delete node")
			}
		}
	}
	// 创建包节点
	packageNode := model.NewPackageNode(pkgName, callerPkgPath)
	s.NodeCollection.PackageList.Add(packageNode)
	// 创建文件节点
	fileNode := model.NewFileNode(fileName, fileRelativePath)
	s.NodeCollection.FileList.Add(fileNode)

	// 保存后续文件扫描需要的文件列表
	s.PathList.Add(fileRelativePath)

	funcNode = model.FunctionNode{
		Name:      funcName,
		File:      fileRelativePath,
		StartLine: startLine,
		EndLine:   endLine,
	}

	// 创建 包-contains->文件关系
	s.LinkCollection.HasFileLinkList.Add(model.PkgToFileLink{Pkg: packageNode, File: fileNode})

	// 创建 文件-contains-> 函数关系
	s.LinkCollection.HasFunctionLinkList.Add(model.FileToFuncLink{File: fileNode, Func: funcNode})

	return funcNode, nil

}
