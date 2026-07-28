resource "ctyun_private_zone_record" "%[1]s" {
  zone_id     = "%[2]s"
  type        = "%[3]s"
  value_list = [%[4]s]
  ttl         = %[5]d
  name        = "%[6]s"
  description = "%[7]s"
  enabled     = %[8]t
}

