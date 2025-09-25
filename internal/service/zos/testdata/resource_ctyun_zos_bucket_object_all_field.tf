resource "ctyun_zos_bucket_object" "%[1]s" {
  bucket = "%[2]s"
  key = "%[3]s"
  content = "%[4]s"
  acl = "%[5]s"
  tags = {
    %[6]s
  }
  storage_type = "%[7]s"
  cache_control = "%[8]s"
  content_encoding = "%[9]s"
  content_type = "%[10]s"
}
