package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	Port     int    `koanf:"port"`
	Host     string `koanf:"host"`
	DBName   string `koanf:"db_name"`
}

type MySQLDB struct {
	config Config
	db     *sql.DB
}

func (m MySQLDB) Conn() *sql.DB{
	return  m.db
}


func New (config Config) *MySQLDB{
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
	config.Username,config.Password,config.Host,config.Port,config.DBName)

	db,err:=sql.Open("mysql",dsn)

	if err!=nil{
		panic(fmt.Errorf("can't open mysql db:%v",err))
	}

	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("can't connect to mysql db: %v", err))
	}

	db.SetConnMaxLifetime(time.Minute*3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &MySQLDB{config:config,db:db}

}