# 配置文件: resource_ctyun_mysql_audit.tf
resource "ctyun_mysql_audit" "%[1]s" {
  instance_id  = "%[2]s"
  audit_switch = %[3]t
}