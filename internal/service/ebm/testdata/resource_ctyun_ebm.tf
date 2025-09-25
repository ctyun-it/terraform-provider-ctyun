resource "ctyun_ebm" "%[1]s" {
  instance_name = "%[2]s"
  hostname = "%[3]s"
  password = "%[4]s"
  status = "%[5]s"
  cycle_type = "month"
  cycle_count = 1
  device_type = "%[6]s"
  image_uuid = "%[7]s"
  security_group_ids = ["%[8]s"]
  vpc_id = "%[9]s"
  system_volume_raid_uuid = "%[10]s"
  data_volume_raid_uuid = "%[11]s"
  subnet_id = "%[12]s"
  az_name = "%[13]s"
}
