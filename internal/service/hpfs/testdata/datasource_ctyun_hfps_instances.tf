data "ctyun_hpfs_instances" "%[1]s" {
  sfs_status = "%[2]s"
  sfs_protocol = "%[3]s"
  az_name = "%[4]s"
}
