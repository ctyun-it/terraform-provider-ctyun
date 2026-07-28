resource "ctyun_vpc_peer_connection" "%[1]s" {
  name           = "%[2]s"
  description    = "%[3]s"
  request_vpc_id = "%[4]s"
  accept_vpc_id  = "%[5]s"
  tags           = %[6]s
}

