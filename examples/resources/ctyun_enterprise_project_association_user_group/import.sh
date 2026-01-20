# 导入企业项目关联用户组资源
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_enterprise_project_association_user_group.[导入配置名称] [enterprise_project_id],[user_group_id]
# 示例
terraform import ctyun_enterprise_project_association_user_group.enterprise_project_user_group_example 376f2f85-ff34-c4e0-4f5b-320dd427a271,group1234567890