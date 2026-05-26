# 导入mysql实例
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_mysql_instance.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_mysql_instance.instance_example 3dd1482933a243f9bd4e8ecb3cafbddb,bb9fdb42056f11eda1610242ac110002