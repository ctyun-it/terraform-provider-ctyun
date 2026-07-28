# 执行云助手命令 + 查询执行结果
resource "ctyun_cloud_assistant_run_command" "%[1]s" {
  instance_ids    = "%[2]s"
  command_name    = "%[3]s"
  command_type    = "%[4]s"
  command_content = "%[5]s"
  description     = "%[6]s"
  save_command    = %[7]t
  timeout         = %[8]d
}

# 查询云助手命令执行结果
data "ctyun_cloud_assistant_invocation_results" "%[9]s" {
  invoked_id = ctyun_cloud_assistant_run_command.%[1]s.id
  page_no    = %[10]d
  page_size  = %[11]d
}
