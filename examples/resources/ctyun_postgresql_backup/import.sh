# 导入PostgreSQL备份
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_postgresql_backup.[导入配置名称] [instance_id],[name],<region_id>
# 示例
terraform import ctyun_postgresql_backup.backup_example 20d1e9aa262246bf86b2915c6364715e,pgsql-964_20260508095105",<region_id>