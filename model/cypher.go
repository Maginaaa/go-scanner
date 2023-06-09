package model

import (
	"fmt"
	"log"
)

const (
	createDomain      = "MERGE (:Domain { name: '%s'})"
	createMicroServer = "MERGE (:MicroServer { name: '%s', domain: '%s', path: '%s' })"
	createApiPath     = "MERGE (:Api {path: '%s', type: '%s'})"
	createGrpcApi     = "MERGE (:GrpcApi {path: '%s', server: '%s'})"
	createPackage     = "MERGE (:Package { name: '%s', path: '%s' , import_path: '%s'})"
	createFile        = "MERGE (:File { name: '%s', path: '%s' })"
	createFunction    = "MERGE (:Function {name:'%s', file: '%s', package: '%s', rec : '%s', content : %s})"
	createStruct      = "MERGE (:Struct {name:'%s', file: '%s' ,folder: '%s', content: %s})"
)

func CreateDomainCy(dm *DomainNode) string {
	if dm.Name == "" {
		log.Println("DomainNode domainName is empty")
		return ""
	}
	return fmt.Sprintf(createDomain, dm.Name)
}

func CreateMicroServerCy(ms *MicroServerNode) string {
	if ms.Name == "" || ms.Domain == "" || ms.Path == "" {
		log.Println("MicroServerNode microServerName or microServerDomain or microServerPath is empty")
		return ""
	}
	return fmt.Sprintf(createMicroServer, ms.Name, ms.Domain, ms.Path)
}

func CreateApiPathCy(api *ApiNode) string {
	if api.Path == "" || api.Type == "" {
		log.Println("ApiNode apiPath or apiType is empty")
		return ""
	}
	return fmt.Sprintf(createApiPath, api.Path, api.Type)
}

func CreateGrpcApiPathCy(api *GrpcApiNode) string {
	if api.Path == "" || api.Server == "" {
		log.Println("GrpcApiNode apiPath or apiServer is empty")
		return ""
	}
	return fmt.Sprintf(createGrpcApi, api.Path, api.Server)
}

func CreatePackageCy(pkg *PackageNode) string {
	if pkg.Name == "" || pkg.Path == "" {
		log.Println("PackageNode pkgName or pkgPath is empty")
		return ""
	}
	return fmt.Sprintf(createPackage, pkg.Name, pkg.Path, pkg.ImportPath)
}

func CreateFileCy(file *FileNode) string {
	if file.Name == "" || file.Path == "" {
		log.Println("FileNode fileName or filePath is empty")
		return ""
	}
	return fmt.Sprintf(createFile, file.Name, file.Path)
}

func CreateFunctionCy(node *FunctionNode) string {
	if node.Name == "" || node.File == "" || node.Package == "" {
		log.Print("FunctionNode functionName or functionFile or functionPackage is empty")
		return ""
	}
	return fmt.Sprintf(createFunction, node.Name, node.File, node.Package, node.Rec, node.Content)
}

func CreateStructCy(st *StructNode) string {
	if st.Name == "" || st.File == "" {
		log.Println("StructNode structName or structFile is empty")
		return ""
	}
	return fmt.Sprintf(createStruct, st.Name, st.File, st.Folder, st.Content)
}

const (
	domainToServer      = "MATCH (d:Domain { name: '%s' }) MATCH (s:MicroServer { name: '%s'}) MERGE (d)-[:HAS_SERVER]->(s)"
	serverToPkg         = "MATCH (s:MicroServer { name: '%s' }) MATCH (p:Package { name: '%s',path : '%s'}) MERGE (s)-[:HAS_PACKAGE]->(p)"
	apiImplFunction     = "MATCH (a:Api {path: '%s', type: '%s'}) MATCH (f:Function{name:'%s', package: '%s'}) MERGE (a)-[:MAPPING]->(f)"
	apiRequestFunction  = "MATCH (a:Api {path: '%s', type: '%s'}) MATCH (f:Function{name:'%s', package: '%s'}) MERGE (f)-[:REQUEST]->(a)"
	grpcApiImplFunction = "MATCH (a:GrpcApi {path: '%s', server: '%s'}) MATCH (f:Function{name:'%s', package: '%s'}) MERGE (a)-[:GRPC_MAPPING]->(f)"
	pkgToFile           = "MATCH (p:Package { name: '%s', path: '%s' }) MATCH (f:File { name: '%s', path: '%s' }) MERGE (p)-[:HAS_FILE]->(f)"
	fileToFunction      = "MATCH (f1:File { name: '%s', path: '%s' }) MATCH (f2:Function {name:'%s', file:'%s'}) MERGE (f1)-[:HAS_FUNCTION]->(f2)"
	fileToStruct        = "MATCH (f:File { name: '%s', path: '%s' }) MATCH (s:Struct {name:'%s', file:'%s'}) MERGE (f)-[:HAS_STRUCT]->(s)"
	funcCallFunc        = "MATCH (f1:Function {name:'%s', file: '%s', rec : '%s'}) MATCH (f2:Function {name:'%s', file: '%s', rec : '%s'}) MERGE (f1)-[:CALL]->(f2)"
	funcReceiver        = "MATCH (f:Function {name:'%s', file: '%s', rec : '%s' }) MATCH (s:Struct {name:'%s', file: '%s'}) MERGE (f)-[:RECEIVER]->(s)"
	funcParam           = "MATCH (f:Function {name:'%s', file: '%s', rec : '%s' }) MATCH (s:Struct {name:'%s', file: '%s'}) MERGE (f)-[:PARAM]->(s)"
	funcReturn          = "MATCH (f:Function {name:'%s', file: '%s', rec : '%s' }) MATCH (s:Struct {name:'%s', file: '%s'}) MERGE (f)-[:RETURN]->(s)"
)

func DomainToServerCy(link *DomainToServerLink) string {
	if link.Domain.Name == "" || link.Server.Name == "" {
		log.Println("DomainToServerLink domainName or serverName is empty")
		return ""
	}
	return fmt.Sprintf(domainToServer, link.Domain.Name, link.Server.Name)
}

func ServerToPkgCy(link *ServerToPkgLink) string {
	if link.Server.Name == "" || link.Pkg.Name == "" {
		log.Println("ServerToPkgLink serverName or pkgName is empty")
		return ""
	}
	return fmt.Sprintf(serverToPkg, link.Server.Name, link.Pkg.Name, link.Pkg.Path)
}

func ApiImplFunctionCy(link *ApiImplFuncLink) string {
	if link.Api.Path == "" || link.Api.Type == "" || link.Func.Name == "" || link.Func.Package == "" {
		log.Println("ApiImplFunctionCy apiPath or apiType or funcName or funcPackage is empty")
		return ""
	}
	return fmt.Sprintf(apiImplFunction, link.Api.Path, link.Api.Type, link.Func.Name, link.Func.Package)
}

func ApiRequestFunctionCy(link *ApiRequestFuncLink) string {
	if link.Api.Path == "" || link.Api.Type == "" || link.Func.Name == "" || link.Func.Package == "" {
		log.Println("ApiRequestFuncLink apiPath or apiType or funcName or funcPackage is empty")
		return ""
	}
	return fmt.Sprintf(apiRequestFunction, link.Api.Path, link.Api.Type, link.Func.Name, link.Func.Package)
}

func GrpcApiImplFunctionCy(link *GrpcApiImplFuncLink) string {
	if link.GrpcApi.Path == "" || link.GrpcApi.Server == "" || link.Func.Name == "" || link.Func.Package == "" {
		log.Println("GrpcApiImplFunctionCy apiPath or apiServer or funcName or funcPackage is empty")
		return ""
	}
	return fmt.Sprintf(grpcApiImplFunction, link.GrpcApi.Path, link.GrpcApi.Server, link.Func.Name, link.Func.Package)
}

func PkgToFileCy(link *PkgToFileLink) string {
	if link.Pkg.Name == "" || link.Pkg.Path == "" || link.File.Name == "" || link.File.Path == "" {
		log.Println("PkgToFileLink pkgName or pkgPath or fileName or filePath is empty")
		return ""
	}
	return fmt.Sprintf(pkgToFile, link.Pkg.Name, link.Pkg.Path, link.File.Name, link.File.Path)
}

func FileToFunctionCy(link *FileToFuncLink) string {
	if link.File.Name == "" || link.File.Path == "" || link.Func.Name == "" || link.Func.File == "" {
		log.Println("FileToFuncLink fileName or filePath or funcName or funcFile is empty")
		return ""
	}
	return fmt.Sprintf(fileToFunction, link.File.Name, link.File.Path, link.Func.Name, link.Func.File)
}

func FileToStructCy(link *FileToStructLink) string {
	if link.File.Name == "" || link.File.Path == "" || link.Struct.Name == "" || link.Struct.File == "" {
		log.Println("pkgToFileLink fileName or filePath or structName or structFile is empty")
		return ""
	}
	return fmt.Sprintf(fileToStruct, link.File.Name, link.File.Path, link.Struct.Name, link.Struct.File)
}

func FuncCallFuncCy(link *FuncCallFuncLink) string {
	if link.Caller.Name == "" || link.Caller.File == "" || link.Callee.Name == "" || link.Callee.File == "" {
		log.Println("FuncCallFuncLink callerName or callerFile or calleeName or calleeFile is empty")
		return ""
	}
	return fmt.Sprintf(funcCallFunc, link.Caller.Name, link.Caller.File, link.Caller.Rec,
		link.Callee.Name, link.Callee.File, link.Callee.Rec)
}

func FuncReceiverCy(link *FuncReceiverLink) string {
	if link.Func.Name == "" || link.Func.File == "" || link.Struct.Name == "" || link.Struct.File == "" {
		log.Println("FuncReceiverLink funcName or funcFile or structName or structFile is empty")
		return ""
	}
	return fmt.Sprintf(funcReceiver, link.Func.Name, link.Func.File, link.Func.Rec, link.Struct.Name, link.Struct.File)
}

func FuncParamCy(link *FuncParamLink) string {
	if link.Func.Name == "" || link.Func.File == "" || link.Struct.Name == "" || link.Struct.File == "" {
		log.Println("FuncParamLink funcName or funcFile or structName or structFile is empty")
		return ""
	}
	return fmt.Sprintf(funcParam, link.Func.Name, link.Func.File, link.Func.Rec, link.Struct.Name, link.Struct.File)
}

func FuncReturnCy(link *FuncReturnLink) string {
	if link.Func.Name == "" || link.Func.File == "" || link.Struct.Name == "" || link.Struct.File == "" {
		log.Println("FuncReturnLink funcName or funcFile or structName or structFile is empty")
		return ""
	}
	return fmt.Sprintf(funcReturn, link.Func.Name, link.Func.File, link.Func.Rec, link.Struct.Name, link.Struct.File)
}
