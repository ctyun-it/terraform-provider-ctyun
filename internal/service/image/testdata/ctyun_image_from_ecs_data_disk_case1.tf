
# 从云主机数据盘创建私有镜像
resource "ctyun_image_from_ecs" "%[1]s" {

  image_type="data_disk"
  # 必选参数：镜像名称（2~32字符，仅数字、字母、-组成，不以数字或-开头/结尾）
  image_name = "%[2]s"

  description = "%[3]s"
  # 必选参数：云主机ID（状态需为stopped，部分资源池支持running）
  instance_id = "%[4]s"

  # 必选参数：数据盘ID（需挂载于指定云主机）
  data_disk_id = "%[5]s"

  # 可选参数：企业项目ID（默认0，即default项目）
  project_id = "0"

  # 可选参数：标签列表（最多10个，键值不可重复）
  labels = [
    {
      label_key   = "environment"
      label_value = "production"
    },
    {
      label_key   = "department"
      label_value = "devops"
    }
  ]

  # 可选参数：是否启用镜像完整性校验（默认false，仅部分资源池支持）
  enable_image_integrity_check = false
}