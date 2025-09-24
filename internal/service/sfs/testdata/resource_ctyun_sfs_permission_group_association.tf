resource "ctyun_sfs_permission_group_association" "%[1]s" {
  permission_group_fuid = "%[2]s"
  sfs_uid               = "%[3]s"
  vpc_id                = "%[4]s"
}


