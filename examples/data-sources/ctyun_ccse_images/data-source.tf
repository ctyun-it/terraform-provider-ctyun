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

data "ctyun_ecs_flavors" "test" {
  cpu    = 4
  ram    = 8
  arch   = "x86"
  series = "C"
  type   = "CPU_C7"
}

# 云主机
data "ctyun_ccse_images" "ccse_ecs_image" {
  instance_type = "ecs"
  flavor_name   = data.ctyun_ecs_flavors.test.flavors[0].name
}

output "image1" {
  value = data.ctyun_ccse_images.ccse_ecs_image
}


# 物理机
data "ctyun_ccse_images" "ccse_ebm_image" {
  instance_type = "ebm"
  flavor_name   = "physical.s5.2xlarge4"
}

output "image2" {
  value = data.ctyun_ccse_images.ccse_ebm_image
}