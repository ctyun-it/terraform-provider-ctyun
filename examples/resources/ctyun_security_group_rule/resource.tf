terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 可参考index.md，在环境变量中配置ak、sk、资源池ID、可用区名称
provider "ctyun" {
  env = "prod"
}

#resource "ctyun_security_group_rule" "security_group_rule_ingress" {
#  security_group_id = "sg-nyq29edx57"
#  direction         = "ingress"
#  action            = "accept"
#  priority          = 50
#  protocol          = "tcp"
#  ether_type        = "ipv4"
#  dest_cidr_ip      = "0.0.0.0/0"
#  description       = "80-90端口1"
#  range             = "80-90"
#}
#
#resource "ctyun_security_group_rule" "security_group_rule_egress" {
#  security_group_id = "sg-nyq29edx57"
#  direction         = "egress"
#  action            = "drop"
#  priority          = 50
#  protocol          = "udp"
#  ether_type        = "ipv4"
#  dest_cidr_ip      = "0.0.0.0/0"
#  description       = "3306端口1"
#  range             = "3306"
#}

resource "ctyun_security_group_rule" "security_group_rule_egress_any" {
  security_group_id = "sg-p5gue21vy8"
  direction         = "egress"
  action            = "accept"
  priority          = 50
  protocol          = "udp"
  range             = "9999"
  ether_type        = "ipv4"
  dest_cidr_ip      = "192.168.0.0/16"
  description       = "tcp协议9"
}