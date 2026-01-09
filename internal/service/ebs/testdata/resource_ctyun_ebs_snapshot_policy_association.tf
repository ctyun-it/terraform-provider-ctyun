resource "ctyun_ebs_snapshot_policy_association" "%[1]s" {
  snapshot_policy_id = %[2]s
  disk_id = "%[3]s"
}
