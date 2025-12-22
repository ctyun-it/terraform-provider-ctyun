data "ctyun_iam_user_groups" "%[1]s" {
  name      = "%[2]s"
  page_size = 1000
  page_no   = 1
}