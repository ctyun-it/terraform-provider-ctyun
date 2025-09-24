resource "ctyun_zos_bucket" "%[1]s" {
  bucket = "%[2]s"
  acl = "%[3]s"
  az_policy = "%[4]s"
  storage_type = "%[5]s"
  is_encrypted = true
}
