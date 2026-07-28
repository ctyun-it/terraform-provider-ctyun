terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

provider "ctyun" {
  env = "prod"
}

resource "ctyun_acl" "example" {
  vpc_id             = "vpc-idexample1"
  name               = "example-acl"
  description        = "Example ACL created for demonstration"
  enabled            = true
  apply_to_public_lb = false
}

resource "ctyun_acl_rule" "example" {
  acl_id                 = ctyun_acl.example.id
  direction              = "ingress"
  priority               = 100
  protocol               = "tcp"
  ip_version             = "ipv4"
  destination_port       = "8080:8085"
  source_port            = "8080:8085"
  source_ip_address      = "192.168.1.0/24"
  destination_ip_address = "192.168.2.0/24"
  action                 = "accept"
  enabled                = true
  description            = "Example ACL rule"
}