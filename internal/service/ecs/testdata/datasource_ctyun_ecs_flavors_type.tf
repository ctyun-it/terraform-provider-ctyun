data "ctyun_ecs_flavors" "%[1]s" {
  cpu    = %[2]d
  ram    = %[3]d
  arch   = "%[4]s"
  series = "%[5]s"
  type   = "%[6]s"
}
