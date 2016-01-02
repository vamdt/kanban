package crawl

import (
	"database/sql"
	"flag"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/golang/glog"
)

var mysql string

func init() {
	flag.StringVar(&mysql, "mysql", "root@/stock", "mysql uri")
}

func NewMysqlStore() (ms *MysqlStore, err error) {
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
	ms = &MysqlStore{db: db}
	return
}

type MysqlStore struct {
	db *sql.DB
}

func (p *MysqlStore) Close() {
	if p.db != nil {
		p.db.Close()
	}
}

func (p *MysqlStore) createDayTdataTable(table string) {
	sql := "CREATE TABLE IF NOT EXISTS `" + table + "` (" +
		"`time` DATETIME NOT NULL," +
		"`open` INT(11) NOT NULL DEFAULT 0," +
		"`high` INT(11) NOT NULL DEFAULT 0," +
		"`low` INT(11) NOT NULL DEFAULT 0," +
		"`close` INT(11) NOT NULL DEFAULT 0," +
		"`volume` INT(11) NOT NULL DEFAULT 0," +
		"PRIMARY KEY (`time`)" +
		")"
	_, err := p.db.Exec(sql)
	if err != nil {
		glog.Warningln("create table err", table, err)
	}
}

func (p *MysqlStore) LoadTDatas(table string) (res []Tdata, err error) {
	rows, err := p.db.Query("SELECT `time`,`open`,`high`,`low`,`close`,`volume` FROM `" + table + "` ORDER BY time")
	if err != nil {
		glog.Warningln(err)
		return
	}
	defer rows.Close()
	d := Tdata{}
	for rows.Next() {
		if err := rows.Scan(&d.Time, &d.Open, &d.High, &d.Low, &d.Close, &d.Volume); err != nil {
			glog.Warningln(err)
		}
		res = append(res, d)
	}
	if err := rows.Err(); err != nil {
		glog.Warningln(err)
	}
	return
}

func (p *MysqlStore) SaveTData(table string, data *Tdata) (err error) {
	for i := 0; i < 2; i++ {
		_, err = p.db.Exec("INSERT INTO `"+table+"`(`time`,`open`,`high`,`low`,`close`,`volume`) values(?,?,?,?,?,?)",
			data.Time, data.Open, data.High, data.Low, data.Close, data.Volume)
		if err == nil {
			break
		}
		p.createDayTdataTable(table)
	}
	if err != nil {
		glog.Warningln("insert tdata error", err, *data)
	}
	return
}

func (p *MysqlStore) createTickTable(table string) {
	sql := "CREATE TABLE IF NOT EXISTS `" + table + "` (" +
		"`time` DATETIME(3) NOT NULL," +
		"`price` INT(11) NOT NULL DEFAULT 0," +
		"`change` INT(11) NOT NULL DEFAULT 0," +
		"`volume` INT(11) NOT NULL DEFAULT 0," +
		"`turnover` INT(11) NOT NULL DEFAULT 0," +
		"`type` INT(11) NOT NULL DEFAULT 0," +
		"PRIMARY KEY (`time`)" +
		")"
	_, err := p.db.Exec(sql)
	if err != nil {
		glog.Warningln("create table err", table, err)
	}
}

func (p *MysqlStore) LoadTicks(table string) (res []Tick, err error) {
	d := Tick{}
	rows, err := p.db.Query("SELECT `time`,`price`,`change`,`volume`,`turnover`,`type` FROM `" + table + "` ORDER BY time")
	if err != nil {
		glog.Warningln(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&d.Time, &d.Price, &d.Change, &d.Volume, &d.Turnover, &d.Type); err != nil {
			glog.Warningln(err)
		}
		res = append(res, d)
	}
	if err := rows.Err(); err != nil {
		glog.Warningln(err)
	}
	return
}

func (p *MysqlStore) SaveTick(table string, tick *Tick) (err error) {
	for i := 0; i < 2; i++ {
		_, err = p.db.Exec("INSERT INTO `"+table+"`(`time`,`price`,`change`,`volume`,`turnover`,`type`) values(?,?,?,?,?,?)",
			tick.Time, tick.Price, tick.Change, tick.Volume, tick.Turnover, tick.Type)
		if err == nil {
			break
		}
		p.createTickTable(table)
	}
	if err != nil {
		glog.Warningf("insert tick error %+v %+v", err, *tick)
	}
	return
}
