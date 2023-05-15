package model

import (
	"fmt"
	"log"
)

const (
	createDomain      = "MERGE (:Domain { name: '%s'})"
	createMicroServer = "MERGE (:MicroServer { name: '%s', domain: '%s', path: '%s' })"
	createApiPath     = "MERGE (:Api {path: '%s', type: '%s'})"
	createPackage     = "MERGE (:Package { name: '%s', path: '%s' })"
	createFile        = "MERGE (:File { name: '%s', path: '%s' })"
	createFunction    = "MERGE (:Function {name:'%s', file: '%s', folder: '%s', rec : '%s', content : %s})"
	createStruct      = "MERGE (:Struct {name:'%s', file: '%s' ,folder: '%s', content: %s})"
)

func CreateDomainCy(dm *DomainNode) string {
	if dm.Name == "" {
		log.Fatalln("DomainNode domainName is empty")
		return ""
	}
	return fmt.Sprintf(createDomain, dm.Name)
}

func CreateMicroServerCy(ms *MicroServerNode) string {
	if ms.Name == "" || ms.Domain == "" || ms.Path == "" {
		log.Fatalln("MicroServerNode microServerName or microServerDomain or microServerPath is empty")
		return ""
	}
	return fmt.Sprintf(createMicroServer, ms.Name, ms.Domain, ms.Path)
}

func CreateApiPathCy(api *ApiNode) string {
	if api.Path == "" || api.Type == "" {
		log.Fatalln("ApiNode apiPath or apiType is empty")
		return ""
	}
	return fmt.Sprintf(createApiPath, api.Path, api.Type)
}

func CreatePackageCy(pkg *PackageNode) string {
	if pkg.Name == "" || pkg.Path == "" {
		log.Fatalln("PackageNode pkgName or pkgPath is empty")
		return ""
	}
	return fmt.Sprintf(createPackage, pkg.Name, pkg.Path)
}

func CreateFileCy(file *FileNode) string {
	if file.Name == "" || file.Path == "" {
		log.Fatalln("FileNode fileName or filePath is empty")
		return ""
	}
	return fmt.Sprintf(createFile, file.Name, file.Path)
}

func CreateFunctionCy(node *FunctionNode) string {
	if node.Name == "" || node.File == "" || node.Folder == "" {
		log.Fatalln("FunctionNode funcName or funcFile or folderPath is empty")
		return ""
	}
	return fmt.Sprintf(createFunction, node.Name, node.File, node.Folder, node.Rec, node.Content)
}

func CreateStructCy(st *StructNode) string {
	if st.Name == "" || st.File == "" {
		log.Fatalln("StructNode structName or structFile is empty")
		return ""
	}
	return fmt.Sprintf(createStruct, st.Name, st.File, st.Folder, st.Content)
}

const (
	domainToServer = "MATCH (d:Domain { name: '%s' }) MATCH (s:MicroServer { name: '%s'}) MERGE (d)-[:HAS_SERVER]->(s)"
	serverToPkg    = "MATCH (s:MicroServer { name: '%s' }) MATCH (p:Package { name: '%s',path : '%s'}) MERGE (s)-[:HAS_PACKAGE]->(p)"
	apiToFunction  = "MATCH (a:Api {path: '%s', type: '%s'}) MATCH (f:Function{name:'%s', folder: '%s'}) MERGE (a)-[:MAPPING]->(f)"
	pkgToFile      = "MATCH (p:Package { name: '%s', path: '%s' }) MATCH (f:File { name: '%s', path: '%s' }) MERGE (p)-[:HAS_FILE]->(f)"
	fileToFunction = "MATCH (f1:File { name: '%s', path: '%s' }) MATCH (f2:Function {name:'%s', file:'%s'}) MERGE (f1)-[:HAS_FUNCTION]->(f2)"
	fileToStruct   = "MATCH (f:File { name: '%s', path: '%s' }) MATCH (s:Struct {name:'%s', file:'%s'}) MERGE (f)-[:HAS_STRUCT]->(s)"
	funcCallFunc   = "MATCH (f1:Function {name:'%s', file: '%s', rec : '%s'}) MATCH (f2:Function {name:'%s', file: '%s', rec : '%s'}) MERGE (f1)-[:CALL]->(f2)"
	funcReceiver   = "MATCH (f:Function {name:'%s', file: '%s', rec : '%s' }) MATCH (s:Struct {name:'%s', file: '%s'}) MERGE (f)-[:RECEIVER]->(s)"
	funcParam      = "MATCH (f:Function {name:'%s', file: '%s', rec : '%s' }) MATCH (s:Struct {name:'%s', file: '%s'}) MERGE (f)-[:PARAM]->(s)"
	funcReturn     = "MATCH (f:Function {name:'%s', file: '%s', rec : '%s' }) MATCH (s:Struct {name:'%s', file: '%s'}) MERGE (f)-[:RETURN]->(s)"
)

func DomainToServerCy(link *DomainToServerLink) string {
	if link.Domain.Name == "" || link.Server.Name == "" {
		log.Fatalln("DomainToServerLink domainName or serverName is empty")
		return ""
	}
	return fmt.Sprintf(domainToServer, link.Domain.Name, link.Server.Name)
}

func ServerToPkgCy(link *ServerToPkgLink) string {
	if link.Server.Name == "" || link.Pkg.Name == "" {
		log.Fatalln("ServerToPkgLink serverName or pkgName is empty")
		return ""
	}
	return fmt.Sprintf(serverToPkg, link.Server.Name, link.Pkg.Name, link.Pkg.Path)
}

func ApiToFunctionCy(link *ApiToFuncLink) string {
	if link.Api.Path == "" || link.Api.Type == "" || link.Func.Name == "" || link.Func.Folder == "" {
		log.Fatalln("ApiToFuncLink apiPath or apiType or funcName or folderPath is empty")
		return ""
	}
	return fmt.Sprintf(apiToFunction, link.Api.Path, link.Api.Type, link.Func.Name, link.Func.Folder)
}

func PkgToFileCy(link *PkgToFileLink) string {
	if link.Pkg.Name == "" || link.Pkg.Path == "" || link.File.Name == "" || link.File.Path == "" {
		log.Fatalln("PkgToFileLink pkgName or pkgPath or fileName or filePath is empty")
		return ""
	}
	return fmt.Sprintf(pkgToFile, link.Pkg.Name, link.Pkg.Path, link.File.Name, link.File.Path)
}

func FileToFunctionCy(link *FileToFuncLink) string {
	if link.File.Name == "" || link.File.Path == "" || link.Func.Name == "" || link.Func.File == "" {
		log.Fatalln("FileToFuncLink fileName or filePath or funcName or funcFile is empty")
		return ""
	}
	return fmt.Sprintf(fileToFunction, link.File.Name, link.File.Path, link.Func.Name, link.Func.File)
}

func FileToStructCy(link *FileToStructLink) string {
	if link.File.Name == "" || link.File.Path == "" || link.Struct.Name == "" || link.Struct.File == "" {
		log.Fatalln("pkgToFileLink fileName or filePath or structName or structFile is empty")
		return ""
	}
	return fmt.Sprintf(fileToStruct, link.File.Name, link.File.Path, link.Struct.Name, link.Struct.File)
}

func FuncCallFuncCy(link *FuncCallFuncLink) string {
	if link.Caller.Name == "" || link.Caller.File == "" || link.Callee.Name == "" || link.Callee.File == "" {
		log.Fatalln("FuncCallFuncLink callerName or callerFile or calleeName or calleeFile is empty")
		return ""
	}
	return fmt.Sprintf(funcCallFunc, link.Caller.Name, link.Caller.File, link.Caller.Rec,
		link.Callee.Name, link.Callee.File, link.Callee.Rec)
}

func FuncReceiverCy(link *FuncReceiverLink) string {
	if link.Func.Name == "" || link.Func.File == "" || link.Struct.Name == "" || link.Struct.File == "" {
		log.Fatalln("FuncReceiverLink funcName or funcFile or structName or structFile is empty")
		return ""
	}
	return fmt.Sprintf(funcReceiver, link.Func.Name, link.Func.File, link.Func.Rec, link.Struct.Name, link.Struct.File)
}

func FuncParamCy(link *FuncParamLink) string {
	if link.Func.Name == "" || link.Func.File == "" || link.Struct.Name == "" || link.Struct.File == "" {
		log.Fatalln("FuncParamLink funcName or funcFile or structName or structFile is empty")
		return ""
	}
	return fmt.Sprintf(funcParam, link.Func.Name, link.Func.File, link.Func.Rec, link.Struct.Name, link.Struct.File)
}

func FuncReturnCy(link *FuncReturnLink) string {
	if link.Func.Name == "" || link.Func.File == "" || link.Struct.Name == "" || link.Struct.File == "" {
		log.Fatalln("FuncReturnLink funcName or funcFile or structName or structFile is empty")
		return ""
	}
	return fmt.Sprintf(funcReturn, link.Func.Name, link.Func.File, link.Func.Rec, link.Struct.Name, link.Struct.File)
}
