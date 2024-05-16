terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

provider "ctyun" {
  ak        = var.ak
  sk        = var.sk
  region_id = var.region_id
  az_name   = var.az_name
  env       = "prod"
}

module "default_env" {
  source     = "./modules/default_env"
  name       = "poc-test-5"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC2oDo+G8IoXwTguiQbFxxi43Ru7Kz1CnshBJyzf2erBz3YTPKx90afHAIl4eWF9sqBm6lFo+gF3rm4/WNObQA2OXM3+ZDnVcfej4O1b/C2DFKlX43Em4ejejzlfezf3Ai120Rwd9qqCAOojaEBl9fbEvVPrNzZ7eiuSL/GSsJ4jPmwG358eIdZX78jIb4QENAllUOufRLRCLpjoQFWRijK7Y10C9EY6uQn2a5CuXhWrZblzl3wOKvh+5zJzj2Om7NICEbxW2zrTAojYAbkpftnkcr4zh+H6gI9czMupJSOWEOyPhienJbiJ/H6XnlTmPCu3fxU4Eg7kwR6+N98AjWvtonQxKsOuvdXgQkwi7gdvpsCLJlx5vpamx5MICvoIAsanDOoInL/bUXUvmVxEmL5vLlHSpOv+gntpVShtOXvqkaeHkC9gq8crLzuKameZyn7bAzwa1NuBQe2ceqnYM8iFkQ88HWgf5fFpbcyMb9jRQSNuuktan9FFryxs9ib9/s= 44872@MinChiang"
}

module "mysql_config" {
  source                           = "./modules/mysql_config"
  mysql_vpc_id                     = module.default_env.default_vpc_id
  mysql_subnet_id                  = module.default_env.default_subnet_id
  mysql_default_security_group_ids = [module.default_env.default_security_group_id]
  mysql_security_group_name        = "mysql-poc-5"
  mysql_ecs_name                   = "mysql-poc-ecs-5"
  mysql_ecs_password               = "Denny__22222aa"

  mysql_eip_bandwidth   = 5
  mysql_eip_cycle_count = 1
  mysql_eip_cycle_type  = "month"
  mysql_eip_name        = "mysql-poc-eip-5"


  mysql_bandwidth_bandwidth   = 50
  mysql_bandwidth_cycle_type  = "on_demand"
  mysql_bandwidth_name        = "mysql-poc-bandwidth-5"
  mysql_bandwidth_cycle_count = null

  mysql_ebs_cycle_count = 1
  mysql_ebs_cycle_type  = "month"
  mysql_ebs_mode        = "VBD"
  mysql_ebs_name        = "mysql-poc-ebs-5"
  mysql_ebs_size        = 100
  mysql_ebs_type        = "SATA"
  mysql_key_pair_name   = module.default_env.default_key_pair_name
}
