# 管理云助手命令（带自定义参数）
resource "ctyun_cloud_assistant_command" "%[1]s" {
  command_name      = "%[2]s"
  command_type      = "%[3]s"
  command_content   = "%[4]s"
  description       = "%[5]s"
  timeout           = %[6]d
  enabled_parameter = %[7]t
  parameter = %[8]s
}
