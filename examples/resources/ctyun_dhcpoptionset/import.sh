# 导入DHCP选项集
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_dhcpoptionset.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_dhcpoptionset.dhcp_option_set_example dhcpopt-12345,region-67890