resource "ctyun_iam_policy_association_user_group" "%[1]s" {
  user_group_id   = "%[2]s"
  policy_id = "%[3]s"
}