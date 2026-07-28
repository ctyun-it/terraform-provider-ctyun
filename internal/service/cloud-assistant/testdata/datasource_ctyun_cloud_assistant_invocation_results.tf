# 查询云助手命令执行结果
data "ctyun_cloud_assistant_invocation_results" "%[1]s" {
  invoked_id = "%[2]s"
  page_no    = %[3]d
  page_size  = %[4]d
}
