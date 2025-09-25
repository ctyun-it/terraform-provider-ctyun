resource "ctyun_security_group" "%[1]s" {
  vpc_id = "%[4]s"
  name        = "%[2]s"
  description = "%[3]s"
}
