resource "ctyun_subnet_association_acl" "%[1]s" {
  project_id = "%[2]s"
  acl_id     = "%[3]s"
  subnet_id  = "%[4]s"
}
