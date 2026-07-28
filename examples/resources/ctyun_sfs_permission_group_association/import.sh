# 导入SFS权限组关联
# [] 标记的参数为必填参数
# <> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_sfs_permission_group_association.[导入配置名称] [sfs_id],[vpc_id],[permission_group_id],<region_id>
# 示例
terraform import ctyun_sfs_permission_group_association.permission_group_association_example sfs-789012,vpc-123456,xxxxxxxx,<region-345678>