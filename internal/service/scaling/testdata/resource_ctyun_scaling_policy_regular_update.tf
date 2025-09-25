# 定时策略
resource "ctyun_scaling_policy" "%[1]s" {
  group_id       = %[2]d
  name           = "%[3]s"
  policy_type    = "%[4]s"
  operate_unit   = "%[5]s"
  operate_count  = %[6]d
  action         = "%[7]s"
  execution_time = "%[8]s"
  status         = "%[9]s"
}
