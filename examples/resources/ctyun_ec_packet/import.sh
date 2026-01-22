# 导入云间高速带宽包
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_ec_packet.[导入配置名称] [ec_id],[packet_id]
# 示例
terraform import ctyun_ec_packet.ec_packet_example ec-12345678,pkt-87654321