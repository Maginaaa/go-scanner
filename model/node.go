package model

type NodeCollection struct {
	ApiList     ApiNodeList
	FuncList    FunctionNodeList
	StructList  StructNodeList
	FileList    FileNodeList
	PackageList PackageNodeList
	MicroServer MicroServerNode
	Domain      DomainNode
}

func NewNodeCollection() *NodeCollection {
	return &NodeCollection{
		ApiList:     NewApiNodeList(),
		FuncList:    NewFuncNodeList(),
		StructList:  NewStructNodeList(),
		FileList:    NewFileNodeList(),
		PackageList: NewPackageNodeList(),
	}
}

func (c NodeCollection) ToCypherList() []string {
	cypherList := make([]string, 0)
	cypherList = append(cypherList, c.Domain.ToCypher())
	cypherList = append(cypherList, c.MicroServer.ToCypher())
	cypherList = append(cypherList, c.ApiList.ToCypher()...)
	cypherList = append(cypherList, c.FuncList.ToCypher()...)
	cypherList = append(cypherList, c.StructList.ToCypher()...)
	cypherList = append(cypherList, c.FileList.ToCypher()...)
	cypherList = append(cypherList, c.PackageList.ToCypher()...)
	return cypherList
}

type ApiNode struct {
	NodeId int64  `json:"id"`
	Path   string `json:"path"`
	Type   string `json:"type"`
}

func (n *ApiNode) CreateCy() string {
	return CreateApiPathCy(n)
}

type ApiNodeList struct {
	nodes []ApiNode
}

func NewApiNodeList() ApiNodeList {
	return ApiNodeList{nodes: make([]ApiNode, 0)}
}

func (l *ApiNodeList) ToCypher() []string {
	set := NewSet()
	for _, node := range l.nodes {
		set.Add(node.CreateCy())
	}
	return set.KeySet()
}

func (l *ApiNodeList) Add(node ApiNode) {
	l.nodes = append(l.nodes, node)
}

func (l *ApiNodeList) Len() int {
	return len(l.nodes)
}

func (l *ApiNodeList) Append(list ApiNodeList) {
	l.nodes = append(l.nodes, list.nodes...)
}

type FunctionNode struct {
	NodeId  int64  `json:"id"`
	File    string `json:"file"`
	Folder  string `json:"folder"`
	Name    string `json:"name" `
	Content string `json:"content"`
	Rec     string `json:"rec"`
}

func (n *FunctionNode) ToCypher() string {
	return CreateFunctionCy(n)
}

type FunctionNodeList struct {
	nodes []FunctionNode
}

func (l *FunctionNodeList) Add(node FunctionNode) {
	if node.Name == "" {
		panic("function name is empty")
	}
	l.nodes = append(l.nodes, node)
}

func (l *FunctionNodeList) Len() int {
	return len(l.nodes)
}

func (l *FunctionNodeList) ToCypher() []string {
	set := NewSet()
	for _, node := range l.nodes {
		set.Add(node.ToCypher())
	}
	return set.KeySet()
}

func NewFuncNodeList() FunctionNodeList {
	return FunctionNodeList{nodes: make([]FunctionNode, 0)}
}

type StructNode struct {
	NodeId  int64  `json:"id"`
	File    string `json:"file"`
	Folder  string `json:"folder"`
	Content string `json:"content"`
	Name    string `json:"name"`
}

func (n *StructNode) ToCypher() string {
	return CreateStructCy(n)
}

type StructNodeList struct {
	nodes []StructNode
}

func (l *StructNodeList) ToCypher() []string {
	set := NewSet()
	for _, node := range l.nodes {
		set.Add(node.ToCypher())
	}
	return set.KeySet()
}

func (l *StructNodeList) Add(node StructNode) {
	if node.Name == "" {
		panic("struct name is empty")
	}
	l.nodes = append(l.nodes, node)
}

func (l *StructNodeList) Len() int {
	return len(l.nodes)
}

func NewStructNodeList() StructNodeList {
	return StructNodeList{nodes: make([]StructNode, 0)}
}

type FileNode struct {
	NodeId int64  `json:"id"`
	Name   string `json:"name"`
	Path   string `json:"path"`
}

func (n *FileNode) ToCypher() string {
	return CreateFileCy(n)
}

func NewFileNode(name, path string) FileNode {
	return FileNode{NodeId: 0, Name: name, Path: path}
}

type FileNodeList struct {
	nodes []FileNode
}

func NewFileNodeList() FileNodeList {
	return FileNodeList{nodes: make([]FileNode, 0)}
}

func (l *FileNodeList) ToCypher() []string {
	set := NewSet()
	for _, node := range l.nodes {
		set.Add(node.ToCypher())
	}
	return set.KeySet()
}

func (l *FileNodeList) Add(node FileNode) {
	l.nodes = append(l.nodes, node)
}

func (l *FileNodeList) Len() int {
	return len(l.nodes)
}

type PackageNode struct {
	NodeId int64  `json:"id"`
	Name   string `json:"name"`
	Path   string `json:"path"`
}

func NewPackageNode(name, path string) PackageNode {
	return PackageNode{NodeId: 0, Name: name, Path: path}
}

func (n *PackageNode) ToCypher() string {
	return CreatePackageCy(n)
}

type PackageNodeList struct {
	nodes []PackageNode
}

func (l *PackageNodeList) List() []PackageNode {
	return l.nodes
}

func NewPackageNodeList() PackageNodeList {
	return PackageNodeList{nodes: make([]PackageNode, 0)}
}

func (l *PackageNodeList) ToCypher() []string {
	set := NewSet()
	for _, node := range l.nodes {
		set.Add(node.ToCypher())
	}
	return set.KeySet()
}

func (l *PackageNodeList) Add(node PackageNode) {
	l.nodes = append(l.nodes, node)
}

func (l *PackageNodeList) Len() int {
	return len(l.nodes)
}

type MicroServerNode struct {
	NodeId int64  `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
	Path   string `json:"path"`
}

func (n *MicroServerNode) ToCypher() string {
	return CreateMicroServerCy(n)
}

func (n *MicroServerNode) Set(node MicroServerNode) {
	n.NodeId = node.NodeId
	n.Name = node.Name
	n.Domain = node.Domain
	n.Path = node.Path
}

//type MicroServerNodeList struct {
//	Nodes []MicroServerNode `json:"nodes"`
//}
//
//func NewMicroServerNodeList() MicroServerNodeList {
//	return MicroServerNodeList{Nodes: make([]MicroServerNode, 0)}
//}
//
//func (f *MicroServerNodeList) ToCypherList() []string {
//	set := NewSet()
//	for _, node := range f.Nodes {
//		set.Add(node.CreateCy())
//	}
//	return set.KeySet()
//}
//
//func (f *MicroServerNodeList) Add(node MicroServerNode) {
//	f.Nodes = append(f.Nodes, node)
//}

type DomainNode struct {
	NodeId int64  `json:"id"`
	Name   string `json:"name"`
}

func (n *DomainNode) CreateCy() string {
	return CreateDomainCy(n)
}

func (n *DomainNode) Set(node DomainNode) {
	n.NodeId = node.NodeId
	n.Name = node.Name
}

func (n *DomainNode) ToCypher() string {
	return CreateDomainCy(n)
}

//type DomainNodeList struct {
//	Nodes []DomainNode `json:"nodes"`
//}
//
//func NewDomainNodeList() DomainNodeList {
//	return DomainNodeList{Nodes: make([]DomainNode, 0)}
//}
//
//func (f *DomainNodeList) ToCypherList() []string {
//	set := NewSet()
//	for _, node := range f.Nodes {
//		set.Add(node.CreateCy())
//	}
//	return set.KeySet()
//}
//
//func (f *DomainNodeList) Add(node DomainNode) {
//	f.Nodes = append(f.Nodes, node)
//}

type Node struct {
	Id     int64                  `json:"id"` // 函数节点id
	Labels []string               `json:"labels"`
	Props  map[string]interface{} `json:"props"`
}
