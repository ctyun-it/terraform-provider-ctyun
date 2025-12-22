resource "ctyun_redis_param_template" "%[1]s" {
  name         = "%[2]s"
  description  = "%[3]s"
  cache_mode   = "%[4]s"
  sys_template = %[5]t

  params = %[6]s
}

