resource "ctyun_elb_loadbalancer" "%[1]s" {
  subnet_id     = "%[2]s"
  name          = "%[3]s"
  sla_name      = "%[4]s"
  resource_type = "%[5]s"
  vpc_id        = "%[6]s"
  description   = "%[7]s"
  %[8]s
  %[9]s
  %[10]s
}

