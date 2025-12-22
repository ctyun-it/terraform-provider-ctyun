
provider "ctyun" {
  alias           = "default"
  ak              = ""                                    # 如果此值不填，则默认读取环境变量中的CTYUN_AK
  sk              = ""
  env             = "prod"
}

resource "ctyun_vpc_peer_connection" "%[1]s" {
  # provider = ctyun.default
  project_id     = "%[2]s"
  name           = "%[3]s"
  description    = "%[4]s"
  request_vpc_id = "%[5]s"
  accept_vpc_id  = "%[6]s"
  accept_email   = "%[7]s"
}

