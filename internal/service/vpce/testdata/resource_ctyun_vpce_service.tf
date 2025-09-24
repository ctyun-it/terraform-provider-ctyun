resource "ctyun_vpce_service" "%[1]s" {
  name  = "%[2]s"
  vpc_id = "%[5]s"
  subnet_id = "%[6]s"
  auto_connection = true
  type = "interface"
  instance_id = "%[7]s"
  instance_type = "vm"
  rules = [{
    protocol = "TCP"
    endpoint_port = %[3]s
    server_port = 1
  }]
  whitelist_email = [
    "%[4]s",
  ]
}
