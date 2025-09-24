resource "ctyun_ccse_cluster" "%[1]s" {
  base_info = {
    cluster_name = "%[2]s"
    cluster_series = "%[3]s"
    vpc_id     = "%[4]s"
    subnet_id  = "%[5]s"
    cluster_domain = "www.ctyun.com"
    network_plugin = "cubecni"
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
    series_type = "managedbase"
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
