# 导入DHCP选项集与VPC关联
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_dhcpoptionset_association_vpc.[导入配置名称] [dhcp_option_sets_id],<region_id>
# 示例
terraform import ctyun_dhcpoptionset_association_vpc.dhcp_vpc_assoc_example dhcpopt-12345,region-67890