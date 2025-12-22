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

resource "ctyun_acl" "example" {
  vpc_id             = "vpc-exampleid1"
  name               = "example-acl"
  description        = "Example ACL created for demonstration"
  enabled            = true
  apply_to_public_lb = false
}

data "ctyun_acls" "example" {
  id        = ctyun_acl.example.id
  page_no   = 1
  page_size = 20
}

output "ctyun_acls_example" {
  value = data.ctyun_acls.example
}