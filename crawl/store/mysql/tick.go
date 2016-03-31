package mysql

import (
	"database/sql"
	"strings"
	"time"

	. "../../base"

	"github.com/golang/glog"
)

func (p *Mysql) createTickTable(table string) {
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

func (p *Mysql) LoadTicks(table string) (res []Tick, err error) {
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

func (p *Mysql) checkTicks(table string, ticks []Tick) ([]Tick, error) {
	var stmt *sql.Stmt
	var err error
	for i := 0; i < 2; i++ {
		stmt, err = p.db.Prepare("SELECT 1 FROM `" + table + "` WHERE `time`=?")
		if err == nil {
			break
		}
		if strings.Index(err.Error(), "Error 1146:") > -1 {
			p.createTickTable(table)
			continue
		}
	}
	if err != nil {
		return ticks, err
	}
	defer stmt.Close()

	unsave := []Tick{}
	for i, c := 0, len(ticks); i < c; i++ {
		has := false
		stmt.QueryRow(ticks[i].Time).Scan(&has)
		if !has {
			unsave = append(unsave, ticks[i])
		}
	}
	return unsave, nil
}

func (p *Mysql) saveTicks(table string, ticks []Tick) error {
	var stmt *sql.Stmt
	var err error
	for i := 0; i < 2; i++ {
		stmt, err = p.db.Prepare("INSERT INTO `" + table + "`(`time`,`price`,`change`,`volume`,`turnover`,`type`) values(?,?,?,?,?,?)")
		if err == nil {
			break
		}
		if strings.Index(err.Error(), "Error 1146:") > -1 {
			p.createTickTable(table)
			continue
		}
	}
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i, c := 0, len(ticks); i < c; i++ {
		tick := &ticks[i]
		_, e := stmt.Exec(tick.Time, tick.Price, tick.Change, tick.Volume, tick.Turnover, tick.Type)
		if e != nil {
			msg := e.Error()
			if strings.Index(msg, "Error 1062:") > -1 {
				// duplicate
				continue
			}
			if strings.Index(msg, "Error 1615:") > -1 {
				// Prepared statement needs to be re-prepared
				i--
				continue
			}
			if strings.Index(msg, "Error 1461:") > -1 {
				i--
				time.Sleep(time.Millisecond * 100)
				continue
			}
			err = e
		}
	}
	if err != nil {
		glog.Warningf("insert tick error %v", err)
	}
	return err
}

func (p *Mysql) SaveTicks(table string, ticks []Tick) error {
	var err error
	unsave := ticks
	saved := 0
	for i := 0; i < 10 && len(unsave) > 0; i++ {
		err = p.saveTicks(table, unsave)
		sum := len(unsave)
		unsave, err = p.checkTicks(table, unsave)
		saved = sum - len(unsave)
		if saved > 0 {
			i--
		}
	}

	return err
}
func (p *Mysql) HasTickData(table string, t time.Time) bool {
	has := false
	t = t.Truncate(time.Hour * 24)
	end := t.AddDate(0, 0, 1)
	p.db.QueryRow("SELECT 1 FROM `"+table+"` WHERE `time` BETWEEN ? AND ? LIMIT 0,1", t, end).Scan(&has)
	return has
}
