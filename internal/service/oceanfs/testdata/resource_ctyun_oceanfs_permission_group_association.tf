resource "ctyun_oceanfs_permission_group_association" "%[1]s" {
  permission_group_id = "%[2]s"
  oceanfs_id              = "%[3]s"
  vpc_id              = "%[4]s"
}

