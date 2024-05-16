terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

resource "ctyun_eip" "eip_test" {
  count       = var.eip_instance_count
  name        = "${var.name}-${count.index}"
  cycle_type  = var.cycle_type
  cycle_count = var.cycle_count
  bandwidth   = var.bandwidth
}

resource "ctyun_eip_association" "eip_association_test" {
  count = var.eip_association_count
  eip_id         = ctyun_eip.eip_test[count.index].id
  association_id = var.instance_ids[count.index]
}