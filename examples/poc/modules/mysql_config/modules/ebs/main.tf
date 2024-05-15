terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
}

resource "ctyun_ebs" "ebs_test" {
  count = var.ebs_instance_count
  name = "${var.name}-${count.index}"
  mode = var.mode
  type = var.type
  size = var.size
  cycle_type = var.cycle_type
  cycle_count = var.cycle_count
}

resource "ctyun_ebs_association_ecs" "ebs_association_test" {
    count = var.ebs_association_ecs_count
    disk_id = ctyun_ebs.ebs_test[count.index].id
    instance_id = var.instance_ids[count.index]
}