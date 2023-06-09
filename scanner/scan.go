package scanner

import (
	"errors"
	"github.com/Maginaaa/go-scanner/config"
	"github.com/Maginaaa/go-scanner/model"
	"log"
	"time"
)

type Scanner struct {
	Domain           string                   // 域名， 默认为服务名
	ProjectName      string                   // 项目名(必传，go.mod文件中的module)
	ProjectPath      string                   // 项目包名，默认同等于ProjectName
	MicroServerName  string                   // 微服务名,如果是单包扫描，同等于服务名
	MicroServerPath  string                   // 服务目录， 默认为服务名(考虑到项目内含有多服务的情况，如果是单包，同等于服务名)
	RootPath         string                   // 文件保存路径(必传)
	FilterInit       bool                     // 是否过滤init函数
	FilterDependency bool                     // 是否过滤三方依赖包
	CustomFilterPath []string                 // 过滤一些自定义规则，需传入正则表达式，符合规则的路径将被过滤
	CustomFunc       []func(scanner *Scanner) // 自定义处理逻辑
	NodeCollection   *model.NodeCollection
	LinkCollection   *model.LinkCollection
	PathList         model.Set
	Neo4jConfig      config.Neo4jConfig
}

func (s *Scanner) init() error {
	if s.RootPath == "" {
		return errors.New("RootPath is empty")
	}
	if s.ProjectName == "" {
		return errors.New("ProjectName is empty")
	}
	if s.Domain == "" {
		s.Domain = s.ProjectName
	}
	if s.ProjectPath == "" {
		s.ProjectPath = s.ProjectName
	}
	if s.MicroServerPath == "" {
		s.MicroServerPath = s.ProjectPath
	}
	if s.MicroServerName == "" {
		s.MicroServerName = s.ProjectName
	}
	s.configInit()
	s.NodeCollection = model.NewNodeCollection()
	s.LinkCollection = model.NewLinkCollection()
	s.PathList = model.NewSet()
	// 域节点、服务节点和关系的初始化
	s.dataInit()
	return nil
}

func (s *Scanner) ServerScanner() (err error) {
	start := time.Now()
	// 步骤 1
	err = s.init()
	log.Printf("%s 服务开始处理\n", s.MicroServerName)
	if err != nil {
		return err
	}
	// 步骤 2
	err = s.callGraph()
	if err != nil {
		return err
	}

	// 步骤 3
	s.fileContentScanner()

	// 步骤 4 自定义处理逻辑
	if len(s.CustomFunc) > 0 {
		for _, f := range s.CustomFunc {
			f(s)
		}
	}

	log.Printf("%s 处理耗时完整耗时：%s\n", s.MicroServerName, time.Since(start))
	return nil
}

// 域节点、服务节点及关系初始化
func (s *Scanner) dataInit() {
	domainNodeInit(s)
	serverNodeInit(s)
	domainToServerInit(s)
}

func domainNodeInit(s *Scanner) {
	s.NodeCollection.Domain.Set(model.DomainNode{Name: s.Domain})
}

func serverNodeInit(s *Scanner) {
	s.NodeCollection.MicroServer.Set(model.MicroServerNode{
		Name:   s.MicroServerName,
		Domain: s.Domain,
		Path:   s.MicroServerPath,
	})
}

func domainToServerInit(s *Scanner) {
	s.LinkCollection.HasServerLink.Set(model.DomainToServerLink{
		Server: s.NodeCollection.MicroServer,
		Domain: s.NodeCollection.Domain,
	})
}

func (s *Scanner) configInit() {
	config.Neo4jCfg.Url = s.Neo4jConfig.Url
	config.Neo4jCfg.UserName = s.Neo4jConfig.UserName
	config.Neo4jCfg.Password = s.Neo4jConfig.Password
}
