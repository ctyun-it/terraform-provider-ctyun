# 导入VPCE Service Reverse Rule
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_vpce_service_reverse_rule.[导入配置名称] [id],[endpoint_service_id],<region_id>
# 示例
terraform import ctyun_vpce_service_reverse_rule.example rule-123,service-123,region-123