# 导入弹性公网IP绑定关系
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_eip_association.example [eip_id],<region_id>
# 示例
terraform import ctyun_eip_association.example fb5e6128-c29a-4d2d-ac87-8a5a2c8579a6,bb9fdb42056f11eda1610242ac110002