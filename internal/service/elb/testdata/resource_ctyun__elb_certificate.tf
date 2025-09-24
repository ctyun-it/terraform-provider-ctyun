resource "ctyun_elb_certificate" "%[1]s" {
  name        = "%[2]s"
  type        = "%[3]s"
  certificate = %[4]s
  %[5]s
}
