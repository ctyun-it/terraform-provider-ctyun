# 导入SFS权限规则
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_sfs_permission_rule.[导入配置名称] [id],[permission_group_id],<region_id>
# 示例
terraform import ctyun_sfs_permission_rule.sfs_permission_rule_example rule-12345,group-12345,bb9fdb42056f11eda1610242ac110002