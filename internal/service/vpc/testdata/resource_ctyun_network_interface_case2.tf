resource "ctyun_port" "%[1]s" {
  name                       = "%[2]s"
  description                = "%[3]s"
  subnet_id                  = "%[4]s"
  primary_ip_address         = "%[5]s"
  security_group_ids         = ["%[6]s"]
  secondary_private_ip_count = %[7]d
  ipv6_address_count         = %[8]d
}
