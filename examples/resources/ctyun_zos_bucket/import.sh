# 导入ZOS Bucket
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_zos_bucket.[导入配置名称] [bucket],<region_id>
# 示例
terraform import ctyun_zos_bucket.example my-bucket-name,region-123