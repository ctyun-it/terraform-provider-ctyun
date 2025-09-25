# ctyun_ccse_cluster (Resource)
**详细说明请见文档：https://www.ctyun.cn/document/10083472/10656137**



## 样例

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `base_info` (Attributes) 集群基础信息 (see [below for nested schema](#nestedatt--base_info))
- `slave_host` (Attributes) slave节点基本信息 (see [below for nested schema](#nestedatt--slave_host))

### Optional

- `master_host` (Attributes) master节点基本信息，专有版必填，托管版时不传 (see [below for nested schema](#nestedatt--master_host))
- `region_id` (String) 资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID

### Read-Only

- `external_kube_config` (String) 外网连接信息
- `id` (String) ID
- `internal_kube_config` (String) 内网连接信息
- `master_order_id` (String) 主订单号
- `name` (String) 名称

<a id="nestedatt--base_info"></a>
### Nested Schema for `base_info`

Required:

- `cluster_domain` (String) 集群本地域名
- `cluster_name` (String) 集群名字
- `cluster_series` (String) 集群系列，支持cce.standard（专有版），cce.managed（托管版），您可查看<a href="https://www.ctyun.cn/document/10083472/10892150">产品定义</a>
- `cluster_version` (String) 集群版本，支持1.31.6，1.29.3，1.27.8，您可查看<a href="https://www.ctyun.cn/document/10083472/10650447">集群版本说明</a>
- `container_runtime` (String) 容器运行时,可选containerd、docker，您可查看<a href="https://www.ctyun.cn/document/10083472/10902208">容器运行时说明</a>
- `cycle_type` (String) 订购周期类型，取值范围：month：按月，year：按年、on_demand：按需。当此值为month或者year时，cycle_count为必填
- `deploy_type` (String) 部署模式，单可用区为single，多可用区为multi
- `elb_prod_code` (String) ApiServer的ELB类型，支持standardI（标准I型），standardII（标准II型），enhancedI（增强I型），enhancedII（增强II型），higherI（高阶I型），您可查看<a href="https://www.ctyun.cn/document/10026756/10032048">ELB类型规格说明</a>
- `end_port` (Number) 节点服务终止端口，可选范围30000-65535，startPort到endPort范围需大于20
- `kube_proxy` (String) kubeProxy类型：iptables或ipvs。您可查看<a href="https://www.ctyun.cn/document/10083472/10915725">iptables与IPVS如何选择</a>
- `network_plugin` (String) 网络插件，可选calico和cubecni，calico需要申请白名单。您可查看<a href="https://www.ctyun.cn/document/10083472/10520760">容器网络插件说明</a>
- `start_port` (Number) 节点服务开始端口，可选范围30000-65535
- `subnet_id` (String) 子网ID
- `timezone` (String) 时区，例如Asia/Shanghai (UTC+08:00)
- `vpc_id` (String) 虚拟私有云ID

Optional:

- `auto_renew` (Boolean) 是否自动续订，默认非自动续订，当cycle_type不等于on_demand时才可填写，按月购买，自动续订周期为1个月；按年购买，自动续订周期为1年。
- `cycle_count` (Number) 订购时长，该参数在cycle_type为month或year时才生效，当cycle_type=month，支持订购1-11个月；当cycle_type=year，支持订购1-3年
- `enable_api_server_eip` (Boolean) 是否开启ApiServerEip，默认false，若开启将自动创建按需计费类型的eip。
- `enable_snat` (Boolean) 是否开启nat网关，默认false，若开启将自动创建按需计费类型的nat网关。
- `install_als` (Boolean) 是否安装日志插件，默认false
- `install_als_cube_event` (Boolean) 是否安装事件采集插件，默认false
- `install_ccse_monitor` (Boolean) 是否安装监控插件，默认false
- `install_nginx_ingress` (Boolean) 是否安装nginx_ingress插件，默认false
- `ip_vlan` (Boolean) 基于IPVLAN做弹性网卡共享，默认false，当指定为true时，主机镜像只有使用CtyunOS系统才能生效
- `nat_gateway_spec` (String) 当enable_snat=true时填写，nat网关规格：small，medium，large，xlarge，可参考<a href="https://www.ctyun.cn/document/10026759/10043996">产品规格说明</a>
- `network_policy` (Boolean) 是否提供基于策略的网络访问控制，默认false
- `nginx_ingress_lb_spec` (String) install_nginx_ingress=true必填，支持规格：standardI（标准I型） ,standardII（标准II型）, enhancedI（增强I型）, enhancedII（增强II型） , higherI（高阶I型），可参考<a href="https://www.ctyun.cn/document/10026756/10032048">规格详情</a>
- `nginx_ingress_network` (String) install_nginx_ingress=true必填，nginx ingress访问方式：external（公网），internal（内网），当选择公网时将自动创建eip额外产生eip相关费用
- `pod_cidr` (String) pod网络cidr，使用cubecni作为网络插件时，podCidr不填，服务端会取vpcCidr。使用calico作为网络插件时，podCidr与vpcCidr和serviceCidr不能重叠。
- `pod_subnet_id_list` (Set of String) pod子网ID列表，网络插件选择cubecni必传，需要属于所选VPC，最多支持10个子网
- `project_id` (String) 企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID
- `security_group_id` (String) 安全组ID，需属于所选vpc。使用自定义安全组时，需要配置如下规则，参考<a href="https://www.ctyun.cn/document/10083472/10915714">集群安全组规则配置</a>
- `series_type` (String) 托管版集群规格，托管版集群必填。支持managedbase（单实例），managedpro（多实例）。单/多实例指控制面是否高可用，生产环境建议使用多实例
- `service_cidr` (String) 服务cidr，默认10.96.0.0/16。网络插件为calico时，podCidr与vpcCidr与serviceCidr不能重叠。选择cubecni时，podCidr（vpcCidr）与serviceCidr不能重叠。


<a id="nestedatt--slave_host"></a>
### Nested Schema for `slave_host`

Required:

- `az_infos` (Attributes List) 可用区信息，包括可用区编码和该可用区下worker节点数量，支持的可用区可通过ctyun_regions查询 (see [below for nested schema](#nestedatt--slave_host--az_infos))
- `instance_type` (String) 实例类型，支持ecs（云主机）、ebm（裸金属）
- `item_def_name` (String) 实例规格名称，使用至少4C8G以上的规格，云主机规格通过ctyun_ecs_flavors查询，裸金属规格通过ctyun_ebm_device_types查询
- `mirror_type` (Number) 镜像类型，支持传0（私有），1（公有），可查看<a href="https://www.ctyun.cn/document/10026730/10030151">镜像概述</a>

Optional:

- `data_disks` (Attributes List) 数据盘信息 (see [below for nested schema](#nestedatt--slave_host--data_disks))
- `mirror_id` (String) 镜像id，worker节点为ecs类型必填，可查看<a href="https://www.ctyun.cn/document/10083472/11004475">节点规格和节点镜像</a>
- `mirror_name` (String) 镜像名称，worker节点为ebm类型必填，可查看<a href="https://www.ctyun.cn/document/10083472/11004475">节点规格和节点镜像</a>
- `sys_disk` (Attributes) 系统盘信息 (see [below for nested schema](#nestedatt--slave_host--sys_disk))

<a id="nestedatt--slave_host--az_infos"></a>
### Nested Schema for `slave_host.az_infos`

Required:

- `az_name` (String) worker可用区编码
- `size` (Number) 该可用区下worker节点数量


<a id="nestedatt--slave_host--data_disks"></a>
### Nested Schema for `slave_host.data_disks`

Required:

- `size` (Number) 数据盘大小，单位为G，支持范围10-20000
- `type` (String) 数据盘类型，支持SATA、SAS、SSD


<a id="nestedatt--slave_host--sys_disk"></a>
### Nested Schema for `slave_host.sys_disk`

Required:

- `size` (Number) 系统盘大小，单位为G，支持范围80-2040
- `type` (String) 系统盘类型，支持SATA、SAS、SSD



<a id="nestedatt--master_host"></a>
### Nested Schema for `master_host`

Required:

- `az_infos` (Attributes List) 可用区信息，包括可用区编码和该可用区下master节点数量，支持的可用区可通过ctyun_regions查询 (see [below for nested schema](#nestedatt--master_host--az_infos))
- `item_def_name` (String) 实例规格名称，使用至少4C8G以上的规格，仅支持云主机，可通过ctyun_ecs_flavors查询

Optional:

- `data_disks` (Attributes List) 数据盘信息 (see [below for nested schema](#nestedatt--master_host--data_disks))
- `sys_disk` (Attributes) 系统盘信息 (see [below for nested schema](#nestedatt--master_host--sys_disk))

<a id="nestedatt--master_host--az_infos"></a>
### Nested Schema for `master_host.az_infos`

Required:

- `az_name` (String) master可用区编码
- `size` (Number) 该可用区下master节点数量


<a id="nestedatt--master_host--data_disks"></a>
### Nested Schema for `master_host.data_disks`

Required:

- `size` (Number) 数据盘大小，单位为G，支持范围10-20000
- `type` (String) 数据盘类型，支持SATA、SAS、SSD


<a id="nestedatt--master_host--sys_disk"></a>
### Nested Schema for `master_host.sys_disk`

Required:

- `size` (Number) 系统盘大小，单位为G，支持范围80-2040
- `type` (String) 系统盘类型，支持SATA、SAS、SSD