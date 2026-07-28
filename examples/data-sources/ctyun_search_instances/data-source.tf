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

# 查询 OpenSearch 实例列表（默认配置）
data "ctyun_search_instances" "all_instances" {
  region_id  = "200000002401"
  # page_no   = 1     # 默认为 1
  # page_size = 10    # 默认为 10
}

# 分页查询 OpenSearch 实例
data "ctyun_search_instances" "instances_page2" {
  region_id  = "200000002401"
  page_no    = 2
  page_size  = 20
}

# 输出查询结果
output "all_instances_info" {
  value = {
    id        = data.ctyun_search_instances.all_instances.id
    region_id = data.ctyun_search_instances.all_instances.region_id
    total     = data.ctyun_search_instances.all_instances.total
    
    instances = [
      for instance in data.ctyun_search_instances.all_instances.instances : {
        id              = instance.id
        cluster_name    = instance.cluster_name
        status          = instance.status
        create_time     = instance.create_time
        region_id       = instance.region_id
        available_zone_id = instance.available_zone_id
      }
    ]
  }
}

output "page2_instance_names" {
  value = [
    for instance in data.ctyun_search_instances.instances_page2.instances :
    instance.cluster_name
  ]
}
