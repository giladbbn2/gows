package gows

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"
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

func (mysql *Mysql) Query(conn string, sql string, args ...interface{}) ([][]interface{}, int, error) {

	if conn == "" || sql == "" {
		return nil, 0, nil
	}

	db, err := mysql.GetConnection(conn)
	if err != nil {
		return nil, 0, err
	}

	if sql[0] == " "[0] {
		sql = strings.TrimLeft(sql, " ")
	}

	if len(sql) < 3 {
		return nil, 0, nil
	}

	if strings.ToLower(sql[:3]) == "sel" {

		resultSet := make([][]interface{}, 0)

		rows, err := db.Query(sql, args...)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()

		cols, _ := rows.Columns()
		numCols := len(cols)

		for rows.Next() {

			results := make([]interface{}, numCols)
			for j := 0; j < numCols; j++ {
				var tmp interface{}
				results[j] = &tmp
			}

			if err := rows.Scan(results...); err != nil {
				return nil, 0, err
			}

			resultSet = append(resultSet, results)

		}

		return resultSet, len(resultSet), nil

	}

	res, err := db.Exec(sql, args...)
	if err != nil {
		return nil, 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return nil, 0, err
	}

	return nil, int(count), nil

}

func (mysql *Mysql) Exec(conn string, sql string, args ...interface{}) (int, error) {
	_, count, err := mysql.Query(conn, sql, args...)
	return count, err
}

func (mysql *Mysql) ToByteSlice(mysqlResultVal interface{}) ([]byte, bool, error) {

	i := (*(mysqlResultVal.(*interface{})))
	if i == nil {
		return nil, true, nil
	}

	by, ok := i.([]byte)
	if !ok {
		return nil, false, errors.New("can't convert to []byte")
	}

	return by, false, nil

}

func (mysql *Mysql) ToBool(mysqlResultVal interface{}) (bool, bool, error) {

	i := (*(mysqlResultVal.(*interface{})))
	if i == nil {
		return false, true, nil
	}

	by, ok := i.([]byte)
	if !ok || len(by) == 0 {
		return false, false, errors.New("can't convert to bool")
	}

	return by[0] == 1, false, nil

}

func (mysql *Mysql) ToString(mysqlResultVal interface{}) (string, bool, error) {

	i := (*(mysqlResultVal.(*interface{})))
	if i == nil {
		return "", true, nil
	}

	by, ok := i.([]byte)
	if !ok {
		return "", false, errors.New("can't convert to string")
	}

	s := string(by)

	return s, false, nil

}

func (mysql *Mysql) ToInt(mysqlResultVal interface{}) (int, bool, error) {

	i := (*(mysqlResultVal.(*interface{})))
	if i == nil {
		return 0, true, nil
	}

	by, ok := i.([]byte)
	if !ok {
		return 0, false, errors.New("can't convert to int")
	}

	s := string(by)

	// Atoi can be faster than ParseInt
	in, err := strconv.Atoi(s)

	return in, false, err

}

func (mysql *Mysql) ToInt8(mysqlResultVal interface{}) (int8, bool, error) {

	i := (*(mysqlResultVal.(*interface{})))
	if i == nil {
		return 0, true, nil
	}

	by, ok := i.([]byte)
	if !ok {
		return 0, false, errors.New("can't convert to int8")
	}

	s := string(by)

	in, err := strconv.ParseInt(s, 10, 8)

	return int8(in), false, err

}

func (mysql *Mysql) ToInt16(mysqlResultVal interface{}) (int16, bool, error) {

	i := (*(mysqlResultVal.(*interface{})))
	if i == nil {
		return 0, true, nil
	}

	by, ok := i.([]byte)
	if !ok {
		return 0, false, errors.New("can't convert to int16")
	}

	s := string(by)

	in, err := strconv.ParseInt(s, 10, 16)

	return int16(in), false, err

}

func (mysql *Mysql) ToInt32(mysqlResultVal interface{}) (int32, bool, error) {

	i := (*(mysqlResultVal.(*interface{})))
	if i == nil {
		return 0, true, nil
	}

	by, ok := i.([]byte)
	if !ok {
		return 0, false, errors.New("can't convert to int32")
	}

	s := string(by)

	in, err := strconv.ParseInt(s, 10, 32)

	return int32(in), false, err

}

func (mysql *Mysql) ToInt64(mysqlResultVal interface{}) (int64, bool, error) {

	i := (*(mysqlResultVal.(*interface{})))
	if i == nil {
		return 0, true, nil
	}

	by, ok := i.([]byte)
	if !ok {
		return 0, false, errors.New("can't convert to int64")
	}

	s := string(by)

	in, err := strconv.ParseInt(s, 10, 64)

	return in, false, err

}

func (mysql *Mysql) ToUint(mysqlResultVal interface{}) (uint, bool, error) {

	i := (*(mysqlResultVal.(*interface{})))
	if i == nil {
		return 0, true, nil
	}

	by, ok := i.([]byte)
	if !ok {
		return 0, false, errors.New("can't convert to uint")
	}

	s := string(by)

	in, err := strconv.ParseUint(s, 10, 0)

	return uint(in), false, err

}

func (mysql *Mysql) ToUint8(mysqlResultVal interface{}) (uint8, bool, error) {

	i := (*(mysqlResultVal.(*interface{})))
	if i == nil {
		return 0, true, nil
	}

	by, ok := i.([]byte)
	if !ok {
		return 0, false, errors.New("can't convert to uint8")
	}

	s := string(by)

	in, err := strconv.ParseUint(s, 10, 8)

	return uint8(in), false, err

}

func (mysql *Mysql) ToUint16(mysqlResultVal interface{}) (uint16, bool, error) {

	i := (*(mysqlResultVal.(*interface{})))
	if i == nil {
		return 0, true, nil
	}

	by, ok := i.([]byte)
	if !ok {
		return 0, false, errors.New("can't convert to uint16")
	}

	s := string(by)

	in, err := strconv.ParseUint(s, 10, 16)

	return uint16(in), false, err

}

func (mysql *Mysql) ToUint32(mysqlResultVal interface{}) (uint32, bool, error) {

	i := (*(mysqlResultVal.(*interface{})))
	if i == nil {
		return 0, true, nil
	}

	by, ok := i.([]byte)
	if !ok {
		return 0, false, errors.New("can't convert to uint32")
	}

	s := string(by)

	in, err := strconv.ParseUint(s, 10, 32)

	return uint32(in), false, err

}

func (mysql *Mysql) ToUint64(mysqlResultVal interface{}) (uint64, bool, error) {

	i := (*(mysqlResultVal.(*interface{})))
	if i == nil {
		return 0, true, nil
	}

	by, ok := i.([]byte)
	if !ok {
		return 0, false, errors.New("can't convert to uint64")
	}

	s := string(by)

	in, err := strconv.ParseUint(s, 10, 64)

	return in, false, err

}

func (mysql *Mysql) ToFloat32(mysqlResultVal interface{}) (float32, bool, error) {

	i := (*(mysqlResultVal.(*interface{})))
	if i == nil {
		return 0, true, nil
	}

	by, ok := i.([]byte)
	if !ok {
		return 0, false, errors.New("can't convert to float32")
	}

	s := string(by)

	fl, err := strconv.ParseFloat(s, 32)

	return float32(fl), false, err

}

func (mysql *Mysql) ToFloat64(mysqlResultVal interface{}) (float64, bool, error) {

	i := (*(mysqlResultVal.(*interface{})))
	if i == nil {
		return 0, true, nil
	}

	by, ok := i.([]byte)
	if !ok {
		return 0, false, errors.New("can't convert to float64")
	}

	s := string(by)

	fl, err := strconv.ParseFloat(s, 64)

	return fl, false, err

}

func (mysql *Mysql) ToTime(mysqlResultVal interface{}) (time.Time, bool, error) {

	i := (*(mysqlResultVal.(*interface{})))
	if i == nil {
		return time.Unix(0, 0), true, nil
	}

	by, ok := i.([]byte)
	if !ok {
		return time.Unix(0, 0), false, errors.New("can't convert to Time")
	}

	s := string(by)

	if len(s) < 19 {
		return time.Unix(0, 0), false, errors.New("can't convert to Time")
	}

	// effectively make it UTC
	t, err := time.Parse(time.RFC3339, s[:10]+"T"+s[11:]+"Z")

	return t, false, err

}
