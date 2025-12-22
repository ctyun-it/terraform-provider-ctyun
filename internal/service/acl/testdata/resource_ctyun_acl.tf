resource "ctyun_acl" "%[1]s" {
  project_id  = "%[2]s"
  vpc_id      = "%[3]s"
  name        = "%[4]s"
  description = "%[5]s"
}
