# 导入弹性IP
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_eip.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_eip.eip_example eip-12345,region-67890