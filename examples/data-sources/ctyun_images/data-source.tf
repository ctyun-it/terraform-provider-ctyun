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

data "ctyun_images" "ctyun_images_test" {
  name       = "Ubuntu 22.04"
  visibility = "public"
  page_size  = 50
  page_no    = 1
}

output "ctyun_image" {
  value = data.ctyun_images.ctyun_images_test.images
}