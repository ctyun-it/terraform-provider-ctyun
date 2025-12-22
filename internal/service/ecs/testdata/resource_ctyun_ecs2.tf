resource "ctyun_ecs" "%[1]s" {
  instance_name      = "%[2]s"
  display_name       = "%[3]s"
  image_id           = "%[4]s"
  system_disk_type   = "sata"
  system_disk_size   = %[5]d
  vpc_id             = "%[6]s"
  subnet_id          = "%[7]s"
  key_pair_name      = "%[8]s"
  %[9]s
}
