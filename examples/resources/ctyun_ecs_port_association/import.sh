# 导入弹性网卡关联
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_ecs_port_association.[导入配置名称] [instance_id],[port_id],<region_id>
# 示例
terraform import ctyun_ecs_port_association.port_association_example instance-12345678,port-87654321,<region_id>