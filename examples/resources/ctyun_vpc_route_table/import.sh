# 导入VPC路由表
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_vpc_route_table.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_vpc_route_table.route_table_example rtb-12345,region-67890