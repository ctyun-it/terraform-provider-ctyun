terraform {
  required_providers {
    ctyun = {
      source = "www.ctyun.cn/ctyun/ctyun"
    }
  }
}

module "security_group" {
  source              = "./modules/security_group"
  security_group_name = var.mysql_security_group_name
  vpc_id              = var.mysql_vpc_id
}

module "security_group_rules" {
  source            = "./modules/security_group_rules"
  security_group_id = module.security_group.security_group_id
}

module "ecs" {
  source             = "./modules/ecs"
  vpc_id             = var.mysql_vpc_id
  subnet_id          = var.mysql_subnet_id
  security_group_ids = setunion([module.security_group.security_group_id], var.mysql_default_security_group_ids)
  ecs_name           = var.mysql_ecs_name
  ecs_password       = var.mysql_ecs_password
  ecs_instance_count = var.mysql_ecs_count
  ecs_key_pair_name  = var.mysql_key_pair_name
}

module "eip" {
  source             = "./modules/eip"
  bandwidth          = var.mysql_eip_bandwidth
  cycle_count        = var.mysql_eip_cycle_count
  cycle_type         = var.mysql_eip_cycle_type
  name               = var.mysql_eip_name
  eip_instance_count = var.mysql_ecs_count
  instance_ids       = module.ecs.ecs_instance_id
}

module "bandwidth" {
  source                          = "./modules/bandwidth"
  bandwidth                       = var.mysql_bandwidth_bandwidth
  cycle_count                     = var.mysql_bandwidth_cycle_count
  cycle_type                      = var.mysql_bandwidth_cycle_type
  name                            = var.mysql_bandwidth_name
  bandwidth_association_eip_count = var.mysql_ecs_count
  eip_ids                         = module.eip.eip_instance_id
}

module "ebs" {
  source                    = "./modules/ebs"
  cycle_count               = var.mysql_ebs_cycle_count
  cycle_type                = var.mysql_ebs_cycle_type
  mode                      = var.mysql_ebs_mode
  name                      = var.mysql_ebs_name
  size                      = var.mysql_ebs_size
  type                      = var.mysql_ebs_type
  ebs_association_ecs_count = var.mysql_ecs_count
  instance_ids              = module.ecs.ecs_instance_id
}