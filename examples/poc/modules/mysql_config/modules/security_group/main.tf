terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id = var.vpc_id
  name   = "mysql-security-group-${var.security_group_name}"
}