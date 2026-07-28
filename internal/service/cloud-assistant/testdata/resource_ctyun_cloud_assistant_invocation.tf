# 触发云助手命令执行
resource "ctyun_cloud_assistant_invocation" "%[1]s" {
  instance_ids    = "%[2]s"
  command_id      = "%[3]s"
  timeout         = %[4]d
}
