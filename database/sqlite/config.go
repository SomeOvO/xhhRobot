package sqlite

import (
	"fmt"
	"heybox/database"
	"heybox/log"
)

const config_db = "cfg"

func config_set(Name int, Value string, Enable bool) error {
	query := fmt.Sprintf("INSERT INTO %s (name,value,enable) VALUES (?,?,?) ON CONFLICT(name) DO UPDATE SET name=excluded.name,value=excluded.value,enable=excluded.enable", config_db)
	_, err := db.Exec(query, Name, Value, Enable)
	if err != nil {
		log.ErrorLOG("数据库写入出错", err)
		return nil
	}
	return nil
}

func config_get(ConfigName int) *database.TYPE_Config {
	query := fmt.Sprintf("SELECT value,enable FROM %s WHERE name=?", config_db)
	row := db.QueryRow(query, ConfigName)
	var config database.TYPE_Config
	row.Scan(&config.Value, &config.Enable)
	config.Name = ConfigName
	return &config
}
