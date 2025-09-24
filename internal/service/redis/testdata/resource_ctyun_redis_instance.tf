resource "ctyun_redis_instance" "%[1]s" {
  instance_name = "%[2]s"
  version = "%[3]s"
  edition = "%[4]s"
  %[8]s
  %[9]s
  password = "%[10]s"
  engine_version = "%[11]s"
  maintenance_time = "%[12]s"
  protection_status = %[13]s
  shard_mem_size = 8
  vpc_id = "%[5]s"
  subnet_id = "%[6]s"
  security_group_id = "%[7]s"
  cycle_type = "month"
  cycle_count = 1
  auto_renew = true
  auto_renew_cycle_count = 12
}
