# 导入VPCE Service Transit IP
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_vpce_service_transit_ip.[导入配置名称] [transit_ip],[endpoint_service_id],<region_id>
# 示例
terraform import ctyun_vpce_service_transit_ip.example 192.168.1.100,service-123,region-123