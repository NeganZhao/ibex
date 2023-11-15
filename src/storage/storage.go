package storage

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/ulricqin/ibex/src/pkg/ormx"
)

type Config struct {
	Gorm     Gorm
	MySQL    MySQL
	Postgres Postgres
}

type Gorm struct {
	Debug             bool
	DBType            string
	MaxLifetime       int
	MaxOpenConns      int
	MaxIdleConns      int
	TablePrefix       string
	EnableAutoMigrate bool
}

type MySQL struct {
	Address    string
	User       string
	Password   string
	DBName     string
	Parameters string
}

func (a MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		a.User, a.Password, a.Address, a.DBName, a.Parameters)
}

type Postgres struct {
	Address  string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Schema   string
}

func (a Postgres) DSN() string {
	arr := strings.Split(a.Address, ":")
	if len(arr) != 2 {
		panic("pg address(" + a.Address + ") invalid")
	}

	if a.Schema != "" {
		return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s search_path=%s",
			arr[0], arr[1], a.User, a.DBName, a.Password, a.SSLMode, a.Schema)
	}

	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		arr[0], arr[1], a.User, a.DBName, a.Password, a.SSLMode)
}

var DB *gorm.DB

func InitDB(c Config) error {
	db, err := newGormDB(c)
	if err == nil {
		DB = db
	}

	return err
}

func newGormDB(c Config) (*gorm.DB, error) {
	var dsn string
	switch c.Gorm.DBType {
	case "mysql":
		dsn = c.MySQL.DSN()
	case "postgres":
		dsn = c.Postgres.DSN()
	default:
		return nil, errors.New("unknown DBType")
	}

	return ormx.New(ormx.Config{
		Debug:        c.Gorm.Debug,
		DBType:       c.Gorm.DBType,
		DSN:          dsn,
		MaxIdleConns: c.Gorm.MaxIdleConns,
		MaxLifetime:  c.Gorm.MaxLifetime,
		MaxOpenConns: c.Gorm.MaxOpenConns,
		TablePrefix:  c.Gorm.TablePrefix,
	})
}
