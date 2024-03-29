package routers

import (
	"opms/controllers/alarm"
	"opms/controllers/cfg_screen"
	"opms/controllers/cfg_trigger"
	"opms/controllers/config_global"
	"opms/controllers/dbconfig"
	"opms/controllers/cfg_os"
	"opms/controllers/demo"
	"opms/controllers/dr_business"
	"opms/controllers/dr_config"
	"opms/controllers/dr_oper"
	"opms/controllers/logs"
	"opms/controllers/messages"
	"opms/controllers/mssql"
	"opms/controllers/mysql"
	"opms/controllers/oracle"
	"opms/controllers/os"
	"opms/controllers/roles"
	"opms/controllers/screen"
	"opms/controllers/users"
	"opms/controllers/index"
	hc "opms/controllers/healthcheck"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &index.MainController{})

	//用户
	beego.Router("/login", &users.LoginUserController{})
	beego.Router("/logout", &users.LogoutUserController{})
	beego.Router("/system/user/manage", &users.ManageUserController{})
	beego.Router("/system/user/ajax/status", &users.AjaxStatusUserController{})
	beego.Router("/system/user/add", &users.AddUserController{})
	beego.Router("/system/user/edit/:id", &users.EditUserController{})
	beego.Router("/system/user/ajax/delete", &users.AjaxDeleteUserController{})
	beego.Router("/system/user/avatar", &users.AvatarUserController{})
	beego.Router("/system/user/ajax/search", &users.AjaxSearchUserController{}) //搜索用户名匹配
	beego.Router("/system/user/show/:id", &users.ShowUserController{})
	beego.Router("/system/user/profile", &users.EditUserProfileController{})
	beego.Router("/system/user/password", &users.EditUserPasswordController{})
	beego.Router("/system/user/ajax/reset_passwd", &users.AjaxResetPasswordController{})

	beego.Router("/system/my/manage", &users.ShowUserController{})

	//告警
	beego.Router("/alarm/alarm/manage", &alarm.ManageAlarmController{})
	beego.Router("/alarm/history/list", &alarm.AlarmHistoryListController{})

	//消息
	beego.Router("/system/message/manage", &messages.ManageMessageController{})
	beego.Router("/system/message/ajax/delete", &messages.AjaxDeleteMessageController{})
	beego.Router("/system/message/ajax/status", &messages.AjaxStatusMessageController{})

	//角色
	beego.Router("/system/role/manage", &roles.ManageRoleController{})
	beego.Router("/system/role/ajax/delete", &roles.AjaxDeleteRoleController{})
	beego.Router("/system/role/add", &roles.FormRoleController{})
	beego.Router("/system/role/edit/:id", &roles.FormRoleController{})
	//角色成员
	beego.Router("/system/role/user/:id", &roles.ManageRoleUserController{})
	beego.Router("/system/role/user/add/:id", &roles.FormRoleUserController{})
	beego.Router("/system/role/user/ajax/delete", &roles.AjaxDeleteRoleUserController{})

	//角色权限
	beego.Router("/system/role/permission/:id", &roles.ManageRolePermissionController{})
	beego.Router("/system/role/permission/ajax/delete", &roles.AjaxDeleteRolePermissionController{})

	//权限
	beego.Router("/system/permission/manage", &roles.ManagePermissionController{})
	beego.Router("/system/permission/ajax/delete", &roles.AjaxDeletePermissionController{})
	beego.Router("/system/permission/add", &roles.FormPermissionController{})
	beego.Router("/system/permission/edit/:id", &roles.FormPermissionController{})

	//日志
	beego.Router("/system/log/manage", &logs.ManageLogController{})
	beego.Router("/system/log/ajax/delete", &logs.AjaxDeleteLogController{})

	//业务系统配置
	beego.Router("/config/dr_business/manage", &dr_business.ManageBusinessController{})
	beego.Router("/config/dr_business/add", &dr_business.AddBusinessController{})
	beego.Router("/config/dr_business/edit", &dr_business.EditBusinessController{})
	beego.Router("/config/dr_business/ajax/delete", &dr_business.AjaxDeleteBusinessController{})

	//数据库配置
	beego.Router("/config/db/manage", &dbconfig.ManageDBConfigController{})
	beego.Router("/config/db/add", &dbconfig.AddDBConfigController{})
	beego.Router("/config/db/edit/:id", &dbconfig.EditDBConfigController{})
	beego.Router("/config/db/ajax/status", &dbconfig.AjaxStatusDBConfigController{})
	beego.Router("/config/db/ajax/delete", &dbconfig.AjaxDeleteDBConfigController{})
	beego.Router("/config/db/ajax/connect", &dbconfig.AjaxConnectDBConfigController{})

	//操作系统配置
	beego.Router("/config/os/manage", &cfg_os.ManageOSConfigController{})
	beego.Router("/config/os/add", &cfg_os.AddOSConfigController{})
	beego.Router("/config/os/edit/:id", &cfg_os.EditOSConfigController{})
	beego.Router("/config/os/ajax/status", &cfg_os.AjaxStatusOSConfigController{})
	beego.Router("/config/os/ajax/delete", &cfg_os.AjaxDeleteOSConfigController{})
	beego.Router("/config/os/ajax/connect", &cfg_os.AjaxConnectOSConfigController{})

	//容灾配置
	beego.Router("/config/dr_config/manage", &dr_config.ManageDrController{})
	beego.Router("/config/dr_config/add", &dr_config.AddDrController{})
	beego.Router("/config/dr_config/edit/:id", &dr_config.EditDrController{})
	beego.Router("/config/dr_config/ajax/status", &dr_config.AjaxStatusDrConfigController{})
	beego.Router("/config/dr_config/ajax/delete", &dr_config.AjaxDeleteDrConfigController{})

	//全局配置
	beego.Router("/config/config_global/manage", &config_global.ManageGlobalController{})
	beego.Router("/config/config_global/ajax/save", &config_global.AjaxSaveGlobalController{})
	beego.Router("/config/cfg_trigger/manage", &cfg_trigger.ManageTriggerController{})
	beego.Router("/config/cfg_trigger/edit/:id", &cfg_trigger.EditTriggerController{})

	beego.Router("/config/screen/manage", &cfg_screen.ManageScreenController{})
	beego.Router("/config/screen/ajax/save", &cfg_screen.AjaxSaveScreenController{})

	//资产状态
	// beego.Router("/asset/status/manage", &asset.ManageAssetController{})
	beego.Router("/oracle/status/manage", &oracle.ManageOracleController{})
	beego.Router("/oracle/tbs/manage", &oracle.ManageOracleTbsController{})
	beego.Router("/oracle/asm/manage", &oracle.ManageOracleAsmController{})

	beego.Router("/mysql/status/manage", &mysql.ManageMysqlController{})
	beego.Router("/mysql/resource/manage", &mysql.ResourceMysqlController{})
	beego.Router("/mysql/key/manage", &mysql.KeyMysqlController{})

	beego.Router("/mssql/status/manage", &mssql.ManageMssqlController{})

	beego.Router("/os/status/manage", &os.ManageOSController{})
	beego.Router("/os/disk/manage", &os.ManageOSDiskController{})
	beego.Router("/os/io/manage", &os.ManageOSDiskIOController{})

	//操作
	beego.Router("/operation/dr_manage/list", &dr_oper.ListDrController{})
	beego.Router("/operation/dr_manage/detail/:id", &dr_oper.DetailDrController{})
	
	beego.Router("/operation/screen/view/:id", &dr_oper.ScreenDrViewController{})
	beego.Router("/operation/dr_switch/switchover", &dr_oper.AjaxDrSwitchoverController{})
	beego.Router("/operation/dr_switch/failover", &dr_oper.AjaxDrFailoverController{})
	beego.Router("/operation/dr_switch/process", &dr_oper.AjaxDrProcessController{})
	beego.Router("/operation/dr_active/startread", &dr_oper.AjaxDrStartReadController{})
	beego.Router("/operation/dr_active/stopread", &dr_oper.AjaxDrStopReadController{})
	beego.Router("/operation/dr_sync/startsync", &dr_oper.AjaxDrStartSyncController{})
	beego.Router("/operation/dr_sync/stopsync", &dr_oper.AjaxDrStopSyncController{})
	beego.Router("/operation/dr_snapshot/manage", &dr_oper.ManageDrSnapshotController{})
	beego.Router("/operation/dr_snapshot/detail/:id", &dr_oper.DetailSnapshotController{})
	beego.Router("/operation/dr_snapshot/startsnapshot", &dr_oper.AjaxDrStartSnapshotController{})
	beego.Router("/operation/dr_snapshot/stopsnapshot", &dr_oper.AjaxDrStopSnapshotController{})
	beego.Router("/operation/dr_recover/manage", &dr_oper.ManageDrRecoverController{})
	beego.Router("/operation/dr_recover/detail/:id", &dr_oper.DetailRecoverController{})
	beego.Router("/operation/dr_recover/oper/:id", &dr_oper.OperDrRecoverController{})
	beego.Router("/operation/dr_recover/flashback", &dr_oper.AjaxDrFlashbackController{})
	beego.Router("/operation/dr_recover/recover", &dr_oper.AjaxDrRecoverController{})

	//巡检
	beego.Router("/healthcheck/task/manage", &hc.ManageTaskController{})
	beego.Router("/healthcheck/history/manage", &hc.ManageHistoryController{})
	beego.Router("/healthcheck/config/manage", &hc.ManageConfigController{})
	beego.Router("/healthcheck/config/edit", &hc.EditConfigController{})
	beego.Router("/healthcheck/template/manage", &hc.ManageTemplateController{})


	//大屏
	beego.Router("/screen/manage", &screen.ManageScreenController{})

	//UI demo
	beego.Router("/demo/index", &demo.DemoController{})
	beego.Router("/demo/form", &demo.FormController{})
	beego.Router("/demo/base", &demo.BaseController{})
	beego.Router("/demo/dashboard", &demo.DashboardController{})
	beego.Router("/demo/dgscreen", &demo.DgscreenController{})
}
