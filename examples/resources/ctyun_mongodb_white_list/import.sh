# 导入mongodb白名单
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_mongodb_white_list.[导入配置名称] [instance_id],[group_name],<region_id>
# 示例
# 使用必填参数（region_id从环境变量获取）
terraform import ctyun_mongodb_white_list.example inf9f1b75c219d4b1194011b088486a5f7,my-whitelist-group
