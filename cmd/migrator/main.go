package main

import (
	"shs/log"
	"shs/mariadb"
)

func main() {
	err := mariadb.Migrate()
	if err != nil {
		log.Fatalln(err)
	}
}
