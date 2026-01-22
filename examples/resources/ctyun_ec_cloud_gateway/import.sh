# 导入云网关实例
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_ec_cloud_gateway.[导入配置名称] [ec_id],[cgw_id]
# 示例
terraform import ctyun_ec_cloud_gateway.ec_cloud_gateway_example ec-12345678,cgw-87654321