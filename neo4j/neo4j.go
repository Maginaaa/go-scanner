package neo4j

import (
	"fmt"
	"github.com/Maginaaa/go-scanner/config"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
)

type Neo4j struct {
	Url      string
	UserName string
	PassWord string
	driver   neo4j.Driver
	session  neo4j.Session
}

func NewClient() *Neo4j {
	var client Neo4j
	client.DefaultConnect()
	return &client
}

func (n *Neo4j) DefaultConnect() {

	driver, err := neo4j.NewDriver(config.Neo4jCfg.Url, neo4j.BasicAuth(config.Neo4jCfg.UserName, config.Neo4jCfg.Password, ""))
	if driver == nil {
		return
	}
	n.driver = driver

	n.session = driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		err = session.Close()
		if err != nil {
			log.Println("neo4j session close error:", err)
		}
	}(n.session)
}

func (n *Neo4j) Write(cypher string) error {
	if cypher == "" {
		return fmt.Errorf("cypher is empty")
	}
	_, err := n.session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(cypher, nil)
		if err != nil {
			log.Println("wirte to neo4j DB with error:", err)
			return nil, err
		}
		return result.Consume()
	})
	return err
}

func (n *Neo4j) BatchWrite(cypherList []string) error {
	if len(cypherList) == 0 {
		return fmt.Errorf("cypherList is empty")
	}
	_, err := n.session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		for _, cypher := range cypherList {
			_, err := tx.Run(cypher, nil)
			if err != nil {
				log.Println("wirte to neo4j DB with error:", err)
				return nil, err
			}
		}
		return nil, nil
	})
	return err
}

func (n *Neo4j) CloseDriver() {
	err := n.driver.Close()
	if err != nil {
		log.Println("driver close err :", err)
		return
	}
	return
}
