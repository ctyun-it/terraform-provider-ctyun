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

resource "ctyun_image" "image_test" {
  name         = "image-test2"
  file_source  = "https://jiangsu-10.zos.ctyun.cn/bucket-305c/test1/Ubuntu-22.04-x86_64-231229-R3.qcow2"
  os_distro    = "ubuntu"
  os_version   = "22.04"
  architecture = "x86_64"
  boot_mode    = "bios"
  description  = "测试镜像上传1"
  disk_size    = "8"
  type         = "system"
  project_id   = "4f5ef15300724760af59b37cf6409f45"
}
