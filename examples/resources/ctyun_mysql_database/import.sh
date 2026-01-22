# 导入MySQL数据库资源
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_mysql_database.[导入配置名称] [name],[instance_id],<project_id>,<region_id>
# 示例
terraform import ctyun_mysql_database.example testdb,c81c2dc376e34e7887334cbcbd4xxx,prj-1234567890abcdef0,bb9fdb42056f11eda1610242ac110002