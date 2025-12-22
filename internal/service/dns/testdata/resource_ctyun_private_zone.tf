resource "ctyun_private_zone" "%[1]s" {
  name          = "%[2]s"
  description   = "%[3]s"
  proxy_pattern = "%[4]s"
  ttl           = %[5]d
  vpc_id_list   = [%[6]s]
}

