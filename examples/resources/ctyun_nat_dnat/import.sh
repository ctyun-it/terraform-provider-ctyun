# 导入dnat规则
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_nat_dnat.[导入配置名称] [id],[nat_gateway_id],<region_id>
# 示例
terraform import ctyun_nat_dnat.dnat_example 376f2f85-ff34-c4e0-4f5b-320dd427a271,nat-12345678,bb9fdb42056f11eda1610242ac110002