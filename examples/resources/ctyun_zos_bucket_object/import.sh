# 导入ZOS Bucket Object
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_zos_bucket_object.[导入配置名称] [key],[bucket],<region_id>
# 示例
terraform import ctyun_zos_bucket_object.example my-object-key,my-bucket-name,region-123