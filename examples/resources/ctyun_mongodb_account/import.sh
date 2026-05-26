# 导入mongodb账号
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_mongodb_account.[导入配置名称] [instance_id],[name],<region_id>
# 示例
terraform import ctyun_mongodb_account.mongodb_account_example ff532dfa5e3744928bcb16daf50b4b69,user_example