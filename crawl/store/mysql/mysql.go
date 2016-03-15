package mysql

import (
	"database/sql"
	"flag"
	"strings"

	. "../../base"
	"../../store"
	_ "github.com/go-sql-driver/mysql"
)

var mysql string

func init() {
	flag.StringVar(&mysql, "mysql", "root@/stock", "mysql uri")
	store.Register("mysql", &Mysql{})
}

func (p *Mysql) Open() (err error) {
	if p.db != nil {
		p.Close()
	}

	dsn := mysql
	if !strings.Contains(dsn, "?") {
		dsn = dsn + "?"
	}
	if !strings.Contains(dsn, "parseTime=") {
		dsn = dsn + "&parseTime=true"
	}
	if !strings.Contains(dsn, "loc=") {
		dsn = dsn + "&loc=UTC"
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	p.db = db
	p.SetMaxOpenConns(100)
	return
}

func (p *Mysql) SetMaxOpenConns(n int) {
	num := p.getMaxConnections()
	if n < 0 {
		n = 0
	} else {
		if n > num/2 {
			n = num / 2
		}
		if n < 1 {
			n = 1
		}
	}
	p.db.SetMaxOpenConns(n)
}

type Mysql struct {
	db *sql.DB
}

func (p *Mysql) Close() {
	if p.db != nil {
		p.db.Close()
		p.db = nil
	}
}

func (p *Mysql) getMaxConnections() int {
	num := 0
	p.db.QueryRow("SELECT @@max_connections").Scan(&num)
	return num
}
