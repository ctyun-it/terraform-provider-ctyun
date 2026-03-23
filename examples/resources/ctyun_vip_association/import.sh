# 导入高可用虚IP关联
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_vip_association.[导入配置名称] [vip_id],[network_interface_id],<region_id>
或
terraform import ctyun_vip_association.[导入配置名称] [vip_id],[floating_id],<region_id>
# 示例
terraform import ctyun_vip_association.vip_assoc_example vip-12345,port-xxxx,region-22222
或
terraform import ctyun_vip_association.vip_assoc_example vip-12345,eip-xxxx,region-22222