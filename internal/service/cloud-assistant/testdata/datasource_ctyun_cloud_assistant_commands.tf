# 查询云助手命令列表
data "ctyun_cloud_assistant_commands" "%[1]s" {
  page_no   = %[2]d
  page_size = %[3]d
}
