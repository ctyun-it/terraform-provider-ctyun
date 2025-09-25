# Terraform Provider Ctyun

## 依赖项

- 开发依赖项
  - [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.10.5
  - [Go](https://golang.org/doc/install) >= 1.23



## 本地编译

```
git clone https://github.com/ctyun-it/terraform-provider-ctyun.git
cd terraform-provider-ctyun
go build .
```



## Terraform配置项

```
terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
      version = "1.0.1"
    }
  }
}

provider "ctyun" {
  # Configuration options
}
```



## 推荐配置项

- CTYUN_AK = 您的AK，如果此值不在环境变量中配置，则读取provider中的ak
- CTYUN_SK = 您的SK，如果此值不在环境变量中配置，则读取provider中的sk
- CTYUN_REGION_ID = 对应的区域id，如果此值不在环境变量中配置，则读取provider中的region_id
- CTYUN_AZ_NAME = 对应的可用区id，如果此值不在环境变量中配置，则读取provider中的az_name
- CTYUN_ENV = 选用环境，如果此值不在环境变量中配置，则读取provider中的env
- TF_LOG = INFO，terraform的日志输出级别
- TF_LOG_PATH = 路径目录，terraform的日志输出路径



## 文档参考

详见工程中的[website](https://github.com/ctyun-it/terraform-provider-ctyun/tree/main/website)


## 版权声明

Copyright@2024  China Telecom Cloud Technology Co., Ltd. （天翼云科技有限公司）