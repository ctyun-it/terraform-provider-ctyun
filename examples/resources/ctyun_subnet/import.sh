# 导入子网
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_subnet.example [id],<region_id>
# 示例
terraform import ctyun_subnet.example 4a0a1e86-0736-4c33-9478-359a1307a2c8,bb9fdb42056f11eda1610242ac110002