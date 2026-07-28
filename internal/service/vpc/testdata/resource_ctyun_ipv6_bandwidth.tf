resource "ctyun_ipv6_bandwidth" "%[1]s" {
  name = "%[2]s"
  bandwidth = "%[3]s"
  cycle_type = "month"
  cycle_count = 1
}