resource "ctyun_postgresql_security_group" "%[1]s" {
	instance_id = "%[2]s"
	security_group_ids = [%[3]s]
}