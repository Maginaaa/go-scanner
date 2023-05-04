package model

type LinkCollection struct {
	MappingLinkList     ApiToFuncLinkList
	CallLinkList        FuncCallFuncLinkList
	HasFunctionLinkList FileToFuncLinkList
	HasStructLinkList   FileToStructLinkList
	HasFileLinkList     PkgToFileLinkList
	HasPkgLinkList      ServerToPkgLinkList
	HasServerLink       DomainToServerLink
	FuncReceiverList    FuncReceiverLinkList
	FuncParamList       FuncParamLinkList
	FuncReturnList      FuncReturnLinkList
}

func (l LinkCollection) ToCypherList() []string {
	cypher := make([]string, 0)
	cypher = append(cypher, l.HasServerLink.ToCypher())
	cypher = append(cypher, l.CallLinkList.ToCypherList()...)
	cypher = append(cypher, l.MappingLinkList.ToCypherList()...)
	cypher = append(cypher, l.HasFunctionLinkList.ToCypherList()...)
	cypher = append(cypher, l.HasStructLinkList.ToCypherList()...)
	cypher = append(cypher, l.HasFileLinkList.ToCypherList()...)
	cypher = append(cypher, l.HasPkgLinkList.ToCypherList()...)
	cypher = append(cypher, l.FuncReceiverList.ToCypherList()...)
	cypher = append(cypher, l.FuncParamList.ToCypherList()...)
	cypher = append(cypher, l.FuncReturnList.ToCypherList()...)
	return cypher
}

func NewLinkCollection() *LinkCollection {
	return &LinkCollection{
		CallLinkList:        NewFuncCallFuncLinkList(),
		MappingLinkList:     NewApiToFuncLinkList(),
		HasFunctionLinkList: NewFileToFuncLinkList(),
		HasStructLinkList:   NewFileToStructLinkList(),
		HasPkgLinkList:      NewServerToPkgLinkList(),
		HasFileLinkList:     NewPkgToFileLinkList(),
		FuncReceiverList:    NewFuncReceiverLinkList(),
		FuncParamList:       NewFuncParamLinkList(),
		FuncReturnList:      NewFuncReturnLinkList(),
	}
}

type DomainToServerLink struct {
	Domain DomainNode
	Server MicroServerNode
}

func (l *DomainToServerLink) ToCypher() string {
	return DomainToServerCy(l)
}

func (l *DomainToServerLink) Set(link DomainToServerLink) {
	l.Domain = link.Domain
	l.Server = link.Server
}

type DomainToServerLinkList struct {
	Links []DomainToServerLink `json:"links"`
}

func (l *DomainToServerLinkList) Add(link DomainToServerLink) {
	l.Links = append(l.Links, link)
}

func (l *DomainToServerLinkList) ToCypherList() []string {
	cypherList := NewSet()
	for _, link := range l.Links {
		cypherList.Add(link.ToCypher())
	}
	return cypherList.KeySet()
}

type ServerToPkgLink struct {
	Server MicroServerNode
	Pkg    PackageNode
}

func (l *ServerToPkgLink) CreateCy() string {
	return ServerToPkgCy(l)
}

type ServerToPkgLinkList struct {
	Links []ServerToPkgLink `json:"links"`
}

func NewServerToPkgLinkList() ServerToPkgLinkList {
	return ServerToPkgLinkList{Links: make([]ServerToPkgLink, 0)}
}

func (l *ServerToPkgLinkList) Add(link ServerToPkgLink) {
	// TODO：服务和包有多对多情况，需修复
	l.Links = append(l.Links, link)
}

func (l *ServerToPkgLinkList) ToCypherList() []string {
	cypherList := NewSet()
	for _, link := range l.Links {
		cypherList.Add(link.CreateCy())
	}
	return cypherList.KeySet()
}

type PkgToFileLink struct {
	Pkg  PackageNode
	File FileNode
}

func (l *PkgToFileLink) ToCypher() string {
	return PkgToFileCy(l)
}

type PkgToFileLinkList struct {
	Links []PkgToFileLink `json:"links"`
}

func NewPkgToFileLinkList() PkgToFileLinkList {
	return PkgToFileLinkList{Links: make([]PkgToFileLink, 0)}
}

func (l *PkgToFileLinkList) Add(link PkgToFileLink) {
	l.Links = append(l.Links, link)
}

func (l *PkgToFileLinkList) ToCypherList() []string {
	cypherList := NewSet()
	for _, link := range l.Links {
		cypherList.Add(link.ToCypher())
	}
	return cypherList.KeySet()
}

type FileToFuncLink struct {
	File FileNode
	Func FunctionNode
}

func (l *FileToFuncLink) ToCypher() string {
	return FileToFunctionCy(l)
}

type FileToFuncLinkList struct {
	Links []FileToFuncLink `json:"links"`
}

func NewFileToFuncLinkList() FileToFuncLinkList {
	return FileToFuncLinkList{Links: make([]FileToFuncLink, 0)}
}

func (l *FileToFuncLinkList) Add(link FileToFuncLink) {
	l.Links = append(l.Links, link)
}

func (l *FileToFuncLinkList) ToCypherList() []string {
	cypherList := NewSet()
	for _, link := range l.Links {
		cypherList.Add(link.ToCypher())
	}
	return cypherList.KeySet()
}

type FileToStructLink struct {
	File   FileNode
	Struct StructNode
}

func (l *FileToStructLink) ToCypher() string {
	return FileToStructCy(l)
}

type FileToStructLinkList struct {
	Links []FileToStructLink `json:"links"`
}

func NewFileToStructLinkList() FileToStructLinkList {
	return FileToStructLinkList{Links: make([]FileToStructLink, 0)}
}

func (l *FileToStructLinkList) Add(link FileToStructLink) {
	l.Links = append(l.Links, link)
}

func (l *FileToStructLinkList) ToCypherList() []string {
	cypherList := NewSet()
	for _, link := range l.Links {
		cypherList.Add(link.ToCypher())
	}
	return cypherList.KeySet()
}

type FuncReceiverLink struct {
	Func   FunctionNode
	Struct StructNode
}

func (l *FuncReceiverLink) ToCypher() string {
	return FuncReceiverCy(l)
}

type FuncReceiverLinkList struct {
	Links []FuncReceiverLink `json:"links"`
}

func NewFuncReceiverLinkList() FuncReceiverLinkList {
	return FuncReceiverLinkList{Links: make([]FuncReceiverLink, 0)}
}

func (l *FuncReceiverLinkList) Add(link FuncReceiverLink) {
	l.Links = append(l.Links, link)
}

func (l *FuncReceiverLinkList) ToCypherList() []string {
	cypherList := NewSet()
	for _, link := range l.Links {
		cypherList.Add(link.ToCypher())
	}
	return cypherList.KeySet()
}

type FuncParamLink struct {
	Func   FunctionNode
	Struct StructNode
}

func (l *FuncParamLink) ToCypher() string {
	return FuncParamCy(l)
}

type FuncParamLinkList struct {
	Links []FuncParamLink `json:"links"`
}

func NewFuncParamLinkList() FuncParamLinkList {
	return FuncParamLinkList{Links: make([]FuncParamLink, 0)}
}

func (l *FuncParamLinkList) Add(link FuncParamLink) {
	l.Links = append(l.Links, link)
}

func (l *FuncParamLinkList) ToCypherList() []string {
	cypherList := NewSet()
	for _, link := range l.Links {
		cypherList.Add(link.ToCypher())
	}
	return cypherList.KeySet()
}

type FuncReturnLink struct {
	Func   FunctionNode
	Struct StructNode
}

func (l *FuncReturnLink) ToCypher() string {
	return FuncReturnCy(l)
}

type FuncReturnLinkList struct {
	Links []FuncReturnLink `json:"links"`
}

func NewFuncReturnLinkList() FuncReturnLinkList {
	return FuncReturnLinkList{Links: make([]FuncReturnLink, 0)}
}

func (l *FuncReturnLinkList) Add(link FuncReturnLink) {
	l.Links = append(l.Links, link)
}

func (l *FuncReturnLinkList) ToCypherList() []string {
	cypherList := NewSet()
	for _, link := range l.Links {
		cypherList.Add(link.ToCypher())
	}
	return cypherList.KeySet()
}

type ApiToFuncLink struct {
	Func FunctionNode
	Api  ApiNode
}

func (l *ApiToFuncLink) ToCypher() string {
	return ApiToFunctionCy(l)
}

type ApiToFuncLinkList struct {
	Links []ApiToFuncLink `json:"links"`
}

func NewApiToFuncLinkList() ApiToFuncLinkList {
	return ApiToFuncLinkList{Links: make([]ApiToFuncLink, 0)}
}

func (l *ApiToFuncLinkList) Add(link ApiToFuncLink) {
	l.Links = append(l.Links, link)
}

func (l *ApiToFuncLinkList) Append(list ApiToFuncLinkList) {
	l.Links = append(l.Links, list.Links...)
}

func (l *ApiToFuncLinkList) ToCypherList() []string {
	cypherList := NewSet()
	for _, link := range l.Links {
		cypherList.Add(link.ToCypher())
	}
	return cypherList.KeySet()
}

type FuncCallFuncLink struct {
	Caller FunctionNode
	Callee FunctionNode
}

func (f *FuncCallFuncLink) ToCypher() string {
	return FuncCallFuncCy(f)
}

type FuncCallFuncLinkList struct {
	Links []FuncCallFuncLink `json:"links"`
}

func NewFuncCallFuncLinkList() FuncCallFuncLinkList {
	return FuncCallFuncLinkList{Links: make([]FuncCallFuncLink, 0)}
}

func (l *FuncCallFuncLinkList) Add(link FuncCallFuncLink) {
	l.Links = append(l.Links, link)
}

func (l *FuncCallFuncLinkList) Append(list FuncCallFuncLinkList) {
	l.Links = append(l.Links, list.Links...)
}

func (l *FuncCallFuncLinkList) ToCypherList() []string {
	cypherList := NewSet()
	for _, link := range l.Links {
		cypherList.Add(link.ToCypher())
	}
	return cypherList.KeySet()
}
