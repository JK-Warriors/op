package common

import (
	"log"
	"time"
	"strings"

	"github.com/xormplus/xorm"
)


type Dr struct {
	Id          int    `xorm:"int 'id'"`
	Bs_Name     string `xorm:"varchar(200) 'bs_name'"`
	Asset_Type   	int    `xorm:"int 'asset_type'"`
	Db_Id_P     int    `xorm:"int 'db_id_p'"`
	Host_P      string `xorm:"varchar(20) 'host_p'"`
	Port_P      int    `xorm:"int 'port_p'"`
	Alias_P     string `xorm:"varchar(200) 'alias_p'"`
	Inst_Name_P string `xorm:"varchar(50) 'inst_name_p'	"`
	Db_Id_S     int    `xorm:"int 'db_id_s'"`
	Host_S      string `xorm:"varchar(20) 'host_s'"`
	Port_S      int    `xorm:"int 'port_s'"`
	Alias_S     string `xorm:"varchar(200) 'alias_s'"`
	Inst_Name_S string `xorm:"varchar(50) 'inst_name_s'"`
	Db_Name   	string `xorm:"varchar(50) 'db_name'"`
	Is_Shift    int    `xorm:"int 'is_shift'"`
}


type Trigger struct {
	Id       			int `xorm:"int 'id' "`
	Asset_Id   			int `xorm:"int 'asset_id' "`
	Asset_Type   		int `xorm:"int 'asset_type' "`
	Name   				string `xorm:"varchar(200) 'name' "`
	TemplateId   		int 	`xorm:"int 'templateid' "`
	Trigger_Type   		string `xorm:"varchar(200) 'trigger_type' "`
	Severity   			string `xorm:"varchar(200) 'severity' "`
	Expression   		string `xorm:"varchar(1000) 'expression' "`
	Description   		string `xorm:"varchar(2000) 'description' "`
	Status   			int `xorm:"int 'status' "`
	Recovery_Mode   	int `xorm:"int 'recovery_mode' "`
	Recovery_Expression string `xorm:"varchar(1000) 'recovery_expression' "`
	Recovery_Description string `xorm:"varchar(2000) 'recovery_description' "`
	Created  			int64  `xorm:"bigint 'created' "`
}


func AddAlert(mysql *xorm.Engine, db_id int, item_name string, item_value string, tri Trigger){
	var count int
	sql := `select count(1) from pms_alerts where asset_id = ? and templateid = ? and created > UNIX_TIMESTAMP() - 3600`
	_, err := mysql.SQL(sql, db_id, tri.TemplateId).Get(&count)
	if err != nil {
		log.Printf("%s: %s", sql, err.Error())
	}

	send_mail, send_wechat, send_sms := 0, 0, 0
	_, err = mysql.SQL(`select alert_mail from pms_asset_config where id = ?`, db_id).Get(&send_mail)
	_, err = mysql.SQL(`select alert_wechat from pms_asset_config where id = ?`, db_id).Get(&send_wechat)
	_, err = mysql.SQL(`select alert_sms from pms_asset_config where id = ?`, db_id).Get(&send_sms)

	var send_mail_list, send_sms_list string
	_, err = mysql.SQL(`select value from pms_global_options where id = 'send_mail_to_list'`).Get(&send_mail_list)
	_, err = mysql.SQL(`select value from pms_global_options where id = 'send_sms_to_list'`).Get(&send_sms_list)


	if count == 0 {
		description := strings.Replace(tri.Description, "{ItemName}", item_name, -1)
		description = strings.Replace(description, "{ItemValue}", item_value, -1)

		sql := `insert into pms_alerts(asset_id, name, severity, templateid, subject, message, send_mail, send_mail_list, send_wechat, send_sms, send_sms_list, created)
				values(?,?,?,?,?,?,?,?,?,?,?,?)`
	
		_, err := mysql.Exec(sql, db_id, tri.Name, tri.Severity, tri.TemplateId, tri.Name, description, send_mail, send_mail_list, send_wechat, send_sms, send_sms_list, time.Now().Unix())
		if err != nil {
			log.Printf("%s: %s", sql, err.Error())
		}
	}

}

func AddRecoveryAlert(mysql *xorm.Engine, db_id int, item_name string, item_value string, tri Trigger){
	var created int64
	var count int
	sql := `select created from pms_alerts where asset_id = ? and templateid = ? order by id desc limit 1`
	_, err := mysql.SQL(sql, db_id, tri.TemplateId).Get(&created)
	if err != nil {
		log.Printf("%s: %s", sql, err.Error())
	}

	if created > 0{
		sql := `select count(1) from pms_alerts where asset_id = ? and templateid = ? and severity = 'Info' and created > ?`
		_, err := mysql.SQL(sql, db_id, tri.TemplateId, created).Get(&count)
		if err != nil {
			log.Printf("%s: %s", sql, err.Error())
		}

		if count == 0 {
			// no recovery
			description := strings.Replace(tri.Recovery_Description, "{ItemName}", item_name, -1)
			description = strings.Replace(description, "{ItemValue}", item_value, -1)
	
			sql := `insert into pms_alerts(asset_id, name, severity, templateid, message, created)
					values(?,?,?,?,?,?)`
		
			_, err := mysql.Exec(sql, db_id, tri.Name, "Info", tri.TemplateId, description, time.Now().Unix())
			if err != nil {
				log.Printf("%s: %s", sql, err.Error())
			}

		}

	}
}


func GetTriggers(mysql *xorm.Engine, db_id int, trigger_type string) (triconf []Trigger) {
	
	sql := `select id, asset_id, asset_type, name, templateid,trigger_type,severity,expression,description,status,recovery_mode,recovery_expression,recovery_description, created 
			from pms_triggers where asset_id = ? and trigger_type = ?`

	err := mysql.SQL(sql, db_id, trigger_type).Find(&triconf)
	if err != nil {
		log.Printf("%s: %s", sql, err.Error())
	}

	return triconf
}