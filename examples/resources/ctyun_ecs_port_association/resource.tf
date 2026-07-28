terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

provider "ctyun" {
  env = "prod"
}

resource "ctyun_port" "port" {
  name               = "port_for_association_test"
  subnet_id          = "subnet-a0c0c0c011"
  security_group_ids = ["sg-xxxxxxxxx"]
}

resource "ctyun_ecs_port_association" "ecs_port_for_association_test" {
  instance_id = "ae432721-61bf-45b7-b207-7e3256c1c2d6"
  port_id     = ctyun_port.port.id
}



