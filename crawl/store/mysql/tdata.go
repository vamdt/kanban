package mysql

import (
	"database/sql"
	"strings"
	"time"

	. "../../base"

	"github.com/golang/glog"
)

func (p *Mysql) createDayTdataTable(table string) {
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

func (p *Mysql) LoadTDatas(table string, start time.Time) (res []Tdata, err error) {
	rows, err := p.db.Query("SELECT `time`,`open`,`high`,`low`,`close`,`volume` FROM `"+table+"` WHERE `time`>=? ORDER BY time", start)
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

func (p *Mysql) checkTDatas(table string, datas []Tdata) ([]Tdata, error) {
	var stmt *sql.Stmt
	var err error
	for i := 0; i < 2; i++ {
		stmt, err = p.db.Prepare("SELECT 1 FROM `" + table + "` WHERE `time`=?")
		if err == nil {
			break
		}
		if strings.Index(err.Error(), "Error 1146:") > -1 {
			p.createDayTdataTable(table)
			continue
		}
	}
	if err != nil {
		return datas, err
	}
	defer stmt.Close()

	unsave := []Tdata{}
	for i, c := 0, len(datas); i < c; i++ {
		has := false
		stmt.QueryRow(datas[i].Time).Scan(&has)
		if !has {
			unsave = append(unsave, datas[i])
		}
	}

	return unsave, nil
}

func (p *Mysql) saveTDatas(table string, datas []Tdata) error {
	var stmt *sql.Stmt
	var err error
	for i := 0; i < 2; i++ {
		stmt, err = p.db.Prepare("INSERT INTO `" + table + "`(`time`,`open`,`high`,`low`,`close`,`volume`) values(?,?,?,?,?,?)")
		if err == nil {
			break
		}
		if strings.Index(err.Error(), "Error 1146:") > -1 {
			p.createDayTdataTable(table)
			continue
		}
	}
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i, c := 0, len(datas); i < c; i++ {
		data := &datas[i]
		_, e := stmt.Exec(data.Time, data.Open, data.High, data.Low, data.Close, data.Volume)
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
		glog.Warningln("insert tdata error", err)
	}
	return err
}

func (p *Mysql) SaveTDatas(table string, datas []Tdata) error {
	var err error
	unsave := datas
	saved := 0
	for i := 0; i < 10 && len(unsave) > 0; i++ {
		err = p.saveTDatas(table, unsave)
		sum := len(unsave)
		unsave, err = p.checkTDatas(table, unsave)
		saved = sum - len(unsave)
		if saved > 0 {
			i--
		}
	}

	return err
}
