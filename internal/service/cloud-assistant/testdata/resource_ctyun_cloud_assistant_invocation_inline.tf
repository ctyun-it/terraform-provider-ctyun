# 触发云助手命令执行（内联新建命令）
resource "ctyun_cloud_assistant_invocation" "%[1]s" {
  instance_ids    = "%[2]s"
  command_name    = "%[3]s"
  command_type    = "%[4]s"
  command_content = <<-EOF %[5]s EOF
  save_command    = %[6]t
  timeout         = %[7]d
}
