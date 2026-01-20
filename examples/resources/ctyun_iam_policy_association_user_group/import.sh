# 导入IAM策略关联用户组资源
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_iam_policy_association_user_group.[导入配置名称] [id]
# 示例
terraform import ctyun_iam_policy_association_user_group.policy_group_example 376f2f85-ff34-c4e0-4f5b-320dd427a271