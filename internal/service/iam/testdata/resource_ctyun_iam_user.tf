resource "ctyun_iam_user" "%[1]s" {
  email          = "%[2]s"
  phone          = "%[3]s"
  name           = "%[4]s"
  password       = "%[5]s"
  description    = "%[6]s"
  user_group_ids = [
    "%[7]s"
  ]
}
