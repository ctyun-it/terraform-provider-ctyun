resource "ctyun_ecs" "%[1]s" {
  instance_name      = "%[2]s"
  display_name       = "%[3]s"
  flavor_name          = "%[4]s"
  image_id           = "%[5]s"
  system_disk_type   = "SATA"
  system_disk_size   = %[6]d
  vpc_id             = "%[7]s"
  subnet_id          = "%[8]s"
  key_pair_name      = "%[9]s"
  security_product  = "%[10]s"
  cycle_type = "%[11]s"
  status = "%[12]s"
}
