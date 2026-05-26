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



# 查询云助手命令执行结果
data "ctyun_cloud_assistant_invocation_results" "test" {
  page_no    = 1
  page_size  = 10
}
