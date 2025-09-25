resource "ctyun_ebs_backup" "%[1]s" {
  repository_id = "%[2]s"
  disk_id = "%[3]s"
  name  = "%[4]s"
  full_backup = false
}
