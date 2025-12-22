resource "ctyun_vip_association" "%[1]s" {
  vip_id             = "%[2]s"
  resource_type        = "%[3]s"
  network_interface_id = "%[4]s"
  instance_id          = "%[5]s"
  floating_id          = "%[6]s"
}