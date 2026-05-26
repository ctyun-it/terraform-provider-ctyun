# 触发云助手命令执行（内联新建命令）
resource "ctyun_cloud_assistant_invocation" "%[1]s" {
  instance_ids    = "%[2]s"
  command_id      = "%[3]s"
  working_directory = "%[4]s"
  timeout         = %[5]d
}

# 查询云助手命令执行结果
data "ctyun_cloud_assistant_invocation_results" "%[6]s" {
  invoked_id = ctyun_cloud_assistant_invocation.%[1]s.id
  page_no    = %[7]d
  page_size  = %[8]d
}
