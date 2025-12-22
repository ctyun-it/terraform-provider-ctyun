resource "ctyun_vpc_peer_connection" "%[1]s" {
  project_id     = "%[2]s"
  name           = "%[3]s"
  description    = "%[4]s"
  request_vpc_id = "%[5]s"
  accept_vpc_id  = "%[6]s"

  # lifecycle {
  #   ignore_changes = [id]
  # }
}
