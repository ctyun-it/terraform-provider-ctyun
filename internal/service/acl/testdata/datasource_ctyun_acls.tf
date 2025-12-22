data "ctyun_acls" "%[1]s" {
  id         = %[2]s
  project_id = "%[3]s"
  name       = "%[4]s"
  page_no    = %[5]d
  page_size  = %[6]d
}
