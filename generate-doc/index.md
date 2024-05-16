# terraform-provider-ctyun使用说明



## 背景

为了保证信息安全，terraform-provider-ctyun暂未上架到Terraform的[官方镜像源仓库](https://registry.terraform.io/)；但后续会进行推送操作，并且持续更新版本维护；因此，现阶段联调使用我司提供离线的terraform-provider-ctyun进行使用。



## 依赖

- terraform最新版本（v1.6.6），[下载地址](https://developer.hashicorp.com/terraform/install)，请按照官方指引安装
- terraform-provider-ctyun插件，linux请使用terraform-provider-ctyun，windows使用terraform-provider-ctyun.exe
- 配置样例若干
- 插件使用说明：
  - linux：请使用linux_amd64或linux_arm64目录下的插件
  - windows：请使用windows_amd64或windows_arm64目录下的插件
  - Mac：请使用darwin_amd64或darwin_arm64目录下的插件
  - 如需其他环境的编译支撑，请联系天翼云方对接人
- 查看用户自己的ak、sk，[点击进入](https://www.ctyun.cn/console/user/setting)，在账号中心，安全设置，用户AccessKey中查看



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



## 开发调试指南

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



## 最佳实践&建议配置

**强烈建议您在系统环境变量中配置下面设置**

- TF_LOG=INFO，terraform的日志输出级别
- TF_LOG_PATH=terraform的日志输出路径
- TF_CLI_ARGS_apply=-parallelism=并发个数，terraform的最大启动实例个数，建议改小
- CTYUN_AK=您的AK，如果此值不在环境变量中配置，则读取provider中的ak，**推荐**
- CTYUN_SK=您的SK，如果此值不在环境变量中配置，则读取provider中的sk，**推荐**
- CTYUN_REGION_ID=对应的区域id，如果此值不在环境变量中配置，则读取provider中的region_id，**推荐**
- CTYUN_AZ_NAME=对应的可用区id，如果此值不在环境变量中配置，则读取provider中的az_name，**推荐**
- CTYUN_ENV=选用环境，如果此值不在环境变量中配置，则读取provider中的env



## 说明
- **目前terraform-provider-ctyun仅支持4.0的资源池**，3.0的资源池正在积极接入中，请选择下面的4.0资源池进行使用，建议您的测试使用**华东1可用区1**进行测试，`region_id=bb9fdb42056f11eda1610242ac110002,az_name=cn-huadong1-jsnj1A-public-ctcloud`

| 区域名称   | region_id                        | 可用区名称 | az_name                           |
| ---------- | -------------------------------- | ---------- | --------------------------------- |
| 太原4      | 200000002689                     | 可用区1    | cn-sx-tyn4-1a-public-ctcloud      |
| 西南2-贵州 | 200000002927                     | 可用区1    | cn-xinan2-gz-1a-public-ctcloud    |
| 郑州5      | 200000002586                     | 可用区1    | cn-ha-cgo5-1a-public-ctcloud      |
| 青岛20     | 200000001703                     | 可用区1    | cn-sd-qd20-sdqd1A-public-ctcloud  |
| 武汉41     | 200000001781                     | 可用区1    | cn-hb-wh41-hbwh1A-public-ctcloud  |
| 西南1      | 200000002368                     | 可用区1    | cn-xinan1-xn1A-public-ctcloud     |
| 华南2      | 200000002530                     | 可用区1    | cn-huanan2-1A-public-ctcloud      |
| 华北2      | 200000001852                     | 可用区2    | cn-huabei2-tj-2a-public-ctcloud   |
| 南宁23     | 200000001704                     | 可用区1    | cn-gx-nn23-gxnn1A-public-ctcloud  |
| 华北2      | 200000001852                     | 可用区1    | cn-huabei2-tj-1a-public-ctcloud   |
| 上海36     | 200000001790                     | 可用区1    | cn-sh36-sh1A-public-ctcloud       |
| 西南1      | 200000002368                     | 可用区2    | cn-xinan1-xn2A-public-ctcloud     |
| 长沙42     | 200000002401                     | 可用区1    | cn-hn-cs42-hncs1A-public-ctcloud  |
| 南昌5      | 200000002527                     | 可用区1    | cn-jx-nc5-jxnc1A-public-ctcloud   |
| 华东1      | bb9fdb42056f11eda1610242ac110002 | 可用区3    | cn-huadong1-jsnj3A-public-ctcloud |
| 华东1      | bb9fdb42056f11eda1610242ac110002 | 可用区1    | cn-huadong1-jsnj1A-public-ctcloud |
| 华东1      | bb9fdb42056f11eda1610242ac110002 | 可用区2    | cn-huadong1-jsnj2A-public-ctcloud |

- 目前**暂未支持资源的ImportState**，预计未来版本会逐步接入



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
  alias     = "huabei"
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
resource "ctyun_security_group_rule" "security_group_rule_ingress_in_common" {
  security_group_id = "sg-5we39vmesy"
  direction         = "ingress"
  action            = "accept"
  priority          = 60
  protocol          = "any"
  ether_type        = "IPv4"
  dest_cidr_ip      = "0.0.0.0/0"
  description       = "80-90端口"
}

# 通过指定provider方式，在华北2创建资源
resource "ctyun_security_group_rule" "security_group_rule_ingress_in_huabei" {
  provider          = ctyun.huabei
  security_group_id = "sg-8ks24nnukg"
  direction         = "ingress"
  action            = "accept"
  priority          = 60
  protocol          = "any"
  ether_type        = "IPv4"
  dest_cidr_ip      = "0.0.0.0/0"
  description       = "80-90端口"
}
```



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `ak` (String) 身份信息ak
- `az_name` (String) 可用区id，如果是3.0资源池，则此值无需填写；如果是4.0资源池，则填写选用的az_name
- `console_url` (String) 请求分发地址，仅供测试使用，需配合inspect_url_keywords一起使用
- `env` (String) 环境类型env，可选值为：dev：开发环境、test：测试环境、prod：生产环境，默认为生产环境prod
- `inspect_url_keywords` (Set of String) 请求拦截的地址，仅供测试使用，如果填入*则表示拦截所有请求，需配合console_url一起使用
- `project_id` (String) 企业项目id，不填则使用用户默认的企业项目
- `region_id` (String) 资源区域id
- `sk` (String) 身份信息sk