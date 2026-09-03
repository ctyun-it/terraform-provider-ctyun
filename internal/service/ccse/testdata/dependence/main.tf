resource "ctyun_vpc" "vpc_test" {
  name        = "tfvpc-ccse-${local.random_string}"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
  lifecycle {
    ignore_changes = [name]
  }
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tfsubnet-ccse1-${local.random_string}"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns         = [
    "8.8.8.8",
    "8.8.4.4"
  ]
  lifecycle {
    ignore_changes = [name]
  }
}

resource "ctyun_subnet" "subnet_test2" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tfsubnet-ccse2-${local.random_string}"
  cidr        = "192.168.2.0/24"
  description = "terraform测试使用"
  dns         = [
    "8.8.8.8",
    "8.8.4.4"
  ]
  lifecycle {
    ignore_changes = [name]
  }
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tfsg-ccse-${local.random_string}"
  description = "terraform测试使用"
  lifecycle {
    ignore_changes = [name]
  }
}

locals {
  real_vpc_id = ctyun_vpc.vpc_test.id
  real_subnet_id = ctyun_subnet.subnet_test.id
  real_subnet_id2 = ctyun_subnet.subnet_test2.id
}

data "ctyun_ecs_flavors" "ecs_flavor_test" {
  cpu  = 4
  ram  = 8
  arch = "x86"
}

locals {
  available_flavor = [for f in data.ctyun_ecs_flavors.ecs_flavor_test.flavors : f if f.available == true][0]
}

data "ctyun_ccse_images" "ccse_images" {
  instance_type = "ecs"
  flavor_name = local.available_flavor.name
}

data "ctyun_zones" "test" {

}

data "ctyun_ebm_device_types" "test" {
  for_each = { for az in data.ctyun_zones.test.zones: az => az }
  az_name = each.value
}

locals {
  # 筛选全部有余量的规格
  all_device_types = flatten([
    for device_type_inst in values(data.ctyun_ebm_device_types.test) : [
      for dt in device_type_inst.device_types : dt
      # if dt.available == true # 核心过滤条件
    ]
  ])

  # 不支持云盘的弹性裸金属
  cloud_boot_false_list  = [for dt in local.all_device_types : dt if dt.cloud_boot == false]
  first_cloud_boot_false = length(local.cloud_boot_false_list) > 0 ? local.cloud_boot_false_list[0] : null
}

locals {
  cluster_name = "tf-ccse-cluster-${local.random_string}"
}

resource "ctyun_ccse_cluster" "test" {
  lifecycle {
    ignore_changes = [base_info.cluster_name]
  }
  base_info = {
    vpc_id     = local.real_vpc_id
    subnet_id  = local.real_subnet_id
    cluster_name = local.cluster_name
    cluster_domain = "www.ctyun.com"
    network_plugin = "cubecni"
    start_port = 30000
    end_port   = 32767
    elb_prod_code = "standardI"
    pod_subnet_id_list = [local.real_subnet_id]
    cycle_type  = "on_demand"
    container_runtime = "containerd"
    timezone    = "Asia/Shanghai"
    cluster_version = "1.29.3"
    deploy_type   = "single"
    kube_proxy    = "ipvs"
    cluster_series = "cce.managed"
    series_type = "managedpro"
    node_scale = 50
  }


  slave_host = {
    instance_type = "ecs"
    mirror_id     = data.ctyun_ccse_images.ccse_images.images[0].id
    mirror_type   = 1
    item_def_name =  local.available_flavor.name

    az_infos = [
      {
        az_name = data.ctyun_zones.test.zones[0]
        size    = 1
      }
    ]

    sys_disk = {
      type = "XSSD-1"
      size = 80
    }

    data_disks = [
      {
        type = "XSSD-1"
        size = 150
      }
    ]
  }
}

locals {
  chart_name = "node-problem-detector"
  chart_name2 = "cube-cluster-autoscaler"
}

data "ctyun_ccse_plugin_market" "autoscaler" {
  chart_name = local.chart_name2
  chart_version = "1.1.2"
  values_type = "YAML"
  depends_on = [
    ctyun_ccse_cluster.test
  ]
}

resource "ctyun_ccse_plugin" "example1" {
  cluster_id = ctyun_ccse_cluster.test.id
  chart_name = local.chart_name2
  chart_version = "1.1.2"
  values_yaml = data.ctyun_ccse_plugin_market.autoscaler.values
}

data "ctyun_ccse_plugin_market" "test" {
  chart_name = local.chart_name
  depends_on = [
    ctyun_ccse_cluster.test
  ]
}

locals {
  chart_version1 =try(data.ctyun_ccse_plugin_market.test.versions[2].chart_version, "")
  chart_version2 =try(data.ctyun_ccse_plugin_market.test.versions[1].chart_version, "")
}

data "ctyun_ccse_plugin_market" "test1" {
  chart_name = local.chart_name
  chart_version = local.chart_version1
  values_type = "YAML"
  depends_on = [
    ctyun_ccse_cluster.test
  ]
}

data "ctyun_ccse_plugin_market" "test2" {
  chart_name = local.chart_name
  chart_version = local.chart_version2
  values_type = "JSON"
  depends_on = [
    ctyun_ccse_cluster.test
  ]
}

locals {
  # 生成当前时间戳的哈希值
  hash = sha256(timestamp())

  # 从哈希结果中截取字符（转为小写并移除特殊字符）
  random_string = substr(
    replace(
      lower(local.hash),
      "/[^a-z0-9]/",
      ""  # 移除所有非字母数字的字符
    ),
    0, 10  # 截取前16个字符
  )
}

#### 云主机
data "ctyun_images" "image_test" {
  name       = "CentOS Linux 8.4"
  visibility = "public"
  page_no    = 1
  page_size  = 10
}

resource "ctyun_ecs" "ecs_test" {
  instance_name       = "tf-ecs-for-ccse1"
  display_name        = "tf-ecs-for-ccse2"
  flavor_id           =  local.available_flavor.id
  image_id            = data.ctyun_images.image_test.images[0].id
  security_group_ids  = [ctyun_ccse_cluster.test.base_info.security_group_id]
  system_disk_type    = "SAS"
  system_disk_size    = 40
  vpc_id              =  local.real_vpc_id
  password            = var.password
  cycle_type          = "on_demand"
  subnet_id           = local.real_subnet_id
  is_destroy_instance = false
}


variable "password" {
  type      = string
  sensitive = true
}

#### 物理机
#
# data "ctyun_ebm_device_raids" "system_raid" {
#   az_name = local.az2
#   device_type = local.device_type1
#   volume_type = "system"
# }
#
# data "ctyun_ebm_device_raids" "data_raid" {
#   az_name = local.az2
#   device_type = local.device_type1
#   volume_type = "data"
# }
#
# data "ctyun_ebm_device_images" "test" {
#   az_name = local.az2
#   device_type = local.device_type1
#   os_type = "linux"
#   image_type = "standard"
# }
#
# locals {
#   system_raids = data.ctyun_ebm_device_raids.system_raid.raids
#   system_raid_id = length(local.system_raids) > 0 ? local.system_raids[0].uuid : null
#
#   data_raids = data.ctyun_ebm_device_raids.data_raid.raids
#   data_raid_id = length(local.data_raids) > 0 ? local.data_raids[0].uuid : null
# }
#
# data "ctyun_ebm_device_images" "dependence" {
#   az_name = local.az2
#   device_type = local.device_type1
#   os_type = "linux"
#   image_type = "standard"
# }
#
# resource "ctyun_ebm" "ebm_test" {
#   az_name = local.az2
#   instance_name = "tf-ebm-for-ccsedisplay"
#   hostname = "tf-ebm-for-ccse"
#   password = var.password
#   cycle_type = "on_demand"
#   device_type = local.device_type1
#   image_uuid = data.ctyun_ebm_device_images.dependence.images[0].image_uuid
#   security_group_ids = [ctyun_ccse_cluster.test.base_info.security_group_id]
#   system_volume_raid_uuid = local.system_raid_id
#   data_volume_raid_uuid = local.data_raid_id
#   vpc_id = local.real_vpc_id
#   subnet_id = local.real_subnet_id
# }
