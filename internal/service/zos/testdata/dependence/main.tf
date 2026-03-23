resource ctyun_zos_subscribe "test" {

}

resource "ctyun_zos_bucket" "test" {
  bucket = "tf-bucket-for-zos"
  depends_on = [ctyun_zos_subscribe.test]
}