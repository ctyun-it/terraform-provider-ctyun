terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 可参考index.md，在环境变量中配置ak、sk、资源池ID、可用区名称
provider "ctyun" {
  env = "prod"
}


variable "password" {
  type      = string
  sensitive = true
}


resource "ctyun_redis_migration_task" "test" {
  sync_mode     = 1
  conflict_mode = 2
  source_db_info = {
    spu_inst_id  = "425c9173f98b4646a72ce0b986af00b0"
    ip_addr      = "192.168.0.10"
    account_name = "testUser"
    password     = var.password
  }
  target_db_info = {
    spu_inst_id  = "d052e0d8e0204466b2e2673c34247b54"
    ip_addr      = "192.168.0.11"
    account_name = "testUser2"
    password     = var.password
  }

}

//结束运行中的任务
#  resource "ctyun_redis_migration_task" "test" {
#    sync_mode = 1
#    conflict_mode = 2
#    source_db_info = {
#        spu_inst_id = "425c9173f98b4646a72ce0b986af00b0"
#        ip_addr = "192.168.0.10"
#        account_name = "testUser"
#        password = var.password
#    }
#    target_db_info = {
#        spu_inst_id = "d052e0d8e0204466b2e2673c34247b54"
#        ip_addr = "192.168.0.11"
#        account_name = "testUser2"
#        password = var.password
#    }
#     operate_type = 2 //结束运行中的任务
#   }