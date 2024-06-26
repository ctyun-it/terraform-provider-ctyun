# ctyun_security_group_rule (Resource)
**详细说明请见文档：https://www.ctyun.cn/document/10026730/10225510**



## 样例

```terraform
terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

#resource "ctyun_security_group_rule" "security_group_rule_ingress" {
#  security_group_id = "sg-nyq29edx57"
#  direction         = "ingress"
#  action            = "accept"
#  priority          = 50
#  protocol          = "tcp"
#  ether_type        = "IPv4"
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
#  ether_type        = "IPv4"
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `action` (String) 拒绝策略，accept：允许，drop：拒绝
- `direction` (String) 规则方向，egress：出方向，ingress：入方向
- `ether_type` (String) IP类型：ipv4、ipv6
- `protocol` (String) 协议类型: tcp、udp、icmp、any，当此值填写any时，range的值不能设置
- `security_group_id` (String) 安全组id

### Optional

- `description` (String) 描述，长度1-128
- `dest_cidr_ip` (String) 远端地址，为cidr地址格式，如果不填默认为0.0.0.0/0
- `priority` (Number) 优先级：1~100，取值越小优先级越大，默认优先级为50
- `range` (String) 安全组开放的传输层协议相关的源端端口范围，格式如：8000-9000，如果仅开放单一端口则直接填写，如：22，中间不能有空格以及其他特殊字符；如果protocol的值为any，请保证此值留空，如果protocol的值为tcp或udp，此值必填
- `region_id` (String) 资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID

### Read-Only

- `id` (String) id