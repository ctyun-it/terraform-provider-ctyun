# 导入SFS权限组关联
# [] 标记的参数为必填参数
# <> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_sfs_permission_group_association.[导入配置名称] [vpc_id],[sfs_uid],<region_id>
# 示例
terraform import ctyun_sfs_permission_group_association.permission_group_association_example vpc-123456,sfs-789012,<region-345678>