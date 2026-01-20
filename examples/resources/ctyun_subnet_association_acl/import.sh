# 导入子网ACL关联
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_subnet_association_acl.[导入配置名称] [subnet_id],[acl_id],<region_id>
# 示例
terraform import ctyun_subnet_association_acl.subnet_association_acl_example subnet-123456,acl-123456,bb9fdb42056f11eda1610242ac110002