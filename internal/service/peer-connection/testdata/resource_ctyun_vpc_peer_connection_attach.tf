

provider "ctyun" {
  alias           = "test_accpet"
  ak              = ""                                    # 如果此值不填，则默认读取环境变量中的CTYUN_AK
  sk              = ""                                    # 如果此值不填，则默认读取环境变量中的CTYUN_SK
  env             = "prod"                                     # 如果此值不填，则默认读取环境变量中的CTYUN_ENV
}

resource "ctyun_vpc_peer_connection_attach" "%[1]s" {
  provider = ctyun.test_accpet
  peer_connection_id = %[2]s
  operation          = "%[3]s"
}

