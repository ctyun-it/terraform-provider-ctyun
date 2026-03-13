resource "ctyun_ebm" "%[1]s" {
  instance_name = "%[2]s"
  hostname = "%[3]s"
  password = "%[4]s"
  status = "%[5]s"
  cycle_type = "on_demand"
  device_type = "%[6]s"
  image_uuid = "%[7]s"
  vpc_id = "%[8]s"
  system_volume_raid_uuid = "%[9]s"
  data_volume_raid_uuid = "%[10]s"
  subnet_id = "%[11]s"
  az_name = "%[12]s"
}
