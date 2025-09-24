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

resource "ctyun_ccse_node_association" "ykyedkmadc" {
  cluster_id = "c7fdfadd092643c8aa11d2a330f27873"
  instance_type = "ecs"
  instance_id = "02d8a872-7d1a-45ab-9bd8-9b158376ba3a"
  mirror_id = "3d2c356a-685a-4e8c-b904-bb0725bfc220"
  visibility_post_host_script = "YWJj"
  visibility_host_script = "MTIz"
  password = var.password
}

variable "password" {
  type      = string
  sensitive = true
}