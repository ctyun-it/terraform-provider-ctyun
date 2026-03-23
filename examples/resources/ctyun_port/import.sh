# 导入网络接口/弹性网卡
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_port.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_port.port_example port-12345,region-67890