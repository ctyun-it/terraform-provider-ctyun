terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}


module "vpc" {
  source   = "./modules/vpc"
  vpc_name = "vpc-${var.name}"
}

module "subnet" {
  source      = "./modules/subnet"
  subnet_name = "subnet-${var.name}"
  vpc_id      = module.vpc.vpc_id
}

module "security_group" {
  source              = "./modules/security_group"
  security_group_name = "security-group-${var.name}"
  vpc_id              = module.vpc.vpc_id
}

module "security_group_rules" {
  source            = "./modules/security_group_rules"
  security_group_id = module.security_group.security_group_id
}

module "key_pair" {
  source        = "./modules/key_pair"
  key_pair_name = "keypair-${var.name}"
  public_key    = var.public_key
}