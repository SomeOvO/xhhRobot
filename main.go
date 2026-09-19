package main

import (
	"heybox/database"
	"heybox/database/sqlite"
	"heybox/log"
)

func main() {
	log.Mod = "debug"
	log.Init()
	sqlite.Init()
	database.Db.Conn()
}
