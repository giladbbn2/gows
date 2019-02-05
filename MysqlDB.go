package gows

import "database/sql"

type MysqlConnConfig struct {
	Host   string
	User   string
	Pass   string
	DBName string
	Port   string
	DB     *sql.DB
}
