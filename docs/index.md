# terraform-provider-ctyun使用说明

## 依赖

- terraform v1.10.5，[下载地址](https://developer.hashicorp.com/terraform/install)，请按照官方指引安装
- 天翼云账号AK、SK，[点击进入](https://www.ctyun.cn/console/user/setting)，在账号中心，安全设置，用户AccessKey中查看

## 建议配置

-> 强烈建议您在系统环境变量中进行如下设置

- TF_LOG=INFO，terraform的日志输出级别
- TF_LOG_PATH=terraform的日志输出路径
- TF_CLI_ARGS_apply=-parallelism=并发个数，terraform的最大启动实例个数，不配置默认为10
- CTYUN_AK=您的AK，如果此值不在环境变量中配置，则读取provider中的ak
- CTYUN_SK=您的SK，如果此值不在环境变量中配置，则读取provider中的sk
- CTYUN_REGION_ID=对应的区域id，如果此值不在环境变量中配置，则读取provider中的region_id
- CTYUN_AZ_NAME=对应的可用区id，如果此值不在环境变量中配置，则读取provider中的az_name
- CTYUN_ENV=选用环境，如果此值不在环境变量中配置，则读取provider中的env
- 优先级：`resource/datasource中的配置`高于`provider配置`高于`环境变量配置`。


## 本地安装指南

### Windows

- 以Administrator安装为例，在`C:\Users\Administrator\AppData\Roaming`目录中新建terraform.rc

- 在文件中写入

```
provider_installation {
  filesystem_mirror {
    path = "C:/Users/Administrator/AppData/Roaming/provider-cache"
  }
}
```

- 目录结构准备：`C:\Users\<用户>\AppData\Roaming\provider-cache\registry.terraform.io\ctyun-it\ctyun\<version>\<arch>`

- 例如创建目录：`C:\Users\Administrator\AppData\Roaming\provider-cache\registry.terraform.io\ctyun-it\ctyun\1.2.0\windows_amd64`

- 将可执行文件terraform-provider-ctyun粘贴到上述目录。

### Linux\MacOS

- vi ~/.terraformrc

- 在文件中写入下列内容：

```
provider_installation {
  filesystem_mirror {
    path = "/opt/.terraform.d/provider-cache"
  }
}
```

- 创建目录：

```
# Linux 示例
mkdir -p /opt/.terraform.d/provider-cache/registry.terraform.io/ctyun-it/ctyun/1.2.0/linux_amd64

# MacOS 示例
mkdir -p /opt/.terraform.d/provider-cache/registry.terraform.io/ctyun-it/ctyun/1.2.0/darwin_amd64
```

- 将可执行文件复制到目录中

```
# Linux 示例
cp terraform-provider-ctyun /opt/.terraform.d/provider-cache/registry.terraform.io/ctyun-it/ctyun/1.2.0/linux_amd64/

# MacOS 示例
cp terraform-provider-ctyun /opt/.terraform.d/provider-cache/registry.terraform.io/ctyun-it/ctyun/1.2.0/darwin_amd64/
```

- 添加可执行权限：`chmod +x /opt/.terraform.d/provider-cache/registry.terraform.io/ctyun-it/ctyun/1.2.0/darwin_amd64/terraform-provider-ctyun`


## 安装验证

- 新建目录，建立main.tf文件，并且输入下面内容

```
terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

provider "ctyun" {
  ak        = "您的ak"
  sk        = "您的sk"
  region_id = "bb9fdb42056f11eda1610242ac110002"
  az_name   = "cn-huadong1-jsnj1A-public-ctcloud"
  env       = "prod"
}
```

- 执行`terraform init`命令查看是否成功

## Terraform 基本命令

- `terraform plan`：预览配置变更将对基础设施造成的影响（创建、更新、删除）。

- `terraform apply`：执行变更。

- `terraform destroy`：删除所有通过 Terraform 管理的资源。

- `terraform refresh`：同步远程数据到本地state文件。


## 最佳实践

- 在并发创建相同类型的资源时，名称需要加以区分，建议使用如下写法：

```
resource "ctyun_ecs" "ecs_test" {
  count = 130
  instance_name       = "ds-ecs-${count.index + 1}"
  ...
}
```

- PaaS产品实例在terraform destroy时可以删除，但相关联的底层资源不能马上释放，所以删除子网和安全组时会报错。涉及CCSE、Redis、Kafka、RabbitMq、Mysql、PostgreSql、MogoDB。预计完善时间8月底。
- 如果您想要将state文件保存到对象存储，可参考https://developer.hashicorp.com/terraform/language/v1.11.x/backend/s3，示例如下，endpoints中的s3是控制台页面上桶的终端节点：
- 需要配置对象存储的AK和SK到环境变量：AWS_ACCESS_KEY_ID 和 AWS_SECRET_ACCESS_KEY

```
terraform {
  backend "s3" {
    bucket         = "bucket-xxs"
    key            = "bc6a-ce8f8fb792db"
    region         = "jiangsu-10"
    skip_region_validation      = true
    skip_metadata_api_check     = true
    skip_credentials_validation = true
    skip_requesting_account_id  = true
    skip_s3_checksum            = true
    use_path_style = true
    endpoints = {
      s3 = "https://jiangsu-10.zos.ctyun.cn:443"
    }
  }

  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
      version = "1.2.0"
    }
  }
}
```

## 说明
- **目前terraform-provider-ctyun仅支持多可用区资源池**，建议您的测试使用**华东1可用区1**进行测试，`region_id=bb9fdb42056f11eda1610242ac110002,az_name=cn-huadong1-jsnj1A-public-ctcloud`

- 部分多可用区资源池的region_id及az_name如下，全量资源池可使用ctyun_regions查询

| 区域名称   | region_id                        | 可用区名称 | az_name                           |
| ---------- | -------------------------------- | ---------- | --------------------------------- |
| 呼和浩特3   | 200000003573                     | 可用区1    | cn-nm-het3-1a-public-ctcloud      |
| 乌鲁木齐7   | 200000004098                     | 可用区1    | cn-xj-urc7-1a-public-ctcloud      |
| 西南2-贵州  | 200000002927                     | 可用区1    | cn-xinan2-gz-1a-public-ctcloud    |
| 华东1      | bb9fdb42056f11eda1610242ac110002 | 可用区1    | cn-huadong1-jsnj1A-public-ctcloud |
| 华东1      | bb9fdb42056f11eda1610242ac110002 | 可用区2    | cn-huadong1-jsnj2A-public-ctcloud |
| 华东1      | bb9fdb42056f11eda1610242ac110002 | 可用区3    | cn-huadong1-jsnj3A-public-ctcloud |
| 太原4      | 200000002689                     | 可用区1    | cn-sx-tyn4-1A-public-ctcloud      |
| 郑州5      | 200000002586                     | 可用区1    | cn-ha-cgo5-1a-public-ctcloud      |
| 青岛20     | 200000001703                     | 可用区1    | cn-sd-qd20-sdqd1A-public-ctcloud  |
| 武汉41     | 200000001781                     | 可用区1    | cn-hb-wh41-hbwh1A-public-ctcloud  |
| 西南1      | 200000002368                     | 可用区1    | cn-xinan1-xn1A-public-ctcloud     |
| 西南1      | 200000002368                     | 可用区2    | cn-xinan1-xn2A-public-ctcloud     |
| 华南2      | 200000002530                     | 可用区1    | cn-huanan2-1A-public-ctcloud      |
| 华北2      | 200000001852                     | 可用区1    | cn-huabei2-tj1A-public-ctcloud   |
| 华北2      | 200000001852                     | 可用区2    | cn-huabei2-tj-2a-public-ctcloud   |
| 华北2      | 200000001852                     | 可用区3    | cn-huabei2-tj-3a-public-ctcloud   |
| 南宁23     | 200000001704                     | 可用区1    | cn-gx-nn23-gxnn1A-public-ctcloud  |
| 上海36     | 200000001790                     | 可用区1    | cn-sh36-sh1A-public-ctcloud       |
| 长沙42     | 200000002401                     | 可用区1    | cn-hn-cs42-hncs1A-public-ctcloud  |
| 南昌5      | 200000002527                     | 可用区1    | cn-jx-nc5-jxnc1A-public-ctcloud   |
| 上海32     | 200000001625                     | 可用区1    | cn-sh32-sh1A-public-ctcloud       |
| 杭州7      | 200000003329                     | 可用区1    | cn-zj-hgh7-1a-public-ctcloud      |
| 芜湖4      | 200000003327                     | 可用区1    | cn-ah-whi4-1a-public-ctcloud      |
| 庆阳2      | 200000003664                     | 可用区1    | cn-gs-qyi2-1a-public-ctcloud      |
| 香港2      | 200000002374                     | 可用区1    | hk-hk3-1A-public-ctcloud          |
| 香港2      | 200000002374                     | 可用区2    | cn-hk2-hk2A-public-ctcloud        |
| 澳门1      | 200000002533                     | 可用区1    | mo-mo1-1a-public-ctcloud          |
| 印度尼西亚1 | 200000003424                     | 可用区1    | id-jkt1-1a-public-ctcloud         |
| 新加坡4     | 200000002670                     | 可用区1    | sg-SINP4-1A-public-ctcloud        |
| 新加坡4     | 200000002670                     | 可用区2    | sg-SINP4-2A-public-ctcloud        |
| 菲律宾1     | 200000002769                     | 可用区1    | ph-bdy1-1a-public-ctcloud         |


## 样例

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `ak` (String, Sensitive) 身份信息AK
- `az_name` (String) 可用区英文，填写选用资源池的az_name
- `console_url` (String) 请求分发地址，仅供测试使用，需配合inspect_url_keywords一起使用
- `env` (String) 环境类型env，可选值为：dev：开发环境、test：测试环境、prod：生产环境，默认为生产环境prod
- `inspect_url_keywords` (Set of String) 请求拦截的地址，仅供测试使用，如果填入*则表示拦截所有请求，需配合console_url一起使用
- `project_id` (String) 企业项目ID，不填则使用用户默认的企业项目
- `region_id` (String) 资源池ID
- `sk` (String, Sensitive) 身份信息SK

## 开发调试指南（开发者阅读）

### Windows

- 在`C:\Users\用户名\AppData\Roaming`目录中新建terraform.rc

- 在文件中写入

```
provider_installation {

  dev_overrides {
      "ctyun-it/ctyun"="D:/Go/gobin/"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

- 将编译好的可执行文件放到terraform.rc中指定的目录，这里是`D:/Go/gobin/`

### Linux\MacOS

- vi ~/.terraformrc

- 在文件中写入下列内容，注意要将${pwd}替换成terraform-provider-ctyun文件所在目录

```
provider_installation {

  dev_overrides {
      "ctyun-it/ctyun"="${pwd}"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

- 如出现权限问题，可使用以下命令，如：

```
chmod +x /${pwd}/terraform-provider-ctyun
```