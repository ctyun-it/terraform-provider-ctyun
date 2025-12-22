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

resource "ctyun_net_tags" "example" {
  resource_type = "vpc"
  resource_id   = "vpc-example-id"

  tags = [
    {
      key   = "environment"
      value = "production"
    },
    {
      key   = "department"
      value = "devops"
    }
  ]
}