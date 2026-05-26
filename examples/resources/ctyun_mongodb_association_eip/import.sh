# 导入mongodb关联eip
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_mongodb_association_eip.[导入配置名称] [instance_id],[eip_id],<region_id>
# 示例
terraform import ctyun_mongodb_association_eip.mongodb_association_eip_example ff532dfa5e3744928bcb16daf50b4b69,eip-qjefnklk36,region-11111111