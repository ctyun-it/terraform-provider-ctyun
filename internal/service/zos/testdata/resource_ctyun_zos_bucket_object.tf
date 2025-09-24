resource "ctyun_zos_bucket_object" "%[1]s" {
  bucket = "%[2]s"
  key = "%[3]s"
  source = "%[4]s"
  acl = "%[5]s"
  tags = {
    %[6]s
  }
}
