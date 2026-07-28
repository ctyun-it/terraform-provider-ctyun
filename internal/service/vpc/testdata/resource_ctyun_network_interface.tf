

resource "ctyun_port" "%[1]s" {
  name                       = "%[2]s"
  description                = "%[3]s"
  subnet_id                  = "%[4]s"
  security_group_ids        = ["%[5]s"]
  secondary_private_ip_count = 1
}
