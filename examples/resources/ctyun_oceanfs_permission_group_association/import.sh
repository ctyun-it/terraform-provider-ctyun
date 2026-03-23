# 导入OceanFS权限组关联
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_oceanfs_permission_group_association.[导入配置名称] [oceanfs_id],[vpc_id],[permission_group_id],<region_id>
# 示例
terraform import ctyun_oceanfs_permission_group_association.permission_group_association_example oceanfs-123,vpc-456,perm-789,region-bb9fdb42056f11eda1610242ac110002