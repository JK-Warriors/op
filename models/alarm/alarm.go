package alarm

import (
	//"fmt"
	"opms/models"
	//"opms/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Alert struct {
    Id 				int				`orm:"pk;column(Id);"`
    Asset_Id 		int		    	`orm:"column(asset_id);"`
    Name 			string		    `orm:"column(name);"`
    Severity 		string		    `orm:"column(severity);"`
    Templateid 		int		    	`orm:"column(templateid);"`
    Subject 		string		    `orm:"column(subject);"`
    Message 		string		    `orm:"column(message);"`
    Status 			int				`orm:"column(status);"`
    Send_Mail 					int				`orm:"column(send_mail);"`
    Send_Mail_List 				string			`orm:"column(send_mail_list);"`
    Send_Mail_Status 			int				`orm:"column(send_mail_status);"`
    Send_Mail_Retries 			int				`orm:"column(send_mail_retries);"`
    Send_Mail_Error 			string			`orm:"column(send_mail_error);"`
    Send_WeChat 				int				`orm:"column(send_wechat);"`
    Send_WeChat_Status 			int				`orm:"column(send_wechat_status);"`
    Send_WeChat_Retries 		int				`orm:"column(send_wechat_retries);"`
    Send_WeChat_Error 			string			`orm:"column(send_wechat_error);"`
    Send_SMS 					int				`orm:"column(send_sms);"`
    Send_SMS_List 				string			`orm:"column(send_sms_list);"`
    Send_SMS_Status 			int				`orm:"column(send_sms_status);"`
    Send_SMS_Retries 			int				`orm:"column(send_sms_retries);"`
    Send_SMS_Error 				string			`orm:"column(send_sms_error);"`
    Created 					int64			`orm:"column(created);"`
}




func (this *Alert) TableName() string {
	return models.TableName("alerts")
}
func init() {
	orm.RegisterModel(new(Alert))
}


//获取告警列表
func ListAlerts(condArr map[string]string, page int, offset int) (num int64, err error, alerts []Alert) {
	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable(models.TableName("alerts"))
	cond := orm.NewCondition()

	if condArr["search_name"] != "" {
		cond = cond.And("name__icontains", condArr["search_name"])
	}

	qs = qs.SetCond(cond)
	if page < 1 {
		page = 1
	}
	if offset < 1 {
		offset, _ = beego.AppConfig.Int("pageoffset")
	}
	start := (page - 1) * offset

	qs = qs.OrderBy("-id")
	nums, errs := qs.Limit(offset, start).All(&alerts)
	return nums, errs, alerts
}


//统计数量
func CountAlerts(condArr map[string]string) int64 {
	o := orm.NewOrm()
	qs := o.QueryTable(models.TableName("alerts"))
	cond := orm.NewCondition()

	if condArr["search_name"] != "" {
		cond = cond.And("name__icontains", condArr["search_name"])
	}
	
	num, _ := qs.SetCond(cond).Count()
	return num
}

//获取近7天所有告警
func ListAllAlerts() (num int64, err error, alerts []Alert) {
	o := orm.NewOrm()
	o.Using("default")

	sql := `select * from pms_alerts where status = 1 and created > UNIX_TIMESTAMP() - 3600*24*7`
	nums, errs := o.Raw(sql).QueryRows(&alerts)

	return nums, errs, alerts
}

type AlertGroup struct {
    Asset_Id 		int		    	`orm:"column(asset_id);"`
    Count 			int				`orm:"column(alertcount);"`
    Rate 			int				`orm:"column(rate);"`
}
func ListAlertGroup() (num int64, err error, alertgroup []AlertGroup) {
	o := orm.NewOrm()
	o.Using("default")
	var li_count int
	sql := `select count(1) from pms_alerts where status = 1 and created > UNIX_TIMESTAMP() - 3600*24*7`
	errs := o.Raw(sql).QueryRow(&li_count)

	sql = `select asset_id, count(1) alertcount, floor(count(1)*100/?) as rate
			from pms_alerts 
			where status = 1 
			and created > UNIX_TIMESTAMP() - 3600*24*7 
			group by asset_id 
			order by 2 desc 
			limit 5`
	nums, errs := o.Raw(sql,li_count).QueryRows(&alertgroup)

	return nums, errs, alertgroup
}
