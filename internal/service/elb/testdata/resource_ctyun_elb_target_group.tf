resource "ctyun_elb_target_group" "%[1]s" {
  name      = "%[2]s"
  vpc_id    = "%[3]s"
  algorithm = "%[4]s"
  %[5]s
  %[6]s
  %[7]s
  %[8]s
  %[9]s
  %[10]s
  %[11]s
}
