# 导入跨域连接
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_ec_region_peer.[导入配置名称] [id],[ec_id],[packet_id],[src_cgw_id]
# 示例
terraform import ctyun_ec_region_peer.ec_region_peer_example peer-12345678,ec-87654321,pkt-11111111,cgw-22222222