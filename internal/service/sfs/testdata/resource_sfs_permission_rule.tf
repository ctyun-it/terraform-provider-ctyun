resource "ctyun_sfs_permission_rule" "%[1]s" {
  permission_group_id    = "%[2]s"
  auth_addr                = "%[3]s"
  rw_permission            = "%[4]s"
  priority = %[5]d
}


