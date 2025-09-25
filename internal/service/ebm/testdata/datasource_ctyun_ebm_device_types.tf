data "ctyun_ebm_device_types" "%[1]s" {
  device_type = "%[2]s"
}

data "ctyun_ebm_device_images" "%[1]s" {
  device_type = "%[2]s"
  os_type = "linux"
  image_type = "standard"
}

data "ctyun_ebm_device_raids" "%[1]s" {
  device_type = "%[2]s"
  volume_type = "system"
}
