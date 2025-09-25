resource "ctyun_security_group_rule" "%[1]s" {
 security_group_id = "%[2]s"
 direction         = "%[3]s"
 action            = "%[4]s"
 protocol          = "%[5]s"
 ether_type        = "%[6]s"
}