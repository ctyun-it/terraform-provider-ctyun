# 导入私网NAT中转IP
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_private_nat_transit_ip.[导入配置名称] [address],[nat_gateway_id],<region_id>
# 示例
terraform import ctyun_private_nat_transit_ip.transit_ip_example 192.168.1.100,nat-67890,region-bb9fdb42056f11eda1610242ac110002