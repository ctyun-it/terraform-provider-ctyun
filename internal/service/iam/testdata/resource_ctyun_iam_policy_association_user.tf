resource "ctyun_iam_policy_association_user" "%[1]s" {
  user_id   = "%[2]s"
  policy_id = "%[3]s"
}