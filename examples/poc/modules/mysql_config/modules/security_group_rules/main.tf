terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
}

# ipv4限制部分
resource "ctyun_security_group_rule" "security_group_rule_ingress_tcp_accept_3306_ipv4" {
  security_group_id = var.security_group_id
  direction         = "ingress"
  action            = "accept"
  priority          = 1
  protocol          = "tcp"
  ether_type        = "IPv4"
  dest_cidr_ip      = "0.0.0.0/0"
  description       = "开放3306端口"
  range             = "3306"
}

resource "ctyun_security_group_rule" "security_group_rule_egress_tcp_accept_3306_ipv4" {
  security_group_id = var.security_group_id
  direction         = "egress"
  action            = "accept"
  priority          = 1
  protocol          = "tcp"
  ether_type        = "IPv4"
  dest_cidr_ip      = "0.0.0.0/0"
  description       = "开放3306端口"
  range             = "3306"
}

# ipv6限制部分
resource "ctyun_security_group_rule" "security_group_rule_ingress_tcp_accept_ipv6" {
  security_group_id = var.security_group_id
  direction         = "ingress"
  action            = "accept"
  priority          = 1
  protocol          = "tcp"
  ether_type        = "IPv6"
  dest_cidr_ip      = "::/0"
  description       = "开放3306端口"
  range             = "3306"
}

resource "ctyun_security_group_rule" "security_group_rule_egress_tcp_accept_ipv6" {
  security_group_id = var.security_group_id
  direction         = "egress"
  action            = "accept"
  priority          = 1
  protocol          = "tcp"
  ether_type        = "IPv6"
  dest_cidr_ip      = "::/0"
  description       = "开放3306端口"
  range             = "3306"
}

