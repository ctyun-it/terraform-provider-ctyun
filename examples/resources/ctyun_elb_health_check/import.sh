# 导入负载均衡健康检查
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_elb_health_check.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_elb_health_check.elb_health_check_example 376f2f85-ff34-c4e0-4f5b-320dd427a271,<region_id>