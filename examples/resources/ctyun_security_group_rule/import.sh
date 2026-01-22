# 导入安全组规则
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_security_group_rule.[导入配置名称] [id],[security_group_id],<region_id>
# 示例
terraform import ctyun_security_group_rule.security_group_rule_example sgr-12345,sg-67890,region-11111