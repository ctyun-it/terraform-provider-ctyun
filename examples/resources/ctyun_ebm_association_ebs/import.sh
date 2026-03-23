# 导入物理机云硬盘关联
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_ebm_association_ebs.[导入配置名称] [instance_id],[ebs_id],<az_name>,<region_id>
# 示例
terraform import ctyun_ebm_association_ebs.ebm_association_ebs_example 376f2f85-ff34-c4e0-4f5b-320dd427a271,f2d1c4e5-1a2b-3c4d-5e6f-7a8b9c0d1e2f,cn-zj-hgh7-1a-public-ctcloud,200000003329