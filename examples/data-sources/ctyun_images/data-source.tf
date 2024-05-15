terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
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