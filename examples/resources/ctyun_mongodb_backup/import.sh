# 导入mongodb备份
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_mongodb_backup.[导入配置名称] [instance_id],[name],<region_id>
# 示例
# 使用全部参数
terraform import ctyun_mongodb_backup.example f9f1b75c219d4b1194011b088486a5f7,backup-name,region-11111111