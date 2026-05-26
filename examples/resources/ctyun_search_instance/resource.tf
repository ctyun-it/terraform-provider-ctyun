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
