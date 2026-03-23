# 导入高可用虚IP
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_vip.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_vip.vip_example vip-12345,region-67890