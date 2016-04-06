package mysql

import (
	"database/sql"
	"strings"
	"time"

	. "../../base"

	"github.com/golang/glog"
)

const (
	macdTable = "macd"
)

func (p *Mysql) createMacdTable() {
	table := macdTable
	sql := "CREATE TABLE IF NOT EXISTS `" + table + "` (" +
		"`code` VARCHAR(16) NOT NULL DEFAULT ''," +
		"`time` DATETIME NOT NULL," +
		"`typ` INT(11) NOT NULL DEFAULT 0," +
		"`emas` INT(11) NOT NULL DEFAULT 0," +
		"`emal` INT(11) NOT NULL DEFAULT 0," +
		"`dea` INT(11) NOT NULL DEFAULT 0," +
		"PRIMARY KEY (`code`, `time`, `typ`)" +
		")"
	_, err := p.db.Exec(sql)
	if err != nil {
		glog.Warningln("create table err", table, err)
	}
}

func (p *Mysql) LoadMacd(symbol string, typ int, start time.Time) (*Tdata, error) {
	table := macdTable
	d := Tdata{}
	sql := "SELECT `time`, `emas`,`emal`,`dea` FROM `" + table + "` WHERE `code`=? AND `typ`=? AND `time`>=? ORDER BY `time` ASC LIMIT 0,1"
	err := p.db.QueryRow(sql, symbol, typ, start).Scan(&d.Time, &d.Emas, &d.Emal, &d.DEA)
	if err != nil {
		glog.Warningln(err)
		return nil, err
	}
	return &d, nil
}

func (p *Mysql) checkMacds(symbol string, typ int, datas []Tdata) ([]Tdata, error) {
	table := macdTable
	var stmt *sql.Stmt
	var err error
	for i := 0; i < 2; i++ {
		stmt, err = p.db.Prepare("SELECT 1 FROM `" + table + "` WHERE `time`=? AND `code`=? AND `typ`=?")
		if err == nil {
			break
		}
		if strings.Index(err.Error(), "Error 1146:") > -1 {
			p.createMacdTable()
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
		stmt.QueryRow(datas[i].Time, symbol, typ).Scan(&has)
		if !has {
			unsave = append(unsave, datas[i])
		}
	}

	return unsave, nil
}

func (p *Mysql) saveMacds(symbol string, typ int, datas []Tdata) error {
	table := macdTable
	var stmt *sql.Stmt
	var err error
	for i := 0; i < 2; i++ {
		stmt, err = p.db.Prepare("INSERT INTO `" + table + "`(`time`,`code`,`typ`,`emas`,`emal`,`dea`) values(?,?,?,?,?,?)")
		if err == nil {
			break
		}
		if strings.Index(err.Error(), "Error 1146:") > -1 {
			p.createMacdTable()
			continue
		}
	}
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i, c := 0, len(datas); i < c; i++ {
		data := &datas[i]
		_, e := stmt.Exec(data.Time, symbol, typ, data.Emas, data.Emal, data.DEA)
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
		glog.Warningln("insert macd error", err)
	}
	return err
}

func nextTime(t time.Time, typ int) time.Time {
	if typ < LDay {
		t = t.AddDate(0, 1, 0)
	} else {
		t = t.AddDate(1, 0, 0)
	}
	return t.Truncate(24 * time.Hour)
}

func (p *Mysql) SaveMacds(symbol string, typ int, datas []Tdata) error {
	table := macdTable
	var err error
	count := len(datas)
	if count < 1 {
		return nil
	}

	unsave := []Tdata{}
	now := time.Now().Truncate(24 * time.Hour)

	t := datas[0].Time
	p.db.QueryRow("SELECT `time` FROM `"+table+"` WHERE `code`=? AND `typ`=? ORDER BY `time` DESC LIMIT 0,1", symbol, typ).Scan(&t)
	t = nextTime(t, typ)

	for ; count > 0 && !datas[count-1].Time.Before(now); count-- {
	}

	for i := 0; i < count; i++ {
		if datas[i].Time.Before(t) {
			continue
		}
		unsave = append(unsave, datas[i])
		t = nextTime(t, typ)
	}
	saved := 0
	for i := 0; i < 10 && len(unsave) > 0; i++ {
		err = p.saveMacds(symbol, typ, unsave)
		sum := len(unsave)
		unsave, err = p.checkMacds(symbol, typ, unsave)
		saved = sum - len(unsave)
		if saved > 0 {
			i--
		}
	}

	return err
}

func (p *Mysql) GetStartTime(symbol string, typ int) time.Time {
	table := macdTable
	t := Market_begin_day
	for i := 4; i > -1; i-- {
		err := p.db.QueryRow("SELECT `time` FROM `"+table+"` WHERE `code`=? AND `typ`=? ORDER BY `time` DESC LIMIT ?,1", symbol, typ, i).Scan(&t)
		if err == nil {
			break
		}
		if err != sql.ErrNoRows {
			glog.Warningln(err)
			break
		}
	}
	return t
}
