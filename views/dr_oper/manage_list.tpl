<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8" />
<title>{{config "String" "globaltitle" ""}}</title>
{{template "inc/meta.tpl" .}}
</head>
<body class="sticky-header">
<section>
  {{template "inc/left.tpl" .}}
  <!-- main content start-->
  <div class="main-content">
    <!-- header section start-->
    <div class="header-section">
      <a class="toggle-btn"><i class="fa fa-bars"></i></a>
      <!--search start-->
      <!--search end-->
      {{template "inc/user-info.tpl" .}}
    </div>
    <!-- header section end-->
    <!-- page heading start-->
    <div class="page-heading">
      <!-- <h3> 日志管理 </h3>-->
      <ul class="breadcrumb pull-left">
        <li><a href="/operation/dr_manage/list">容灾管理</a></li>
      </ul>
    </div>
    <!-- page heading end-->
    <!--body wrapper start-->
    <div class="wrapper">
      <div class="row">
        <div class="col-sm-12">
          <!-- 主体内容 开始 -->
            <div class="searchdiv">
              <div class="search-form">
                <div class="form-inline">
                  <div class="form-group">
                    <form method="get">
                    <input type="text" name="search_name" placeholder="请输入名称" class="form-control" value="{{.condArr.search_name}}"/>
                    <button class="btn btn-primary" type="submit"> <i class="fa fa-search"></i> 搜索 </button>
                    <a href="/operation/dr_manage/list" class="btn btn-default" type="submit"> <i class="fa fa-reset"></i> 重置 </a>
                    </form>
                  </div>
                </div>
              </div>
            </div>

            <section class="panel">
              <header class="panel-heading"> 容灾列表 / 总数：{{.countBs}}
                <span class="tools pull-right"><a href="javascript:;" class="fa fa-chevron-down"></a>
                <!--a href="javascript:;" class="fa fa-times"></a-->
                </span> 
              </header>
              <div class="panel-body">
                <section id="unseen">
                  <form id="user-form-list">
                    <table class="table table-bordered table-striped table-condensed">
                      <thead>
                        <tr>
                          <th>数据库类型</th>
                          <th>容灾组名称</th>
                          <th>主库</th>
                          <th>主库时间</th>
                          <th>备库</th>
                          <th>备库时间</th>
                          <th>同步状态</th>
                          <th>延时</th>
                          <th>操作</th>
                        </tr>
                      </thead>
                      <tbody>
                      {{range $k,$v := .dr}}
                        <tr>
                          <td>{{ getAssetImage $v.Db_Type_P | str2html}}</td>
                          <td>{{$v.Bs_Name}}</td>
                          <td>{{if eq "" $v.Host_P}}---{{else}}{{if eq 0 $v.Is_Switch}}{{$v.Host_P}}:{{$v.Port_P}}{{else}}{{$v.Host_S}}:{{$v.Port_S}}{{end}}{{end}}</td>
                          <td>{{$v.Db_Time_P}}</td>
						  <td>{{if eq "" $v.Host_P}}---{{else}}{{if eq 0 $v.Is_Switch}}{{$v.Host_S}}:{{$v.Port_S}}{{else}}{{$v.Host_P}}:{{$v.Port_P}}{{end}}{{end}}</td>
                          <td>{{$v.Db_Time_S}}</td>
						  <td>{{checkDrStatusLevel $v.Repl_Status | str2html}}</td>
						  <td>{{if eq "unKnown" $v.Repl_Status}}--{{else}}{{$v.Repl_Delay}}{{end}}</td>
						  <td>
						  	<a name="screen"  {{if eq 1 $v.Db_Type_P}}class="btn btn-primary" href="/operation/screen/view/{{$v.Id}}"{{else}}class="btn btn-default"{{end}}> <i class="fa fa-reset"></i> 灾备大屏 </a> 
                            <a name="btnDrManage" class="btn btn-primary" type="button" href="/operation/dr_manage/detail/{{$v.Id}}"> <i class="fa fa-reset"></i> 容灾管理 </a>
                            <!--<button name="switchover" class="btn btn-primary" type="button" value="Switchover" onclick="checkUser(this)" data-id="{{$v.Id}}"> <i class="fa fa-reset"></i> 维护切换 </button>
                            <button name="failover" class="btn btn-danger" type="button" value="Failover" onclick="checkUser(this)" data-id="{{$v.Id}}"> <i class="fa fa-reset"></i> 灾难切换 </button>-->
                          </td>
                        </tr>
                      {{end}}
                      </tbody>
                    </table>
                  </form>
                  {{template "inc/page.tpl" .}}
                </section>
              </div>
            </section>
          <!-- 主体内容 结束 -->
        </div>
      </div>
    </div>
    <!--body wrapper end-->
    <!--footer section start-->
    {{template "inc/foot-info.tpl" .}}
    <!--footer section end-->
  </div>
  <!-- main content end-->
  <div id="div_layer" style="display:none" ></div>
</section>

{{template "inc/foot.tpl" .}}    
<script>

var mylay = null;
var oTimer = null; 

    
var user_pwd = {{.user.Password}} ;
var div_layer = document.getElementById("div_layer");
var query_url="/operation/dr_switch/process";
var bs_id = -1;


function checkUser(e){
    bs_id = $(e).attr('data-id');

		if(e.value == "Switchover"){
			_message = "确认要开始维护切换吗？";
      		target_url = "/operation/dr_switch/switchover";
			op_type = "SWITCHOVER";
		}
		else if(e.value == "Failover"){
			_message = "确认要开始灾难切换吗？";
      		target_url = "/operation/dr_switch/failover";
			op_type = "FAILOVER";
		}
		else{
			return;
		}	
    

		bootbox.prompt({
		    title: "请确认密码",
		    inputType: 'password',
			buttons: {confirm: {label: "确认"}, cancel: {label: "取消"} },
		    callback: function (result) {
		    	if(result)
		    	{ 
		        if (md5(result) == user_pwd)
		        {
					bootbox.dialog({
						message: _message,
						buttons: {
							cancel: {
								label: '取消',
								className: 'btn-default',
								callback: function () {
								}
							},
							ok: {
								label: '确定',
								className: 'btn-danger',
								callback: function(){
									$.ajax({url: target_url,
											type: "POST",
											data: {"bs_id":bs_id,"op_type":op_type},
											success: function (json) {
												//回调函数，判断提交返回的数据执行相应逻辑
												if (json.code == 0) {
													bootbox.alert({
													message: json.message,
													buttons: {
															ok: {
															label: '确定',
															className: 'btn-success'
															}
														},
														callback: function () {
														window.location.reload();
													}
													});
															
													if(mylay!=null){
													layer.close(mylay);
													}
													clearInterval(oTimer);
												}
												else {
												}
											}
											});

												
											$('#div_layer').html("");			//初始化div
											mylay = layer.open({
												type: 1,
												skin: 'layui-layer-demo layblack', //样式类名
												closeBtn: 0, //不显示关闭按钮
												anim: 1,
												title: '详细过程',
												area: ['450px', '240px'],
												shadeClose: false, //开启遮罩关闭
												content: $('#div_layer')
											});
											
											oTimer = setInterval("queryHandle(query_url, bs_id, op_type)",2000);
								}
							}
						}
					});
		        }
		        else
		        {
		        	bootbox.alert({
		        		message: "密码不对，请确认后重新尝试!",
		        		buttons: {
							        ok: {
							            label: '确定',
							            className: 'btn-success'
							        }
							    }
		        	});
		        }
		      }
		
		    }
		});

}

// 初始化内容
$(function(){
		
});  
  
function queryHandle(url, bs_id, op_type){
    $.post(url, {"bs_id":bs_id, "op_type":op_type}, function(json){
        if(json.on_process == 0){
        		if(json.op_type != ""){
		        		//alert(json.op_result);
		        		
		        		if(json.op_type == "SWITCHOVER"){
							if(json.op_reason == 'null'){
								error_message = "主备切换失败，详细原因请查看相关日志";
							}else{
								error_message = "主备切换失败，原因是：" + json.op_reason;
							}
							
							ok_message = "主备切换成功";
		        		}else if(json.op_type == "FAILOVER"){
							if(json.op_reason == 'null'){
								error_message = "灾难切换失败，详细原因请查看相关日志";
							}else{
								error_message = "灾难切换失败，原因是：" + json.op_reason;
							}
							
							ok_message = "灾难切换成功";
		        		}
        		
        				if(json.op_result == '-1'){
							bootbox.alert({
								message: error_message,
								buttons: {
											ok: {
												label: '确定',
												className: 'btn-success'
											}
										},
									callback: function () {
										window.location.reload();
								}
							});
						        	
							if(mylay!=null){
								layer.close(mylay);
							}
							clearInterval(oTimer);
						        	
        				}else if(json.op_result == '1'){
							bootbox.alert({
									message: ok_message,
									buttons: {
												ok: {
													label: '确定',
													className: 'btn-success'
												}
											},
										callback: function () {
											window.location.reload();
									}
							});
						        	
							if(mylay!=null){
								layer.close(mylay);
							}
							clearInterval(oTimer); 
        				}else{
							if(mylay!=null){
								layer.close(mylay);
							}
							clearInterval(oTimer); 
						}
        		}
        }else if(json.on_process == -1){
            bootbox.alert({
                message: "该系统没有配置容灾库",
                buttons: {
                      ok: {
                        label: '确定',
                        className: 'btn-success'
                      }
                    },
                  callback: function () {
                    window.location.reload();
                }
            });
                    
            if(mylay!=null){
              layer.close(mylay);
            }
            clearInterval(oTimer); 
        }

		if(mylay!=null){
			localJson = $.parseJSON(json.json_process);
			//alert(localJson);
			$("#div_layer").empty();
			$.each(localJson,function(idx,item){   
				//alert("Time:"+item.Time+",Process_desc:"+item.Process_desc);   
        		$("#div_layer").append("<p>" + item.Time + ": " + item.Process_desc + "</p>");
			});  

        	$(".layui-layer-content").scrollTop($(".layui-layer-content")[0].scrollHeight);
        }  
    },'json');  
}

</script>
</body>
</html>
