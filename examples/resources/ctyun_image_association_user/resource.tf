terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 可参考index.md，在环境变量中配置ak、sk、资源池ID、可用区名称
provider "ctyun" {
  env = "prod"
}

# 将镜像分享给其他用户
resource "ctyun_image_association_user" "image_association_user_sharer_test" {
  image_id   = "9a099800-3e1c-45cd-99d1-7e2207a2fb08"
  type       = "share"
  user_email = "448725235@qq.com"
}

## 接受镜像
#resource "ctyun_image_association_user" "image_association_user_receiver_test" {
#  image_id = "9a099800-3e1c-45cd-99d1-7e2207a2fb08"
#  type     = "receive"
#}