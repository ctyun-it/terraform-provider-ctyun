
# 查询触发器列表
data "ctyun_function_triggers" "%[1]s" {
  function_name = "%[2]s"
  page_index    = 1
  page_size     =10
}
