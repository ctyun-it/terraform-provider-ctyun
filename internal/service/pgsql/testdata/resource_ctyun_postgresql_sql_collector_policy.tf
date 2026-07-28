resource "ctyun_postgresql_sql_collector_policy" "%[1]s" {
	instance_id = "%[2]s"
	sql_collector_status = "%[3]s"
	log_interval= %[4]d
}
