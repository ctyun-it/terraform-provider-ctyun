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

data "ctyun_ccse_namespaces" "all" {
  cluster_id = "91434876ed1c42d6a3d77dcc4f414bea"
}

data "ctyun_ccse_namespaces" "test" {
  cluster_id = "91434876ed1c42d6a3d77dcc4f414bea"
  field      = "metadata.name%3Ddefault"
}
