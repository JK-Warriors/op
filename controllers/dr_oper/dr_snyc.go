package dr_oper

import (
	"log"
	"opms/controllers"
	"opms/lib/exception"

	. "opms/models/dr_oper"
	. "opms/models/users"
	"opms/utils"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils/pagination"
	"github.com/godror/godror"
)

//业务系统管理
type ManageDrSyncController struct {
	controllers.BaseController
}

func (this *ManageDrSyncController) Get() {
	//权限检测
	if !strings.Contains(this.GetSession("userPermission").(string), "oper-sync-manage") {
		this.Abort("401")
	}

	page, err := this.GetInt("p")
	if err != nil {
		page = 1
	}

	offset, err1 := beego.AppConfig.Int("pageoffset")
	if err1 != nil {
		offset = 15
	}

	search_name := this.GetString("search_name")
	condArr := make(map[string]string)
	condArr["search_name"] = search_name

	countDr := CountOracleDrConfig(condArr)

	paginator := pagination.SetPaginator(this.Ctx, offset, countDr)
	_, _, dr := ListOracleDr(condArr, page, offset)

	this.Data["paginator"] = paginator
	this.Data["condArr"] = condArr
	this.Data["dr"] = dr
	this.Data["countDr"] = countDr

	userid, _ := this.GetSession("userId").(int64)
	user, _ := GetUser(userid)
	this.Data["user"] = user

	this.TplName = "dr_oper/sync-index.tpl"
}

type AjaxDrStartSyncController struct {
	controllers.BaseController
}

func (this *AjaxDrStartSyncController) Post() {
	//权限检测
	// if !strings.Contains(this.GetSession("userPermission").(string), "oper-switch-view") {
	// 	this.Abort("401")
	// }
	bs_id, _ := this.GetInt("bs_id")
	asset_type, _ := this.GetInt("asset_type")
	op_type := "STARTSYNC"

	//灾备配置检查
	cfg_count, err := CheckDrConfig(bs_id)
	if cfg_count == 0 {
		//没有配置容灾库
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "该系统没有配置容灾库"}
		this.ServeJSON()
		return
	}

	sta_id, err := GetStandbyDBId(bs_id)
	if err != nil {
		utils.LogDebug("GetStandbyDBID failed: " + err.Error())
		return
	}
	utils.LogDebug("GetStandbyDBID successfully.")

	dsn_s, err := GetDsn(sta_id, asset_type)
	if err != nil {
		utils.LogDebug("GetStandbyDsn failed: " + err.Error())
	}
	utils.LogDebug("GetStandbyDsn successfully.")

	on_process, err := GetOnProcess(bs_id)
	if on_process == 1 {
		utils.LogDebug("There is another opration on process.")

		this.Data["json"] = map[string]interface{}{"code": 0, "message": "有另外一个操作正在进行中"}
		this.ServeJSON()
		return
	} else {
		exception.Try(func() {

			utils.LogDebug("操作加锁")
			OperationLock(bs_id, op_type)

			utils.LogDebug("初始化切换实例")
			op_id := utils.SnowFlakeId()
			Init_OP_Instance(op_id, bs_id, asset_type, op_type)
			//asset_type = 5
			utils.LogDebug("正式开始启动同步任务")
			if asset_type == 1 { //oracle
				p_sta, err := godror.ParseConnString(dsn_s)
				if err != nil {
					utils.LogDebugf("%s: %w", dsn_s, err)
				}

				utils.LogDebug("开始启动同步...")
				s_result := OraStartSync(op_id, bs_id, p_sta)
				utils.LogDebug("启动同步结束")

				if s_result == 1 {
					utils.LogDebug("更新启动结果")
					Update_OP_Result(op_id, 1)
				} else {
					utils.LogDebug("备库启动同步失败，更新启动结果")
					Update_OP_Result(op_id, -1)
				}

				OperationUnlock(bs_id, op_type)
			} else if asset_type == 2 { //mysql
				//OraPrimaryToStandby
				//OraStandbyToPrimary

			} else if asset_type == 3 { //sqlserver
				//OraPrimaryToStandby
				//OraStandbyToPrimary

			}

		}).Catch(1, func(e exception.Exception) {
			log.Println(e.Id, e.Msg)
		}).Catch(2, func(e exception.Exception) {
			log.Println(e.Id, e.Msg)
		}).Catch(-1, func(e exception.Exception) {
			log.Println(e.Id, e.Msg)
		}).Finally(func() {
			OperationUnlock(bs_id, op_type)
		})

		this.Data["json"] = map[string]interface{}{"code": 1, "message": "操作完成"}
		this.ServeJSON()
	}
}

type AjaxDrStopSyncController struct {
	controllers.BaseController
}

func (this *AjaxDrStopSyncController) Post() {
	//权限检测
	// if !strings.Contains(this.GetSession("userPermission").(string), "oper-switch-view") {
	// 	this.Abort("401")
	// }
	bs_id, _ := this.GetInt("bs_id")
	asset_type, _ := this.GetInt("asset_type")
	op_type := "STOPSYNC"

	//灾备配置检查
	cfg_count, err := CheckDrConfig(bs_id)
	if cfg_count == 0 {
		//没有配置容灾库
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "该系统没有配置容灾库"}
		this.ServeJSON()
		return
	}

	sta_id, err := GetStandbyDBId(bs_id)
	if err != nil {
		utils.LogDebug("GetStandbyDBID failed: " + err.Error())
	}
	utils.LogDebug("GetStandbyDBID successfully.")

	dsn_s, err := GetDsn(sta_id, asset_type)
	if err != nil {
		utils.LogDebug("GetStandbyDsn failed: " + err.Error())
	}
	utils.LogDebug("GetStandbyDsn successfully.")

	on_process, err := GetOnProcess(bs_id)
	if on_process == 1 {
		utils.LogDebug("There is another opration on process.")

		this.Data["json"] = map[string]interface{}{"code": 0, "message": "有另外一个操作正在进行中"}
		this.ServeJSON()
		return
	} else {
		exception.Try(func() {

			utils.LogDebug("操作加锁")
			OperationLock(bs_id, op_type)

			utils.LogDebug("初始化切换实例")
			op_id := utils.SnowFlakeId()
			Init_OP_Instance(op_id, bs_id, asset_type, op_type)
			//asset_type = 5
			utils.LogDebug("正式开始停止同步任务")
			if asset_type == 1 { //oracle
				p_sta, err := godror.ParseConnString(dsn_s)
				if err != nil {
					utils.LogDebugf("%s: %w", dsn_s, err)
				}

				utils.LogDebug("开始停止同步...")
				s_result := OraStopSync(op_id, bs_id, p_sta)
				utils.LogDebug("停止同步结束")

				if s_result == 1 {
					utils.LogDebug("更新停止结果")
					Update_OP_Result(op_id, 1)
				} else {
					utils.LogDebug("备库停止同步失败，更新停止结果")
					Update_OP_Result(op_id, -1)
				}

				OperationUnlock(bs_id, op_type)
			} else if asset_type == 2 { //mysql
				//OraPrimaryToStandby
				//OraStandbyToPrimary

			} else if asset_type == 3 { //sqlserver
				//OraPrimaryToStandby
				//OraStandbyToPrimary

			}

		}).Catch(1, func(e exception.Exception) {
			log.Println(e.Id, e.Msg)
		}).Catch(2, func(e exception.Exception) {
			log.Println(e.Id, e.Msg)
		}).Catch(-1, func(e exception.Exception) {
			log.Println(e.Id, e.Msg)
		}).Finally(func() {
			OperationUnlock(bs_id, op_type)
		})

		this.Data["json"] = map[string]interface{}{"code": 1, "message": "操作完成"}
		this.ServeJSON()
	}

}
