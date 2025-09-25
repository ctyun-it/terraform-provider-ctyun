# 周期策略
resource "ctyun_scaling_policy" "%[1]s" {
  group_id       = %[2]d
  name           = "%[3]s"
  policy_type    = "%[4]s"
  operate_unit   = "%[5]s"
  operate_count  = %[6]d
  action         = "%[7]s"
  cycle          = "%[8]s"
  day            = %[9]s
  effective_from = "%[10]s"
  effective_till = "%[11]s"
  execution_time = "%[12]s"
}
