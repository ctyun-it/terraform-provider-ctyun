# 导入共享带宽与弹性IP关联
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_bandwidth_association_eip.[导入配置名称] [bandwidth_id],[eip_id],<region_id>
# 示例
terraform import ctyun_bandwidth_association_eip.bandwidth_eip_assoc_example bandwidth-12345,eip-67890,region-11111