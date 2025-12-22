data "ctyun_vpcs" "vpc_test" {
  page_size = 50
}

locals {
  vpcs = [for vpc in data.ctyun_vpcs.vpc_test.vpcs : vpc if vpc.name == "tf-vpc-for-ccse"]
  real_vpc_id = local.vpcs[0].vpc_id
}

data "ctyun_subnets" "subnet_test" {
  vpc_id = local.real_vpc_id
}

locals {
  subnets = [for subnet in data.ctyun_subnets.subnet_test.subnets : subnet if subnet.name == "tf-subnet-for-ccse"]
  real_subnet_id = local.subnets[0].subnet_id

  subnets2 = [for subnet in data.ctyun_subnets.subnet_test.subnets : subnet if subnet.name == "tf-subnet-for-ccse2"]
  real_subnet_id2 = local.subnets2[0].subnet_id
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id      = local.real_vpc_id
  name        = "tf-sg-for-ccse"
  description = "terraform测试使用"
}

data "ctyun_ecs_flavors" "ecs_flavor_test" {
  cpu    = 4
  ram    = 8
  arch   = "x86"
  series = "C"
  type   = "CPU_C7"
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
    mirror_id     = "3f80d8c0-8eb5-4afa-a506-13ba68b61872"
    mirror_type   = 1
    item_def_name = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].name

    az_infos = [
      {
        az_name = "cn-huadong1-jsnj1A-public-ctcloud"
        size    = 1
      }
    ]

    sys_disk = {
      type = "SAS"
      size = 80
    }

    data_disks = [
      {
        type = "SATA"
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
  flavor_id           = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].id
  image_id            = data.ctyun_images.image_test.images[0].id
  security_group_ids  = [ctyun_ccse_cluster.test.base_info.security_group_id]
  system_disk_type    = "sata"
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
# data "ctyun_zones" "test" {
#
# }
#
locals {
  device_type1 = "physical.s5.2xlarge4"      // az2、有本地盘、弹性、不支持云硬盘
  # az2 = data.ctyun_zones.test.zones[1]
}
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
