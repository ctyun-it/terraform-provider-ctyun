# 目标追踪策略
resource "ctyun_scaling_policy" "%[1]s" {
  group_id                          = %[2]d
  name                              = "%[3]s"
  policy_type                       = "%[4]s"
  cooldown                          = %[5]d
  target_metric_name                = "%[6]s"
  target_value                      = %[7]d
  target_scale_out_evaluation_count = %[8]d
  target_scale_in_evaluation_count  = %[9]d
  target_operate_range              = %[10]d
  target_disable_scale_in           = %[11]t
}
