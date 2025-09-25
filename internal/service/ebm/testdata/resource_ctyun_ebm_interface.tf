resource "ctyun_ebm_interface" "%[1]s" {
  instance_id = "%[2]s"
  security_group_ids = ["%[3]s"]
  subnet_id = "%[4]s"
  ipv4 = "192.168.1.18"
  az_name = "%[5]s"
}
