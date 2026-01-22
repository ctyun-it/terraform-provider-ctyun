# 导入MySQL实例与EIP关联资源
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_mysql_association_eip.[导入配置名称] [instance_id],[eip_id],<project_id>,<region_id>
# 示例
terraform import ctyun_mysql_association_eip.example ce5686a3-c557-45ec-8e9a-50981f552580,eip-1234567890abcdef0,prj-1234567890abcdef0,bb9fdb42056f11eda1610242ac110002