# 触发云助手命令执行（内联新建命令 + 参数 + 工作目录）
resource "ctyun_cloud_assistant_invocation" "%[1]s" {
  instance_ids       = "%[2]s"
  command_id   = %[3]s
  working_directory = "%[4]s"
  timeout            = %[5]d
  parameter          = %[6]s
}
