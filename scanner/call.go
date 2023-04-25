package scanner

import (
	"github.com/Maginaaa/go-scanner/model"
	"go/token"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
	"log"
	"path/filepath"
	"strings"
)

// ServerScanner 服务级别的扫描
func (s *Scanner) callGraph() (err error) {

	// 生成Go Packages, 包含扁平的文件列表
	// config配置： Env：启动时的环境变量， Tests：是否扫描单测代码
	dir := filepath.Join(s.RootPath, s.MicroServerPath)
	initial, err := packages.Load(&packages.Config{
		Mode:  packages.LoadAllSyntax,
		Dir:   dir,
		Tests: false,
	})
	if err != nil {
		log.Printf("load: %v\n", err)
		return
	}

	// 生成ssa
	prog, pkgs := ssaUtilAllPackages(initial)
	prog.Build()

	// 找出main package
	var mains []*ssa.Package
	for _, p := range pkgs {
		if p != nil && p.Pkg.Name() == "main" && p.Func("main") != nil {
			mains = append(mains, p)
		}
	}

	// 以main为入口进行深度遍历，使用pointer生成完整调用链路
	config := &pointer.Config{
		Mains:          mains,
		BuildCallGraph: true,
	}
	result, err := pointer.Analyze(config)
	if err != nil {
		log.Printf("pointer.Analyze() err: %v,server: %s, dir: %s\n", err, s.MicroServerName, dir)
		return
	}
	// 遍历调用链路，获取项目所需结构化数据
	if result == nil {
		return
	}
	err = callgraph.GraphVisitEdges(result.CallGraph, func(edge *callgraph.Edge) error {
		err = s.BuildMap(edge)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Printf("callgraph.GraphVisitEdges() err: %v", err)
		return
	}

	s.serverToPkgInit()

	return
}

func ssaUtilAllPackages(initial []*packages.Package) (*ssa.Program, []*ssa.Package) {

	var fset *token.FileSet
	if len(initial) > 0 {
		fset = initial[0].Fset
	}

	prog := ssa.NewProgram(fset, 0)

	isInitial := make(map[*packages.Package]bool, len(initial))
	for _, p := range initial {
		isInitial[p] = true
	}

	ssaMap := make(map[*packages.Package]*ssa.Package)
	packages.Visit(initial, nil, func(p *packages.Package) {
		if p.Types != nil && !p.IllTyped {
			ssaMap[p] = prog.CreatePackage(p.Types, p.Syntax, p.TypesInfo, true)
		}
	})

	var ssaPkgs []*ssa.Package
	for _, p := range initial {
		ssaPkgs = append(ssaPkgs, ssaMap[p])
	}
	return prog, ssaPkgs
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
