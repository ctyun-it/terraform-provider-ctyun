resource "ctyun_bandwidth" "%[1]s" {
  name       = "%[2]s"
  bandwidth  = "%[3]s"
  cycle_count = 1
  cycle_type  = "month"
}
