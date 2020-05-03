package orm

import (
	"strings"
	"time"

	"go-online/lib/ecode"
	"go-online/lib/log"
	xtime "go-online/lib/time"

	// database driver
	// _ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Config postgreSql config.
type Config struct {
	DSN         string         // data source name.
	Active      int            // pool
	Idle        int            // pool
	IdleTimeout xtime.Duration // connect max life time.
}

type ormLog struct{}

func (l ormLog) Print(v ...interface{}) {
	log.Info(strings.Repeat("%v ", len(v)), v...)
}

func init() {
	gorm.ErrRecordNotFound = ecode.NothingFound
}

// NewPostgreSQL new db and retry connection when has error.
func NewPostgreSQL(c *Config) (db *gorm.DB) {
	db, err := gorm.Open("postgres", c.DSN)
	if err != nil {
		log.Error("db dsn(%s) error: %v", c.DSN, err)
		panic(err)
	}
	db.DB().SetMaxIdleConns(c.Idle)
	db.DB().SetMaxOpenConns(c.Active)
	db.DB().SetConnMaxLifetime(time.Duration(c.IdleTimeout) / time.Second)
	db.SetLogger(ormLog{})
	return
}
