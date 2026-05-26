# 执行云助手命令（带自定义参数）
resource "ctyun_cloud_assistant_run_command" "%[1]s" {
  instance_ids       = "%[2]s"
  command_name       = "%[3]s"
  command_type       = "%[4]s"
  command_content    = "%[5]s"
  description        = "%[6]s"
  working_directory  = "%[7]s"
  save_command       = %[8]t
  enabled_parameter  = %[9]t
  timeout            = %[10]d
  parameter          = %[11]s
}
