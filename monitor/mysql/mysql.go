package mysql

import (
	"database/sql"
	"log"
	"opms/monitor/utils"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
)

func GenerateMySQLStats(wg *sync.WaitGroup, mysql *xorm.Engine, db_id int, host string, port int, alias string, is_alert int) {
	//添加异常处理
	defer func() {
		if err := recover(); err != nil{
		   // 出现异常，继续
		   log.Printf("Error: %v", err)
		   (*wg).Done()
		}
	}()
	

	//连接字符串
	dsn, err := GetDsn(mysql, db_id, 2)
	if err != nil {
		utils.LogDebugf("GetDsn failed: %s", err.Error())
	}
	//建立连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		utils.LogDebugf("Open Connection failed: %s", err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		utils.LogDebugf("DB Ping %s failed: %s", alias, err.Error())
		MoveToHistory(mysql, "pms_asset_status", "asset_id", db_id)
		MoveToHistory(mysql, "pms_mysql_status", "db_id", db_id)

		sql := `insert into pms_asset_status(asset_id, asset_type, host, port, alias, connect, created) 
				values(?,?,?,?,?,?,?)`
		_, err = mysql.Exec(sql, db_id, 2, host, port, alias, -1, time.Now().Unix())
		if err != nil {
			log.Printf("%s: %s", sql, err.Error())
		}

		sql = `insert into pms_mysql_status(db_id, host, port, alias, connect, created) 
		values(?,?,?,?,?,?)`
		_, err = mysql.Exec(sql, db_id, host, port, alias, -1, time.Now().Unix())
		if err != nil {
			log.Printf("%s: %s", sql, err.Error())
		}

		if is_alert == 1 {
			AlertConnect(mysql, db_id)
		}
	} else {
		log.Println("ping succeeded")

		//get sqlserver basic infomation
		GatherBasicInfo(db, mysql, db_id, host, port, alias)
		GatherDbStatus(mysql, db_id)
		//GetGlobalStatus(db)

		GatherMetricValue(db, mysql, db_id, host, port, alias)

		if is_alert == 1 {
			AlertBasicInfo(mysql, db_id)
		}
	}

	(*wg).Done()

}

func GatherBasicInfo(db *sql.DB, mysql *xorm.Engine, db_id int, host string, port int, alias string) error {

	connect := 1
	role := GetRole(db)
	version := GetVersion(db)
	//log.Printf("version: %s", version)

	golbalstatus, err := GetGlobalStatus(db)
	if err != nil {
		log.Printf("GetGlobalStatus failed : %s", err.Error())
		return err
	}

	uptime := golbalstatus["Uptime"]
	threads_connected := golbalstatus["Threads_connected"]
	threads_running := golbalstatus["Threads_running"]
	questions := golbalstatus["Questions"]
	com_commit := golbalstatus["Com_commit"]
	com_rollback := golbalstatus["Com_rollback"]
	innodb_log := golbalstatus["innodb_os_log_written"]

	open_tables := golbalstatus["Open_tables"]
	key_blocks_used, err := strconv.Atoi(golbalstatus["Key_blocks_used"])
	key_blocks_unused, err := strconv.Atoi(golbalstatus["Key_blocks_unused"])
	key_blocks_not_flushed := golbalstatus["Key_blocks_not_flushed"]

	key_reads, err := strconv.Atoi(golbalstatus["Key_reads"])
	key_read_requests, err := strconv.Atoi(golbalstatus["Key_read_requests"])
	key_writes, err := strconv.Atoi(golbalstatus["Key_writes"])
	key_write_requests, err := strconv.Atoi(golbalstatus["Key_write_requests"])

	golbalvariables, err := GetGlobalVariables(db)
	if err != nil {
		log.Printf("GetGlobalVariables failed : %s", err.Error())
		return err
	}
	max_connections := golbalvariables["max_connections"]
	max_connect_errors := golbalvariables["max_connect_errors"]
	open_files_limit := golbalvariables["open_files_limit"]
	table_open_cache := golbalvariables["table_open_cache"]
	key_buffer_size := golbalvariables["key_buffer_size"]
	sort_buffer_size := golbalvariables["sort_buffer_size"]
	join_buffer_size := golbalvariables["join_buffer_size"]

	threads_waits := GetProcessWaits(db)

	key_blocks_used_rate := 0
	if (key_blocks_used + key_blocks_unused) != 0 {
		key_blocks_used_rate = key_blocks_used * 100 / (key_blocks_used + key_blocks_unused)
	}

	key_buffer_read_rate := 0
	if key_read_requests != 0 {
		key_buffer_read_rate = key_reads * 100 / key_read_requests
	}

	key_buffer_write_rate := 0
	if key_write_requests != 0 {
		key_buffer_write_rate = key_writes * 100 / key_write_requests
	}

	// storage result
	session := mysql.NewSession()
	defer session.Close()
	// add Begin() before any action
	err = session.Begin()
	//storage stats into pms_asset_status
	//move old data to history table
	MoveToHistory(mysql, "pms_asset_status", "asset_id", db_id)

	sql := `insert into pms_asset_status(asset_id, asset_type, host, port, alias, role, version, connect, sessions, created) 
						values(?,?,?,?,?,?,?,?,?,?)`
	_, err = mysql.Exec(sql, db_id, 1, host, port, alias, role, version, connect, threads_connected, time.Now().Unix())
	if err != nil {
		log.Printf("%s: %s", sql, err.Error())
		err = session.Rollback()
		return err
	}

	//storage stats into pms_mysql_status
	MoveToHistory(mysql, "pms_mysql_status", "db_id", db_id)

	sql = `insert into pms_mysql_status(db_id, host, port, alias, connect, role, uptime, version, questions, com_commit, com_rollback, innodb_log, max_connections, max_connect_errors, open_files_limit, table_open_cache, open_tables, threads_connected, threads_running, threads_waits, key_buffer_size,sort_buffer_size,join_buffer_size,key_blocks_used,key_blocks_unused,key_blocks_not_flushed, key_blocks_used_rate, key_buffer_read_rate, key_buffer_write_rate, created) 
						values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err = mysql.Exec(sql, db_id, host, port, alias, connect, role, uptime, version, questions, com_commit, com_rollback, innodb_log, max_connections, max_connect_errors, open_files_limit, table_open_cache, open_tables, threads_connected, threads_running, threads_waits, key_buffer_size, sort_buffer_size, join_buffer_size, key_blocks_used, key_blocks_unused, key_blocks_not_flushed, key_blocks_used_rate, key_buffer_read_rate, key_buffer_write_rate, time.Now().Unix())
	if err != nil {
		log.Printf("%s: %s", sql, err.Error())
		err = session.Rollback()
		return err
	}

	// add Commit() after all actions
	err = session.Commit()
	return err
}

func GatherMetricValue(db *sql.DB, mysql *xorm.Engine, db_id int, host string, port int, alias string) {
	golbalstatus, err := GetGlobalStatus(db)
	if err != nil {
		log.Printf("GetGlobalStatus failed : %s", err.Error())
		return
	}

	session_total := golbalstatus["Threads_connected"]
	session_actives := golbalstatus["Threads_running"]
	innodb_buffer_reads := golbalstatus["innodb_buffer_pool_reads"]
	innodb_buffer_read_reqs := golbalstatus["innodb_buffer_pool_read_requests"]

	timestamp := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
	StorageMetricData(mysql, db_id, "TotalSessions", timestamp, session_total, "GAUGE")
	StorageMetricData(mysql, db_id, "ActiveSessions", timestamp, session_actives, "GAUGE")

	//get QPS
	qps := GetQPS(mysql, db_id)
	StorageMetricData(mysql, db_id, "Queries Per Second", timestamp, strconv.Itoa(qps), "GAUGE")

	//get TPS
	tps := GetTPS(mysql, db_id)
	StorageMetricData(mysql, db_id, "Transactions Per Second", timestamp, strconv.Itoa(tps), "GAUGE")

	//get Buffer Cache Hit
	bchit := -1
	buffer_reads, err := strconv.Atoi(innodb_buffer_reads)
	if err != nil {
		log.Printf("GetGlobalStatus failed : %s", err.Error())
		buffer_reads = -1
	}
	buffer_reqs, err := strconv.Atoi(innodb_buffer_read_reqs)
	if err != nil {
		log.Printf("GetGlobalStatus failed : %s", err.Error())
		buffer_reqs = -1
	}
	if buffer_reads > 0 && buffer_reqs > 0 {
		bchit = (1 - buffer_reads/buffer_reqs) * 100
	}
	StorageMetricData(mysql, db_id, "Buffer Cache Hit", timestamp, strconv.Itoa(bchit), "GAUGE")

	//get Redo
	log_per_sec := GetLogPerSecond(mysql, db_id)
	StorageMetricData(mysql, db_id, "Log Per Second", timestamp, strconv.Itoa(log_per_sec), "GAUGE")
}

// Storage metric data
func StorageMetricData(mysql *xorm.Engine, db_id int, metric string, timestamp string, value string, counterType string) {

	sql := `insert into pms_metric_data(db_id, metric, timestamp, value, counterType) 
			values(?,?,?,?,?)`
	_, err := mysql.Exec(sql, db_id, metric, timestamp, value, counterType)
	if err != nil {
		log.Printf("StorageMetricData -- %s: %s", sql, err.Error())
	}
}

func MoveToHistory(mysql *xorm.Engine, table_name string, key_name string, key_value int) {
	sql := `insert into ` + table_name + `_his select * from ` + table_name + ` where ` + key_name + ` = ?`
	_, err := mysql.Exec(sql, key_value)
	if err != nil {
		log.Printf("%s: %s", sql, err.Error())
	}

	sql = `delete from ` + table_name + ` where ` + key_name + ` = ?`
	_, err = mysql.Exec(sql, key_value)
	if err != nil {
		log.Printf("%s: %s", sql, err.Error())
	}
}

func GetDsn(db *xorm.Engine, db_id int, asset_type int) (string, error) {
	var dsn string
	var sql string
	if asset_type == 1 {
		sql = `select concat("oracle://",username,":",password ,"@" , host , ":" , port , "/" , instance_name , "?sysdba=1") as dsn 
				from pms_asset_config where id = ? and asset_type = ?`
	} else if asset_type == 2 {
		sql = `select concat(username,":",password,"@tcp(",host,":",port,")/",db_name,"?charset=utf8") from pms_asset_config where id = ? and asset_type = ?`
	} else if asset_type == 3 {
		sql = `select concat("server=",host,"\\",instance_name,";port",port,";database=",case db_name when "" then "master" end,";user id=",username,";password=",password,";encrypt=disable") from pms_asset_config where id = ? and asset_type = ?`
	} else {
		sql = `select "" from pms_asset_config where id = ? and asset_type = ?`
	}

	_, err := db.SQL(sql, db_id, asset_type).Get(&dsn)
	if err != nil {
		log.Fatal(err)
	}

	return dsn, err
}
