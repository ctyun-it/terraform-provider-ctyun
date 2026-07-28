# 导入mysql参数模板
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_mysql_param_template.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_mysql_param_template.param_template_example 837xhs18jd7bctgs6f5b320dd427a271,bb9fdb42056f11eda1610242ac110002