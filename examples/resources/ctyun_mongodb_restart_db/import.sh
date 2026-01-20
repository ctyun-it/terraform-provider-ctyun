# 导入mongodb重启db
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_mongodb_restart_db.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_mongodb_restart_db.mongodb_restart_db_example 12345678-1234-1234-1234-123456789012,bb9fdb42056f11eda1610242ac110002