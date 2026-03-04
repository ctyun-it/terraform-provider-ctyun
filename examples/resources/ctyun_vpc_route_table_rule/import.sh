# 导入VPC路由表规则
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_vpc_route_table_rule.[导入配置名称] [id],[route_table_id],<region_id>
# 示例
terraform import ctyun_vpc_route_table_rule.route_table_rule_example rtr-12345,rtb-67890,region-11111