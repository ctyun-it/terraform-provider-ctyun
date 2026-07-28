---
subcategory: "天翼云OpenSearch服务"
page_title: "CTYUN: ctyun_search_instance"
---

# ctyun_search_instance (Resource)
-> 管理OpenSearch实例



## Example

```terraform
terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 可参考 index.md，在环境变量中配置 ak、sk、资源池 ID、可用区名称
provider "ctyun" {
  env = "prod"
}

# 创建 VPC（如果已有可用的 VPC，可跳过此步）
resource "ctyun_vpc" "example" {
  name        = "tf-vpc-for-opensearch"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  region_id   = "200000002401"
}

# 创建子网
resource "ctyun_subnet" "example" {
  name        = "tf-subnet-for-opensearch"
  cidr        = "192.168.0.0/24"
  description = "terraform测试使用"
  vpc_id      = ctyun_vpc.example.id
  region_id   = "200000002401"
  dns         = ["8.8.8.8", "8.8.4.4"]
}

# 创建安全组
resource "ctyun_security_group" "example" {
  name        = "tf-sg-for-opensearch"
  description = "terraform测试使用"
  vpc_id      = ctyun_vpc.example.id
  region_id   = "200000002401"
}

# 查询可用区
data "ctyun_zones" "test" {
  region_id = "200000002401"
}


variable "password" {
  type      = string
  sensitive = true
}

# 创建 OpenSearch 实例（按需付费）
resource "ctyun_search_instance" "opensearch_ondemand" {
  name               = "tf-opensearch-ondemand"
  region_id          = "200000002401"
  zone_list          = [data.ctyun_zones.test.zones[0]]
  vpc_id             = ctyun_vpc.example.id
  subnet_id          = ctyun_subnet.example.id
  security_group_id  = ctyun_security_group.example.id
  enable_ipv6        = false
  password           = var.password
  cluster_type       = 1  # 1: OpenSearch, 2: Elasticsearch
  os_type            = "CTyun"
  enable_https       = "CLOSE"
  cycle_type         = "on_demand"  # on_demand: 按需付费, month: 包月, year: 包年

  # 节点配置
  node_details = [
    {
      host_num        = 3
      storage_type    = "SSD-genric"
      storage_space   = 40
      flavor_name     = "esearch-4c16g"
      node_group_type = "MASTER"
    },
    {
      host_num        = 3
      storage_type    = "SSD-genric"
      storage_space   = 40
      flavor_name     = "esearch-4c16g"
      node_group_type = "EXCLUSIVE_MASTER"
    },
    {
      host_num        = 3
      storage_type    = "SSD-genric"
      storage_space   = 40
      flavor_name     = "esearch-4c16g"
      node_group_type = "COORDINATE"
    },
    {
      host_num        = 3
      storage_type    = "SSD-genric"
      storage_space   = 40
      flavor_name     = "esearch-4c16g"
      node_group_type = "COLD"
    }
  ]
}

# 创建 OpenSearch 实例（包月付费）
resource "ctyun_search_instance" "opensearch_monthly" {
  name               = "tf-opensearch-monthly"
  region_id          = "200000002401"
  zone_list          = [data.ctyun_zones.test.zones[0]]
  vpc_id             = ctyun_vpc.example.id
  subnet_id          = ctyun_subnet.example.id
  security_group_id  = ctyun_security_group.example.id
  enable_ipv6        = false
  password           = var.password
  cluster_type       = 1  # 1: OpenSearch, 2: Elasticsearch
  os_type            = "CTyun"
  enable_https       = "CLOSE"
  cycle_type         = "month"  # on_demand: 按需付费, month: 包月, year: 包年
  cycle_count        = 1       # cycle_type=month时取1-11, cycle_type=year时取1-5

  # 节点配置
  node_details = [
    {
      host_num        = 3
      storage_type    = "SSD-genric"
      storage_space   = 40
      flavor_name     = "esearch-4c16g"
      node_group_type = "MASTER"
    },
    {
      host_num        = 3
      storage_type    = "SSD-genric"
      storage_space   = 40
      flavor_name     = "esearch-4c16g"
      node_group_type = "EXCLUSIVE_MASTER"
    },
    {
      host_num        = 3
      storage_type    = "SSD-genric"
      storage_space   = 40
      flavor_name     = "esearch-4c16g"
      node_group_type = "COORDINATE"
    },
    {
      host_num        = 3
      storage_type    = "SSD-genric"
      storage_space   = 40
      flavor_name     = "esearch-4c16g"
      node_group_type = "COLD"
    }
  ]
}

# 创建 Elasticsearch 实例（包年付费）
resource "ctyun_search_instance" "elasticsearch_yearly" {
  name               = "tf-elasticsearch-yearly"
  region_id          = "200000002401"
  zone_list          = [data.ctyun_zones.test.zones[0]]
  vpc_id             = ctyun_vpc.example.id
  subnet_id          = ctyun_subnet.example.id
  security_group_id  = ctyun_security_group.example.id
  enable_ipv6        = false
  password           = var.password
  cluster_type       = 2  # 1: OpenSearch, 2: Elasticsearch
  os_type            = "CTyun"
  enable_https       = "CLOSE"
  cycle_type         = "year"  # on_demand: 按需付费, month: 包月, year: 包年
  cycle_count        = 1       # cycle_type=month时取1-11, cycle_type=year时取1-5

  # 节点配置
  node_details = [
    {
      host_num        = 3
      storage_type    = "SSD-genric"
      storage_space   = 40
      flavor_name     = "esearch-4c16g"
      node_group_type = "MASTER"
    },
    {
      host_num        = 3
      storage_type    = "SSD-genric"
      storage_space   = 40
      flavor_name     = "esearch-4c16g"
      node_group_type = "EXCLUSIVE_MASTER"
    },
    {
      host_num        = 3
      storage_type    = "SSD-genric"
      storage_space   = 40
      flavor_name     = "esearch-4c16g"
      node_group_type = "COORDINATE"
    },
    {
      host_num        = 3
      storage_type    = "SSD-genric"
      storage_space   = 40
      flavor_name     = "esearch-4c16g"
      node_group_type = "COLD"
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster_type` (Number) 集群类型：1：OpenSearch，2：Elasticsearch
- `cycle_type` (String) 付费周期类型，month：包月，year：包年，on_demand：按需付费
- `enable_ipv6` (Boolean) 开启IPv6：开启:true 关闭:false
- `name` (String) 实例名称，由大小写字母、数字、下划线(_)或连字符(-)组成，且不以下划线(_)或连字符(-)开头，长度是1-32位
- `node_details` (Attributes Set) 节点组详情列表，支持配置数据节点(MASTER)、专属主节点(EXCLUSIVE_MASTER)、协调节点(COORDINATE)和冷数据节点(COLD)四类节点 (see [below for nested schema](#nestedatt--node_details))
- `os_type` (String) 操作系统类型，ctyun操作系统：CTyun、麒麟操作系统：Kylin
- `security_group_id` (String) 安全组ID
- `subnet_id` (String) 子网ID
- `vpc_id` (String) VPC ID
- `zone_list` (Set of String) 实例所在可用区信息，只能传一个或三个可用区，可通过ctyun_zones查看

### Optional

- `cycle_count` (Number) 订购时长，cycle_type=month时取值为1-11，cycle_type=year时取值为1-5，cycle_type=on_demand时无需填写
- `enable_https` (String) 不开启 https: CLOSE,开启：OPEN；默认 CLOSE
- `password` (String, Sensitive) 组件密码，创建时必填。密码应为数字、大写字母、小写字母、特殊符号 (@$!%*#_~?) 的组合，长度在 12－26 位
- `region_id` (String) 资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID

### Read-Only

- `id` (String) 实例ID
- `status` (String) 实例状态
- `version` (String) 版本信息

<a id="nestedatt--node_details"></a>
### Nested Schema for `node_details`

Required:

- `flavor_name` (String) 实例规格名称，每个资源池可用区下可选择的机型，参考订购页面展示的机型信息，如 s3.medium.2、s3.large.4 等
- `host_num` (Number) 节点数量，MASTER 节点最小为 3，最大为 50；EXCLUSIVE_MASTER 节点最大为 3；COORDINATE 节点最大为 32；COLD 节点最大为 50
- `node_group_type` (String) 节点组类型：MASTER（数据节点）/EXCLUSIVE_MASTER（专属master节点）/COORDINATE（专属协调节点）/COLD（冷数据节点）
- `storage_space` (Number) 存储空间(GB)，MASTER 节点可选：40-6144GB；EXCLUSIVE_MASTER 节点固定 40GB；COORDINATE 节点固定 40GB；COLD 节点可选：40-6144GB
- `storage_type` (String) 存储类型：SSD-genric（通用型SSD）、SAS（高IO）、SSD（超高IO）、XSSD-0、XSSD-1
## 导入

使用以下语法支持导入：

```shell
# 导入 OpenSearch 实例资源
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数，不填则取值环境变量值
terraform import ctyun_search_instance.[导入配置名称] [id][,<region_id>]
# 示例 1: 只指定实例 ID，region_id 从 provider 或环境变量获取
terraform import ctyun_search_instance.opensearch_basic opensearch-xxxxxxxxxxxxxxxxx
# 示例 2: 同时指定实例 ID 和 region_id
terraform import ctyun_search_instance.opensearch_basic opensearch-xxxxxxxxxxxxxxxxx,200000002401
```