# 导入VPCE Service
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_vpce_service.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_vpce_service.example 12345678-1234-1234-1234-123456789012,region-123