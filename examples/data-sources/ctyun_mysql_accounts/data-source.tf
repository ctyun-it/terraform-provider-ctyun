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

# instance id 和 name 在查询时可以根据实际情况进行替换
data "ctyun_mysql_accounts" "examples" {
  instance_id = "2619a9e7300348809f4dbb2e7ead6cc0"
  name        = "root"
}
