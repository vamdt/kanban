package mysql

import (
	"database/sql"
	"flag"
	"strings"

	. "../../base"
	"../../store"
	_ "github.com/go-sql-driver/mysql"

	"github.com/golang/glog"
)

const (
	categoryTable = "category"
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
		if n > num {
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

func (p *Mysql) LoadTDatas(table string) (res []Tdata, err error) {
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

func (p *Mysql) SaveTDatas(table string, datas []Tdata) error {
	p.createDayTdataTable(table)
	var stmt *sql.Stmt
	var err error
	for i := 0; i < 2; i++ {
		stmt, err = p.db.Prepare("INSERT INTO `" + table + "`(`time`,`open`,`high`,`low`,`close`,`volume`) values(?,?,?,?,?,?)")
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
			err = e
		}
	}

	if err != nil {
		glog.Warningln("insert tdata error", err)
	}
	return err
}

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

func (p *Mysql) SaveTicks(table string, ticks []Tick) error {
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
			err = e
		}
	}
	if err != nil {
		glog.Warningf("insert tick error %v", err)
	}
	return err
}

func (p *Mysql) createCategorieTable() {
	table := categoryTable
	sql := "CREATE TABLE IF NOT EXISTS `" + table + "` (" +
		"`id` INT(11) NOT NULL AUTO_INCREMENT," +
		"`pid` INT(11) NOT NULL DEFAULT 0," +
		"`factor` INT(11) NOT NULL DEFAULT 0," +
		"`leaf` TINYINT(1) NOT NULL DEFAULT 0," +
		"`deleted` TINYINT(1) NOT NULL DEFAULT 0," +
		"`name` VARCHAR(128) NOT NULL DEFAULT ''," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY (`pid`, `name`)" +
		")"
	_, err := p.db.Exec(sql)
	if err != nil {
		glog.Warningln("create table err", table, err)
	}
}

func (p *Mysql) LoadCategories() (res []CategoryItemInfo, err error) {
	table := categoryTable
	cols := "`id`,`name`,`pid`,`leaf`,`factor`"
	var rows *sql.Rows
	rows, err = p.db.Query("SELECT " + cols + " FROM `" + table + "`")
	if err != nil {
		glog.Warningln(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		d := CategoryItemInfo{}
		if err = rows.Scan(&d.Id, &d.Name, &d.Pid, &d.Leaf, &d.Factor); err != nil {
			glog.Warningln(err)
			continue
		}
		res = append(res, d)
	}
	if err = rows.Err(); err != nil {
		glog.Warningln(err)
	}
	return
}

func (p *Mysql) getMaxConnections() int {
	num := 0
	p.db.QueryRow("SELECT @@max_connections").Scan(&num)
	return num
}

func (p *Mysql) GetOrInsertCategoryItem(info *CategoryItemInfo) (id int, err error) {
	table := categoryTable
	for i := 0; i < 3; i++ {
		err = p.db.QueryRow("SELECT `id` FROM `"+table+"` WHERE pid=? AND name=?",
			info.Pid, info.Name).Scan(&id)

		if err == nil {
			break
		}

		if strings.Index(err.Error(), "Error 1146:") > -1 {
			p.createCategorieTable()
			continue
		}
		if err == sql.ErrNoRows {
			_, err = p.db.Exec("INSERT INTO `"+table+"`(`pid`,`name`,`leaf`) values(?,?,?)",
				info.Pid, info.Name, info.Leaf)
			continue
		}
		glog.Warningln(err)
	}
	return
}

func (p *Mysql) SaveCategoryItemWithPid(c CategoryItem, pid int) (err error) {
	id := 0
	info := CategoryItemInfo{Name: c.Name, Pid: pid, Leaf: false}
	id, err = p.GetOrInsertCategoryItem(&info)
	if err != nil {
		glog.Warningln(err)
		return
	}

	for _, info := range c.Info {
		info.Pid = id
		info.Leaf = true
		p.GetOrInsertCategoryItem(&info)
	}

	p.SaveCategories(c.Sub, id)
	return
}

func (p *Mysql) SaveCategories(c Category, pid int) (err error) {
	if c == nil {
		return
	}
	for _, cate := range c {
		err = p.SaveCategoryItemWithPid(cate, pid)
	}
	return
}

func (p *Mysql) SaveCategoryItemInfoFactor(datas []CategoryItemInfo) {
	table := categoryTable
	stmt, err := p.db.Prepare("UPDATE `" + table + "` SET `factor`=? WHERE `id`=?")
	if err != nil {
		glog.Warningln(err)
		return
	}
	defer stmt.Close()
	for i, c := 0, len(datas); i < c; i++ {
		if datas[i].Factor < 1 {
			continue
		}
		_, e := stmt.Exec(datas[i].Factor, datas[i].Id)
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
		}
	}
}

func (p *Mysql) changeStar(pid int, symbol string, star bool) {
	table := categoryTable
	info := CategoryItemInfo{Name: "star", Pid: pid, Leaf: false}
	starId, _ := p.GetOrInsertCategoryItem(&info)
	if starId < 1 {
		return
	}
	info.Name = "unstar"
	unstarId, _ := p.GetOrInsertCategoryItem(&info)
	if unstarId < 1 {
		return
	}

	id := 0
	err := p.db.QueryRow("SELECT `id`,`pid` FROM `"+table+"` WHERE `pid` in (?,?) AND name=?",
		starId, unstarId, symbol).Scan(&id, &pid)
	if err == sql.ErrNoRows {
		info.Name = symbol
		info.Pid = unstarId
		if star {
			info.Pid = starId
		}
		info.Leaf = true
		p.GetOrInsertCategoryItem(&info)
		return
	}
	expPid := unstarId
	if star {
		expPid = starId
	}
	if pid != expPid {
		p.db.Exec("UPDATE `"+table+"` SET `pid`=? WHERE `id`=?", expPid, id)
	}
}

func (p *Mysql) Star(pid int, symbol string) {
	p.changeStar(pid, symbol, true)
}

func (p *Mysql) UnStar(pid int, symbol string) {
	p.changeStar(pid, symbol, false)
}

func (p *Mysql) IsStar(pid int, symbol string) bool {
	table := categoryTable
	count := 0
	starSql := "SELECT `id` FROM `" + table + "` WHERE pid=? AND name='star'"
	err := p.db.QueryRow("SELECT COUNT(1) FROM `"+table+"` WHERE `pid`=("+starSql+") AND name=?",
		pid, symbol).Scan(&count)
	return err == nil && count > 0
}

func (p *Mysql) Lucky(pid int, symbol string) string {
	table := categoryTable
	for {
		err := p.db.QueryRow("SELECT `name` FROM `"+table+"` WHERE `leaf`=1 AND name!=? LIMIT 0,1",
			symbol).Scan(&symbol)
		if err != nil {
			glog.Warningln(err)
			break
		}
		if p.IsStar(pid, symbol) {
			continue
		}
		break
	}
	return symbol
}
