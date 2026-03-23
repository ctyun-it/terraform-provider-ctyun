# 导入密钥对
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_keypair.[导入配置名称] [name],<region_id>
# 示例
terraform import ctyun_keypair.keypair_example my-keypair,<region_id>