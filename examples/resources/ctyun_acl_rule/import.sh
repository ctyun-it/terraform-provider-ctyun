# 导入ACL规则
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_acl_rule.[导入配置名称] [id],[acl_id],<region_id>
# 示例
terraform import ctyun_acl_rule.acl_rule_example 376f2f85-ff34-c4e0-4f5b-320dd427a271,acl-123456,bb9fdb42056f11eda1610242ac110002