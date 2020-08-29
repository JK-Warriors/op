<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>{{config "String" "globaltitle" ""}}</title>
{{template "inc/meta.tpl" .}}
</head><body class="sticky-header">
<section> {{template "inc/left.tpl" .}}
  <!-- main content start-->
  <div class="main-content" >
    <!-- header section start-->
    <div class="header-section">
      <!--toggle button start-->
      <a class="toggle-btn"><i class="fa fa-bars"></i></a> {{template "inc/user-info.tpl" .}} </div>
    <!-- header section end-->
    <!-- page heading start-->
    <div class="page-heading">
      <!-- <h3> 组织管理 {{template "users/nav.tpl" .}}</h3>-->
      <ul class="breadcrumb pull-left">
        <li> <a href="/config/db/manage">配置中心</a> </li>
        <li> <a href="/config/db/manage">资产配置</a> </li>
        <li class="active"> {{if gt .dbconf.Id 0}}编辑{{else}}新增{{end}}资产 </li>
      </ul>
      <div class="pull-right"><a href="/config/db/add" class="btn btn-success">+添加资产</a></div>
    </div>
    <!-- page heading end-->
    <!--body wrapper start-->
    <div class="wrapper">
      <div class="row">
        <div class="col-lg-12">
          <section class="panel">
            <header class="panel-heading"> {{.title}} </header>
            <div class="panel-body">
              <form class="form-horizontal adminex-form" id="dbconfig-form">
                <header><b> 基本信息 </b></header>
                <!--<div class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label"><span></span>业务系统名</label>
                  <div class="col-sm-10">
                    <select id="bs_id" name="bs_id" class="form-control">
                      <option value="">请选择系统</option>
                      {{range $k,$v := .bsconf}}
                        <option value="{{$v.Id}}" {{if eq $.dbconf.Bs_Id $v.Id}}selected{{end}}>{{$v.BsName}}</option>
                      {{end}}
                    </select>
                  </div>
                </div>-->
                <div class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label"><span>*</span>资产类型</label>
                  <div class="col-sm-10">
                    <select id="asset_type" name="asset_type" class="form-control">
                      <option value="">请选择类型</option>
                      <option value="1" {{if eq 1 .dbconf.Dbtype}}selected{{end}}>Oracle</option>
                      <option value="2" {{if eq 2 .dbconf.Dbtype}}selected{{end}}>MySQL</option>
                      <option value="3" {{if eq 3 .dbconf.Dbtype}}selected{{end}}>SQLServer</option>
                      <option value="99" {{if eq 99 .dbconf.Dbtype}}selected{{end}}>OS</option>
                    </select>
                  </div>
                </div>
                <div class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label"><span>*</span>主机IP</label>
                  <div class="col-sm-10">
                    <input type="text" id="host" name="host"  value="{{.dbconf.Host}}" class="form-control">
                  </div>
                </div>
                <div id="div_protocol" class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label"><span>*</span>协议</label>
                  <div class="col-sm-10">
                    <select id="protocol" name="protocol" class="form-control">
                      <option value="">请选择协议</option>
                      <option value="snmp" {{if eq "snmp" .dbconf.Protocol}}selected{{end}}>snmp</option>
                      <option value="winrm" {{if eq "winrm" .dbconf.Protocol}}selected{{end}}>winrm</option>
                    </select>
                  </div>
                </div>
                <div class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label"><span>*</span>端口</label>
                  <div class="col-sm-10">
                    <input type="text" id="port" name="port"  value="{{.dbconf.Port}}" class="form-control">
                  </div>
                </div>
                <div class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label"><span>*</span>别名</label>
                  <div class="col-sm-10">
                    <input type="text" id="alias" name="alias"  value="{{.dbconf.Alias}}" class="form-control">
                  </div>
                </div>
                <div id="div_inst_name" class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label"><span>*</span>实例名</label>
                  <div class="col-sm-10">
                    <input type="text" id="inst_name" name="instance_name"  value="{{.dbconf.InstanceName}}" class="form-control">
                  </div>
                </div>
                <div id="div_db_name" class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label"><span>*</span>数据库名</label>
                  <div class="col-sm-10">
                    <input type="text" id="db_name" name="db_name"  value="{{.dbconf.Dbname}}" class="form-control">
                  </div>
                </div>
                <div id="div_username" class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label"><span>*</span>用户名</label>
                  <div class="col-sm-10">
                    <input type="text" id="username" name="username"  value="{{.dbconf.Username}}" class="form-control" placeholder="请填写用户名">
                  </div>
                </div>
                <div id="div_password" class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label"><span>*</span>密码</label>
                  <div class="col-sm-10">
                    <input type="password" id="password" name="password"  value="{{.dbconf.Password}}" class="form-control" placeholder="请填写密码">
                  </div>
                </div>
                <div id="div_role" class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label"><span>*</span>角色</label>
                  <div class="col-sm-10">
                    <select name="role" class="form-control">
                      <option value="1" {{if eq 1 .dbconf.Role}}selected{{end}}>主</option>
                      <option value="2" {{if eq 2 .dbconf.Role}}selected{{end}}>备</option>
                    </select>
                  </div>
                </div>

                <div id="div_os_type" class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label"><span></span>主机类型</label>
                  <div class="col-sm-10">
                    <select id="os_type" name="os_type" class="form-control">
                      <option value="">请选择主机类型</option>
                      <option value="1" {{if eq 1 .dbconf.Ostype}}selected{{end}}>Linux</option>
                      <option value="2" {{if eq 2 .dbconf.Ostype}}selected{{end}}>Windows</option>
                      <option value="3" {{if eq 3 .dbconf.Ostype}}selected{{end}}>AIX</option>
                      <option value="4" {{if eq 4 .dbconf.Ostype}}selected{{end}}>HP-Unix</option>
                      <option value="5" {{if eq 5 .dbconf.Ostype}}selected{{end}}>Solaris</option>
                    </select>
                  </div>
                </div>
                <div id="div_os_protocol" class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label"><span></span>主机协议</label>
                  <div class="col-sm-10">
                    <select id="os_protocol" name="os_protocol" class="form-control">
                      <option value="">请选择主机协议</option>
                      <option value="ssh" {{if eq "ssh" .dbconf.OsProtocol}}selected{{end}}>ssh</option>
                      <option value="telnet" {{if eq "telnet" .dbconf.OsProtocol}}selected{{end}}>telnet</option>
                    </select>
                  </div>
                </div>
                <div id="div_os_username" class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label"><span></span>主机用户名</label>
                  <div class="col-sm-10">
                    <input type="text" id="os_username" name="os_username"  value="{{.dbconf.OsUsername}}" class="form-control" placeholder="请填写主机用户名">
                  </div>
                </div>
                <div id="div_os_password" class="form-group">
                  <label class="col-sm-2 col-sm-2 control-label"><span></span>主机密码</label>
                  <div class="col-sm-10">
                    <input type="password" id="os_password" name="os_password"  value="{{.dbconf.OsPassword}}" class="form-control" placeholder="请填写主机密码">
                  </div>
                </div>
                <div class="form-group">
                  <label class="col-lg-2 col-sm-2 control-label"></label>
                  <div class="col-lg-10">
                    <input type="hidden" id="id" name="id" value="{{.dbconf.Id}}">
                    <button type="button" onclick="checkConnect()" class="btn btn-primary">连接测试</button>
                    <button type="submit" class="btn btn-primary">提 交</button>
                  </div>
                </div>
              </form>
            </div>
          </section>
        </div>
      </div>
    </div>
    <!--body wrapper end-->
    <!--footer section start-->
    {{template "inc/foot-info.tpl" .}}
    <!--footer section end-->
  </div>
  <!-- main content end-->
</section>
{{template "inc/foot.tpl" .}}
<script src="/static/js/jquery-ui-1.10.3.min.js"></script>
<script>
    $(function() {// 初始化内容
        asset_type = {{.dbconf.Dbtype}};
        if(asset_type == "1"){      
            $("#div_inst_name").show();
            $("#div_username").show();
            $("#div_password").show();
            $("#div_protocol").hide();   
            $("#protocol").attr("value","");
            $("#div_db_name").hide();
            $("#db_name").attr("value","");
        }else if($("#asset_type").val() == "2"){
            $("#div_db_name").show();
            $("#div_username").show();
            $("#div_password").show();
            $("#div_protocol").hide();
            $("#protocol").attr("value","");
            $("#div_inst_name").hide();
            $("#inst_name").attr("value","");
        }else if($("#asset_type").val() == "3"){
            $("#div_inst_name").show();
            $("#div_username").show();
            $("#div_password").show();
            $("#div_protocol").hide();   
            $("#protocol").attr("value","");
            $("#div_db_name").hide();
            $("#db_name").attr("value","");
        }else if($("#asset_type").val() == "99"){
            $("#div_protocol").show();
            $("#div_inst_name").hide();
            $("#inst_name").attr("value","");
            $("#div_db_name").hide();
            $("#db_name").attr("value","");
            $("#div_username").hide();
            $("#username").attr("value","");
            $("#div_password").hide();
            $("#password").attr("value","");
            $("#div_os_type").hide();
            $("#os_type").attr("value","");
            $("#div_os_protocol").hide();
            $("#os_os_protocol").attr("value","");
            $("#div_os_username").hide();
            $("#os_username").attr("value","");
            $("#div_os_password").hide();
            $("#os_password").attr("value","");
        }

        id =  {{.dbconf.Id}};
        if (id && id > 0){
        }else{
          if(asset_type == "1"){
            $("#port").attr("value","1521");
          }else if($("#asset_type").val() == "2"){
            $("#port").attr("value","3306");
          }else if($("#asset_type").val() == "3"){
            $("#port").attr("value","1433");
          }
        }
    });  

    $("#asset_type").change(function(){
        if($("#asset_type").val() == "1"){
            $("#div_inst_name").show();
            $("#div_username").show();
            $("#div_password").show();
            $("#div_protocol").hide();   
            $("#protocol").attr("value","");
            $("#div_db_name").hide();
            $("#db_name").attr("value","");
            
            $("#port").attr("value","1521");
        }else if($("#asset_type").val() == "2"){
            $("#div_db_name").show();
            $("#div_username").show();
            $("#div_password").show();
            $("#div_protocol").hide();
            $("#protocol").attr("value","");
            $("#div_inst_name").hide();
            $("#inst_name").attr("value","");

            $("#port").attr("value","3306");
        }else if($("#asset_type").val() == "3"){
            $("#div_inst_name").show();
            $("#div_username").show();
            $("#div_password").show();
            $("#div_protocol").hide();   
            $("#protocol").attr("value","");
            $("#div_db_name").hide();
            $("#db_name").attr("value","");

            $("#port").attr("value","1433");
            $("#inst_name").attr("value","mssqlserver");
        }else if($("#asset_type").val() == "99"){
            $("#div_protocol").show();
            $("#div_inst_name").hide();
            $("#inst_name").attr("value","");
            $("#div_db_name").hide();
            $("#db_name").attr("value","");
            $("#div_username").hide();
            $("#username").attr("value","");
            $("#div_password").hide();
            $("#password").attr("value","");
            $("#div_os_type").hide();
            $("#os_type").attr("value","");
            $("#div_os_protocol").hide();
            $("#os_os_protocol").attr("value","");
            $("#div_os_username").hide();
            $("#os_username").attr("value","");
            $("#div_os_password").hide();
            $("#os_password").attr("value","");
        }
    });


    $("#protocol").change(function(){
        if($("#protocol").val() == "snmp"){
            $("#port").attr("value","161");
        }else if($("#protocol").val() == "winrm"){
            $("#port").attr("value","5985");
        }
    });
    
    function checkConnect(){
        var asset_type = $("#asset_type").val();
        var host = $("#host").val();
        var protocol = $("#protocol").val();
        var port = $("#port").val();
        var inst_name = $("#inst_name").val();
        var db_name = $("#db_name").val();
        var username = $("#username").val();
        var password = $("#password").val();
        
        $.ajax({url: "/config/db/ajax/connect",
                type: "POST",
								data: {"asset_type":asset_type, "host":host,"protocol":protocol, "port":port,"inst_name":inst_name,"db_name":db_name,"username":username,"password":password,},
                success: function (data) {
                    dialogInfo(data.message)
                    setTimeout(function(){ $('#dialogInfo').modal('hide'); }, 5000);
                    if (data.code) {
                    }
                    else {
                    }
                }
        });
    }

    $('#dbconfig-form').validate({
        ignore:'',        
		    rules : {
			      //username:{required: true},
			      role:{required: true},
        },
        messages : {
			      //username:{required: '请填写用户名'},
			      role:{required: '请选择角色'}, 
        },
        submitHandler:function(form) {
            $(form).ajaxSubmit({
                type:'POST',
                dataType:'json',
                success:function(data) {
                    dialogInfo(data.message)
                    if (data.code) {
                       		setTimeout(function(){window.location.href="/config/db/manage"}, 1000);
                    } else {
                        setTimeout(function(){ $('#dialogInfo').modal('hide'); }, 1000);
                    }
                }
            });
        }
    });
</script>
</body>
</html>
