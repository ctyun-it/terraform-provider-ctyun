resource "ctyun_redis_migration_task" "%[1]s" {
  sync_mode = 1
  conflict_mode = 2
  source_db_info = {
      spu_inst_id = "%[2]s"
      ip_addr = "%[3]s"
      account_name = "%[4]s"
      password = "%[5]s"
  }
  target_db_info = {
      spu_inst_id = "%[6]s"
      ip_addr = "%[7]s"
      account_name = "%[8]s"
      password = "%[9]s"
  }
   %[10]s
  }


