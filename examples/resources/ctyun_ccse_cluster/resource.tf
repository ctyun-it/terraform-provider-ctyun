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

resource "ctyun_vpc" "vpc_test" {
  name        = "vpc-test-ccse1"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "subnet-test-ccse1"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns         = [
    "100.95.0.1"
  ]
  enable_ipv6 = true
}

data "ctyun_ecs_flavors" "ecs_flavor_test" {
  cpu    = 4
  ram    = 8
  arch   = "x86"
  series = "C"
  type   = "CPU_C7"
}


resource "ctyun_ccse_cluster" "example" {
  base_info = {
    vpc_id     = ctyun_vpc.vpc_test.id
    subnet_id  = ctyun_subnet.subnet_test.id
    cluster_name = "tf-custer-n213nds9124n3"
    cluster_domain = "www.ccc.com"
    network_plugin = "cubecni"
    start_port = 30000
    end_port   = 65535
    elb_prod_code = "standardI"
    pod_subnet_id_list = [ctyun_subnet.subnet_test.id]
    cycle_type  = "on_demand"
    container_runtime = "containerd"
    timezone    = "Asia/Shanghai"
    cluster_version = "1.29.3"
    deploy_type   = "single"
    kube_proxy    = "iptables"
    cluster_series = "cce.standard"
    enable_api_server_eip = true
    enable_snat= true          
    nat_gateway_spec = "small"
    install_als_cube_event = true
    install_als= true          
    install_ccse_monitor = true
    install_nginx_ingress = true
    nginx_ingress_lb_spec = "standardI"
    nginx_ingress_network = "external"
    ip_vlan = true
    network_policy= true
  }

  master_host = {
    item_def_name =  data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].name

    sys_disk = {
      type = "SSD"
      size = 100
    }

    data_disks = [
      {
        type = "SSD"
        size = 200
      }
    ]

    az_infos = [
      {
        az_name = "cn-huadong1-jsnj1A-public-ctcloud"
        size    = 1
      }
    ]
  }

  slave_host = {
    instance_type = "ecs"
    mirror_id     = "3f80d8c0-8eb5-4afa-a506-13ba68b61872"
    mirror_type   = 1
    item_def_name = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].name

    az_infos = [
      {
        az_name = "cn-huadong1-jsnj2A-public-ctcloud"
        size    = 1
      }
    ]

    sys_disk = {
      type = "SATA"
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

# resource "ctyun_ccse_cluster" "example2" {
#   base_info = {
#     vpc_id     = ctyun_vpc.vpc_test.id
#     subnet_id  = ctyun_subnet.subnet_test.id
#     cluster_name = "auto-sec-grqq33"
#     cluster_domain = "www.ccc.s"
#     network_plugin = "cubecni"
#     start_port = 30000
#     end_port   = 65535
#     elb_prod_code = "standardI"
#     pod_subnet_id_list = [ctyun_subnet.subnet_test.id]
#     cycle_type  = "on_demand"
#     container_runtime = "containerd"
#     timezone    = "Asia/Shanghai"
#     cluster_version = "1.29.3"
#     deploy_type   = "single"
#     kube_proxy    = "ipvs"
#     cluster_series = "cce.managed"
#     series_type = "managedbase"
#   }
#
#
#   slave_host = {
#     instance_type = "ecs"
#     mirror_id     = "3f80d8c0-8eb5-4afa-a506-13ba68b61872"
#     mirror_type   = 1
#     item_def_name = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].name
#
#     az_infos = [
#       {
#         az_name = "cn-huadong1-jsnj2A-public-ctcloud"
#         size    = 1
#       }
#     ]
#
#     sys_disk = {
#       type = "SATA"
#       size = 80
#     }
#
#     data_disks = [
#       {
#         type = "SATA"
#         size = 150
#       }
#     ]
#   }
# }


# 裸金属
# resource "ctyun_ccse_cluster" "example" {
#   base_info = {
#     vpc_id     = ctyun_vpc.vpc_test.id
#     subnet_id  = ctyun_subnet.subnet_test.id
#     cluster_name = "fe-ccse3dsfsdqq3"
#     cluster_domain = "www.ccc.s"
#     network_plugin = "calico"
#     start_port = 30000
#     end_port   = 65535
#     elb_prod_code = "standardI"
#     pod_cidr = "172.26.0.0/16"
#     pod_subnet_id_list = [ctyun_subnet.subnet_test.id]
#     cycle_type  = "on_demand"
#     container_runtime = "containerd"
#     timezone    = "Asia/Shanghai"
#     cluster_version = "1.29.3"
#     deploy_type   = "single"
#     kube_proxy    = "iptables"
#     cluster_series = "cce.standard"
#     enable_api_server_eip = true
#     enable_snat= true
#     nat_gateway_spec = "small"
#     install_als_cube_event = true
#     install_als= true
#     install_ccse_monitor = true
#     install_nginx_ingress = true
#     nginx_ingress_lb_spec = "standardI"
#     nginx_ingress_network = "external"
#     ip_vlan= true
#     network_policy= true
#   }
#
#   master_host = {
#     item_def_name =  data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].name
#     az_infos = [
#       {
#         az_name = "cn-huadong1-jsnj1A-public-ctcloud"
#         size    = 1
#       }
#     ]
#   }
#
#   slave_host = {
#     instance_type = "ebm"
#     mirror_name     = "CTyunOS23.01@cpu_ccse_img_4.0_09"
#     mirror_type   = 1
#     item_def_name = "physical.s5.2xlarge4"
#
#     az_infos = [
#       {
#         az_name = "cn-huadong1-jsnj2A-public-ctcloud"
#         size    = 1
#       }
#     ]
#
#   }
# }