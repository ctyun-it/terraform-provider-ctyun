terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
}

resource "ctyun_bandwidth" "bandwidth_test" {
  cycle_type = var.cycle_type
  cycle_count = var.cycle_count
  bandwidth = var.bandwidth
  name = var.name
}

resource "ctyun_bandwidth_association_eip" "bandwidth_association_test" {
    count = var.bandwidth_association_eip_count
    bandwidth_id = ctyun_bandwidth.bandwidth_test.id
    eip_id = var.eip_ids[count.index]
}