# 导入mongodb实例
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_mongodb_instance.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_mongodb_instance.mongodb_instance_example f9f1b75c219d4b1194011b088486a5f7