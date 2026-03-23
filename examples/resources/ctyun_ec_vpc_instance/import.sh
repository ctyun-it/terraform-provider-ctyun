# 导入VPC网络实例
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_ec_vpc_instance.[导入配置名称] [id],[ec_id],[cgw_id]
# 示例
terraform import ctyun_ec_vpc_instance.ec_vpc_instance_example vpc-ins-12345678,ec-87654321,cgw-11111111