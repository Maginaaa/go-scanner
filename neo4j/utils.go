package neo4j

import (
	"log"
)

func BatchWrite(cys []string) {
	client := NewClient()
	defer client.CloseDriver()
	for _, cy := range cys {
		err := client.Write(cy)
		if err != nil {
			log.Printf("CypherBatchWrite to DB with error: %s, nodeList: cys \n", err)
		}
	}

}

func Write(cypher string) {
	client := NewClient()
	defer client.CloseDriver()
	func(cypher string) {
		err := client.Write(cypher)
		if err != nil {
			log.Printf("CypherWrite to DB with error: %s, cypher: %s \n", err, cypher)
		}
	}(cypher)

}
