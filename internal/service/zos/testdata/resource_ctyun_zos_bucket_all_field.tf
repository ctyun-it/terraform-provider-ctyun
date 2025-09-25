resource "ctyun_zos_bucket" "%[1]s" {
  bucket = "%[2]s"
  acl = "%[3]s"
  az_policy = "%[4]s"
  storage_type = "%[5]s"
  is_encrypted = true
  tags = %[6]s
  version_enabled = true
  log_enabled = true
  log_bucket = "%[7]s"
  log_prefix = "%[8]s"
  retention_mode = "COMPLIANCE"
  %[9]s
}

