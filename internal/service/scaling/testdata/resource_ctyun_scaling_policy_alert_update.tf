# 告警策略
resource "ctyun_scaling_policy" "%[1]s" {
  group_id                    = %[2]d
  name                        = "%[3]s"
  policy_type                 = "%[4]s"
  operate_unit                = "%[5]s"
  operate_count               = %[6]d
  action                      = "%[7]s"
  cooldown                    = %[8]d
  trigger_name                = "%[9]s"
  trigger_metric_name         = "%[10]s"
  trigger_statistics          = "%[11]s"
  trigger_comparison_operator = "%[12]s"
  trigger_threshold           = %[13]d
  trigger_period              = "%[14]s"
  trigger_evaluation_count    = %[15]d
  status                      = "%[16]s"
}
