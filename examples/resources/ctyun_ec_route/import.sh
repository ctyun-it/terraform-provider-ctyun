# 导入云间高速路由
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_ec_route.[导入配置名称] [id],[ec_id],[cgw_id],[rtb_id]
# 示例
terraform import ctyun_ec_route.ec_route_example route-12345678,ec-87654321,cgw-11111111,rtb-22222222