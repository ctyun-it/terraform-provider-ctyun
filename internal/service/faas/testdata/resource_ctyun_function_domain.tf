resource "ctyun_function_domain" "%[1]s" {
  domain_name = "%[2]s"
  protocol    = "HTTP"
  description = "%[3]s"
  cname_check = false
}
