terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id = var.vpc_id
  name   = "mysql-security-group-${var.security_group_name}"
}