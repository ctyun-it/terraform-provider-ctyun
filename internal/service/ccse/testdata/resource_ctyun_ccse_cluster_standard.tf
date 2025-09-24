resource "ctyun_ccse_cluster" "%[1]s" {
  base_info = {
    cluster_name = "%[2]s"
    cluster_series = "%[3]s"
    vpc_id     = "%[4]s"
    subnet_id  = "%[5]s"
    cluster_domain = "www.ctyun.com"
    network_plugin = "calico"
    pod_cidr = "172.26.0.0/16"
    pod_subnet_id_list = ["%[5]s"]
    start_port = 30000
    end_port   = 65535
    elb_prod_code = "standardI"
    cycle_type  = "on_demand"
    container_runtime = "containerd"
    timezone    = "Asia/Shanghai"
    cluster_version = "1.29.3"
    deploy_type   = "single"
    kube_proxy    = "iptables"
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
    item_def_name =  "%[6]s"

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
    item_def_name = "%[6]s"

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
