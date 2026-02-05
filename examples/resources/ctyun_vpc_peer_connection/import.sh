# 导入VPC对等连接
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_vpc_peer_connection.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_vpc_peer_connection.vpc_peer_connection_example vpr-12345678,200000001852