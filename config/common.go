package config

type Neo4jConfig struct {
	Url      string
	UserName string
	Password string
}

var Neo4jCfg = &Neo4jConfig{}
