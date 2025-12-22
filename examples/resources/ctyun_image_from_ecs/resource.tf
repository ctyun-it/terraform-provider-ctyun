terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

provider "ctyun" {
  env = "prod"
}


# 从云主机整机创建私有镜像
resource "ctyun_image_from_ecs" "entire_machine" {
  # 整机镜像配置
  image_type = "entire_machine"
  # 必选参数：镜像名称（2~32字符，仅数字、字母、-组成，不以数字或-开头/结尾）
  image_name = "entire-machine-image-update1"
  # 必选参数：云主机ID（状态需为running或stopped，至少有1块数据盘）
  description = "terrform tf文件命令测试-更新1"
  instance_id = "ae432721-61bf-45b7-b207-7e3256c1c2d6"
  # 可选参数：云主机备份存储库ID（非多可用区资源池时必填）
  # repository_id = "repo-1234567890abcdef"

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




resource "ctyun_image_from_ecs" "data_disk" {

  image_type = "data_disk"
  # 必选参数：镜像名称（2~32字符，仅数字、字母、-组成，不以数字或-开头/结尾）
  image_name = "data-disk-image-update"

  description = "tf 脚本方式测试-更新"
  # 必选参数：云主机ID（状态需为stopped，部分资源池支持running）
  instance_id = "ae432721-61bf-45b7-b207-7e3256c1c2d6"

  # 必选参数：数据盘ID（需挂载于指定云主机）
  data_disk_id = "ae432721-61bf-45b7-b207-7e3256c1c2d6"

  # 可选参数：企业项目ID（默认0，即default项目）
  project_id = "0"

  # 可选参数：标签列表（最多10个，键值不可重复）
  # labels = [
  #   {
  #     label_key   = "environment"
  #     label_value = "production"
  #   },
  #   {
  #     label_key   = "department"
  #     label_value = "devops"
  #   }
  # ]

  # 可选参数：是否启用镜像完整性校验（默认false，仅部分资源池支持）
  enable_image_integrity_check = false
}


resource "ctyun_image_from_ecs" "system_disk" {
  # 必选参数：镜像名称（2~32字符，仅数字、字母、-组成，不以数字或-开头/结尾）

  image_type = "system_disk"
  image_name = "system-disk-image-update"

  # 必选参数：云主机ID（状态需为stopped，部分资源池支持running）
  description = "系统盘镜像 脚本方式测试-更新"
  instance_id = "ae432721-61bf-45b7-b207-7e3256c1c2d6"


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

  # 可选参数：最小内存限制（GiB，0表示不限制，取值：0/1/2/4/8/16/32/64/128/256/512）
  minimum_ram = 4

  # 可选参数：最大内存限制（GiB，需≥最小内存，取值同上）
  maximum_ram = 32
}
