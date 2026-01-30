# 导入SDWAN网络实例
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_ec_sdwan_instance.[导入配置名称] [sdwan_id],[ec_id],[cgw_id]
# 示例
terraform import ctyun_ec_sdwan_instance.ec_sdwan_instance_example sdwan-12345678,ec-11111111,cgw-87654321