resource "ctyun_ebs_snapshot_policy" "%[1]s" {
    name           = "%[2]s"
    repeat_weekdays            = "0,1,2"
    repeat_times            = "0,1,2"
    retention_time        = 2
    enabled  = %[3]t
}
