resource "ctyun_iam_policy" "%[1]s"  {
  name        = "%[2]s"
  description = "%[3]s"
  range       = "%[4]s"
  content     = {
    version   = "1.1"
    statement = [
      {
        effect   = "%[5]s"
        action   = ["%[6]s"]
        resource = ["*"]
      }
    ]
  }
}
