# 导入SD-WAN ACL规则
# [] 标记的参数为必填参数
# <> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_sdwan_acl_rule.[导入配置名称] [acl_id],[id]
# 示例
terraform import ctyun_sdwan_acl_rule.sdwan_acl_rule_example acl-123456,aclr-789012