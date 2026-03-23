# 导入VPCE Service Connection
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_vpce_service_connection.[导入配置名称] [endpoint_service_id],[endpoint_id],<region_id>
# 示例
terraform import ctyun_vpce_service_connection.example service-123,endpoint-123,region-123