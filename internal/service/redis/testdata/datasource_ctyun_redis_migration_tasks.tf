data "ctyun_redis_migration_tasks" "%[1]s"{
    %[2]s
    page_num = 1
    page_size = 10
}