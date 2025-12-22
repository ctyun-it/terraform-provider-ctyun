resource "ctyun_enterprise_project_association_user_group" "%[1]s" {
  enterprise_project_id = %[2]s
  user_group_id         = "%[3]s"
  policy_ids = [
    "%[4]s"
  ]
}