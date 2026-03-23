# 导入物理机网卡
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_ebm_interface.[导入配置名称] [instance_id],[interface_id],<az_name>,<region_id>
# 示例
terraform import ctyun_ebm_interface.interface_example 376f2f85-ff34-c4e0-4f5b-320dd427a271,if-123456789,cn-zj-hgh7-1a-public-ctcloud,200000003329