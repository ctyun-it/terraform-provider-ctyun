terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# ipv4限制部分
resource "ctyun_security_group_rule" "security_group_rule_ingress_tcp_drop_ipv4" {
  security_group_id = var.security_group_id
  direction         = "ingress"
  action            = "drop"
  priority          = 50
  protocol          = "tcp"
  ether_type        = "IPv4"
  dest_cidr_ip      = "0.0.0.0/0"
  description       = "所有端口"
  range             = "1-65535"
}

resource "ctyun_security_group_rule" "security_group_rule_ingress_udp_drop_ipv4" {
  security_group_id = var.security_group_id
  direction         = "ingress"
  action            = "drop"
  priority          = 50
  protocol          = "udp"
  ether_type        = "IPv4"
  dest_cidr_ip      = "0.0.0.0/0"
  description       = "所有端口"
  range             = "1-65535"
}

resource "ctyun_security_group_rule" "security_group_rule_ingress_icmp_drop_ipv4" {
  security_group_id = var.security_group_id
  direction         = "ingress"
  action            = "drop"
  priority          = 50
  protocol          = "icmp"
  ether_type        = "IPv4"
  dest_cidr_ip      = "0.0.0.0/0"
  description       = "所有端口"
  range             = "1-65535"
}

resource "ctyun_security_group_rule" "security_group_rule_egress_tcp_drop_ipv4" {
  security_group_id = var.security_group_id
  direction         = "egress"
  action            = "drop"
  priority          = 50
  protocol          = "tcp"
  ether_type        = "IPv4"
  dest_cidr_ip      = "0.0.0.0/0"
  description       = "所有端口"
  range             = "1-65535"
}

resource "ctyun_security_group_rule" "security_group_rule_egress_udp_drop_ipv4" {
  security_group_id = var.security_group_id
  direction         = "egress"
  action            = "drop"
  priority          = 50
  protocol          = "udp"
  ether_type        = "IPv4"
  dest_cidr_ip      = "0.0.0.0/0"
  description       = "所有端口"
  range             = "1-65535"
}

resource "ctyun_security_group_rule" "security_group_rule_egress_icmp_drop_ipv4" {
  security_group_id = var.security_group_id
  direction         = "egress"
  action            = "drop"
  priority          = 50
  protocol          = "icmp"
  ether_type        = "IPv4"
  dest_cidr_ip      = "0.0.0.0/0"
  description       = "所有端口"
  range             = "1-65535"
}

# ipv6限制部分
resource "ctyun_security_group_rule" "security_group_rule_ingress_tcp_drop_ipv6" {
  security_group_id = var.security_group_id
  direction         = "ingress"
  action            = "drop"
  priority          = 50
  protocol          = "tcp"
  ether_type        = "IPv6"
  dest_cidr_ip      = "::/0"
  description       = "所有端口"
  range             = "1-65535"
}

resource "ctyun_security_group_rule" "security_group_rule_ingress_udp_drop_ipv6" {
  security_group_id = var.security_group_id
  direction         = "ingress"
  action            = "drop"
  priority          = 50
  protocol          = "udp"
  ether_type        = "IPv6"
  dest_cidr_ip      = "::/0"
  description       = "所有端口"
  range             = "1-65535"
}

resource "ctyun_security_group_rule" "security_group_rule_egress_tcp_drop_ipv6" {
  security_group_id = var.security_group_id
  direction         = "egress"
  action            = "drop"
  priority          = 50
  protocol          = "tcp"
  ether_type        = "IPv6"
  dest_cidr_ip      = "::/0"
  description       = "所有端口"
  range             = "1-65535"
}

resource "ctyun_security_group_rule" "security_group_rule_egress_udp_drop_ipv6" {
  security_group_id = var.security_group_id
  direction         = "egress"
  action            = "drop"
  priority          = 50
  protocol          = "udp"
  ether_type        = "IPv6"
  dest_cidr_ip      = "::/0"
  description       = "所有端口"
  range             = "1-65535"
}


resource "ctyun_security_group_rule" "security_group_rule_egress_tcp_accept_22" {
  security_group_id = var.security_group_id
  direction         = "egress"
  action            = "accept"
  priority          = 1
  protocol          = "tcp"
  ether_type        = "IPv4"
  dest_cidr_ip      = "0.0.0.0/0"
  description       = "22端口"
  range             = "22"
}

resource "ctyun_security_group_rule" "security_group_rule_ingress_tcp_accept_22" {
  security_group_id = var.security_group_id
  direction         = "ingress"
  action            = "accept"
  priority          = 1
  protocol          = "tcp"
  ether_type        = "IPv4"
  dest_cidr_ip      = "0.0.0.0/0"
  description       = "22端口"
  range             = "22"
}