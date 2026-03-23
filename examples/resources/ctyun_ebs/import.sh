# 导入云硬盘(ebs)
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_ebs.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_ebs.example d-ebs1234567890,bb9fdb42056f11eda1610242ac110002