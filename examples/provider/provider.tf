terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 完整配置
provider "ctyun" {
  ak                   = "您的ak"                                    # 如果此值不填，则默认读取环境变量中的CTYUN_AK
  sk                   = "您的sk"                                    # 如果此值不填，则默认读取环境变量中的CTYUN_SK
  region_id            = "bb9fdb42056f11eda1610242ac110002"         # 如果此值不填，则默认读取环境变量中的CTYUN_REGION_ID
  az_name              = "cn-huadong1-jsnj1A-public-ctcloud"        # 如果此值不填，则默认读取环境变量中的CTYUN_AZ_NAME
  env                  = "prod"                                     # 如果此值不填，则默认读取环境变量中的CTYUN_ENV
  project_id           = "您的project_id"                            # 如果此值不填，则默认读取环境变量中的CTYUN_PROJECT_ID
  console_url          = "目标consoleUrl"                            # 如果此值不填，则默认读取环境变量中的CTYUN_CONSOLE_URL，仅在非生产环境使用
  inspect_url_keywords = [
    # 如果此值不填，则默认读取环境变量中的CTYUN_INSPECT_URL_KEYWORDS，仅在非生产环境使用
    "拦截的url地址1",
    "拦截的url地址2",
  ]
}

# 下面例子为多provider配置，可以用于不同资源池的配置
# 选用华北2、可用区2为可选资源池
provider "ctyun" {
  alias     = "huabei"                          # 别名
  region_id = "200000001852"
  az_name   = "cn-huabei2-tj-2a-public-ctcloud"
}

# 使用测试环境我的资源池
provider "ctyun" {
  alias     = "test"
  ak        = "您的ak"                              # 如果此值不填，则默认读取环境变量中的CTYUN_AK
  sk        = "您的sk"                              # 如果此值不填，则默认读取环境变量中的CTYUN_SK
  region_id = "81f7728662dd11ec810800155d307d5b"   # 如果此值不填，则默认读取环境变量中的CTYUN_REGION_ID
  az_name   = "az2"                                # 如果此值不填，则默认读取环境变量中的CTYUN_AZ_NAME
  env       = "test"                               # 如果此值不填，则默认读取环境变量中的CTYUN_ENV
}

# 不指定provider选用默认的provider
resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

# 通过指定provider方式，在华北2创建资源
resource "ctyun_vpc" "vpc_test" {
  provider    = ctyun.huabei
  name        = "tf-vpc"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}