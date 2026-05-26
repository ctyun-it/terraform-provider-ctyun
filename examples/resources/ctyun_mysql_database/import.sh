# 导入MySQL数据库资源
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_mysql_database.[导入配置名称] [instance_id],[name],<region_id>
# 示例
terraform import ctyun_mysql_database.example 3dd1482933a243f9bd4e8ecb3cafbddb,db1,bb9fdb42056f11eda1610242ac110002