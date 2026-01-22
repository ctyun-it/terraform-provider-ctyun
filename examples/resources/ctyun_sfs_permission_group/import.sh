# 导入SFS权限组
# [] 标记的参数为必填参数
# <> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_sfs_permission_group.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_sfs_permission_group.permission_group_example pg-123456789,<region-123456>