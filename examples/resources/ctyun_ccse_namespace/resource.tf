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

resource "ctyun_ccse_namespace" "test" {
  cluster_id  = "91434876ed1c42d6a3d77dcc4f414bea"
  values_yaml = <<EOF
apiVersion: v1
kind: Namespace
metadata:
  name: test
EOF
}
