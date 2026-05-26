# 执行云助手命令（Shell类型）
resource "ctyun_cloud_assistant_run_command" "%[1]s" {
  instance_ids    = "%[2]s"
  command_name    = "%[3]s"
  command_type    = "%[4]s"
  command_content = "%[5]s"
  description     = "%[6]s"
  save_command    = %[7]t
  timeout         = %[8]d
}
