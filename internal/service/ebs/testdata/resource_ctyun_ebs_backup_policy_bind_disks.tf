resource "ctyun_ebs_backup_policy_bind_disks" "%[1]s" {
  policy_id = %[2]s
  disk_id_list = "%[3]s"
}
