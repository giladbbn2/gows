package gows

import (
	"database/sql"
	"errors"
	"strconv"
)

type Mysql struct{}

func NewMysql() *Mysql {
	return new(Mysql)
}

func (mysql *Mysql) AddConnection(conn string, host string, user string, pass string, dbname string, port string) {

	mysqlConnections[conn] = &MysqlConnConfig{Host: host, User: user, Pass: pass, DBName: dbname, Port: port, DB: nil}

}

func (mysql *Mysql) OpenConnection(conn string) (*sql.DB, error) {

	mysqldb, ok := mysqlConnections[conn]
	if !ok {
		return nil, errors.New("connection not found")
	}

	db, err := sql.Open("mysql", mysqldb.User+":"+mysqldb.Pass+"@tcp("+mysqldb.Host+":"+mysqldb.Port+")/"+mysqldb.DBName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil

}

func (mysql *Mysql) CloseConnection(conn string) {

	mysqldb, ok := mysqlConnections[conn]
	if !ok {
		return
	}

	if mysqldb.DB == nil {
		return
	}

	err := mysqldb.DB.Close()
	if err != nil {
		return
	}

}

func (mysql *Mysql) AddAndOpenConnection(conn string, host string, user string, pass string, dbname string, port string) (*sql.DB, error) {

	mysql.AddConnection(conn, host, user, pass, dbname, port)

	db, err := mysql.OpenConnection(conn)
	if err != nil {
		return nil, err
	}

	mysqlConnections[conn].DB = db

	return db, nil

}

func (mysql *Mysql) GetConnection(conn string) (*sql.DB, error) {

	mysqldb, ok := mysqlConnections[conn]
	if !ok {
		return nil, errors.New("connection not found")
	}

	return mysqldb.DB, nil

}

func (mysql *Mysql) Query(conn string, sql string) ([][]interface{}, error) {

	var resultsAll = make([][]interface{}, 0)

	db, err := mysql.GetConnection(conn)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultSetId = 0

	for true {

		if resultSetId > 0 {

			// allow branch prediction more easily

			if !rows.NextResultSet() {
				break
			}

		}

		cols, _ := rows.Columns()
		numCols := len(cols)

		for rows.Next() {

			results := make([]interface{}, numCols)
			for j := 0; j < numCols; j++ {
				var tmp interface{}
				results[j] = &tmp
			}

			if err := rows.Scan(results...); err != nil {
				return nil, err
			}

			resultsAll = append(resultsAll, results)

		}

		resultSetId++

	}

	return resultsAll, nil

}

func (mysql *Mysql) ToByteSlice(mysqlResultVal interface{}) ([]byte, error) {
	return (*(mysqlResultVal.(*interface{}))).([]byte), nil
}

func (mysql *Mysql) ToString(mysqlResultVal interface{}) (string, error) {
	return string((*(mysqlResultVal.(*interface{}))).([]byte)), nil
}

func (mysql *Mysql) ToInt(mysqlResultVal interface{}) (int, error) {
	return strconv.Atoi(string((*(mysqlResultVal.(*interface{}))).([]byte)))
}
