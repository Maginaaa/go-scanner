# go-scanner

目前仅支持go1.18及以下的项目扫描

基础使用方式如下：
`ProjectName`填写项目名称,`RootPath`填写项目存放路径(不含项目文件)
例如： `/etc/workspace/go-scanner`
```go
import "github.com/Maginaaa/go-scanner/scanner"

s := scanner.Scanner{
    ProjectName:      "go-scanner", // 项目名，同等于go.mod内的module
    RootPath:         "/etc/workspace",
}
s.ServerScanner()
```
在`s.NodeCollection`与`s.LinkCollection`中，有完整的节点信息与调用信息，工具内置了neo4j的写入方法，如需使用
```go
import "github.com/Maginaaa/go-scanner/neo4j"

neo4j.BatchWrite(s.NodeCollection.ToCypherList())
neo4j.BatchWrite(s.LinkCollection.ToCypherList())
```


如果文件名与项目名不一致，需要同时传入ProjectName与ProjectPath， 例如：
```go
s := scanner.Scanner{
    ProjectName:      "github.com/Maginaaa/go-scanner",
    ProjectPath:      "go-scanner",    // 文件名
    RootPath:         "/etc/workspace",
}
s.ServerScanner()
```

如果是微服务项目，需要传入服务路径，
例如当前项目有两个微服务,分别为`/etc/workspace/go-scanner/serverA`与`/etc/workspace/go-scanner/serverB`
扫描服务A的方式如下：
```go
s := scanner.Scanner{
    ProjectName:      "go-scanner",    // 文件名
    MicroServerPath:  "serverA"  // 微服务路径
    RootPath:         "/etc/workspace",
}
s.ServerScanner()
```