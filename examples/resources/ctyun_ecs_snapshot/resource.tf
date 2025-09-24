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

resource "ctyun_ecs" "jdxutuzpfr" {
  instance_name      = "tf-test-ecs"
  display_name       = "tf-test-init-ecs"
  flavor_id          = "9b4b5e39-db25-f2c8-3914-76881ee77d5c"
  image_id           = "fa3f3784-34f9-4f6b-80a1-dd173d486bd6"
  system_disk_type   = "sata"
  system_disk_size   = 60
  vpc_id             = "vpc-0ein2p8bs8"
  subnet_id          = "subnet-0oiyrpu8nk"
  key_pair_name      = "tf-keypair-for-ecs"
  cycle_type         = "on_demand"
}

resource "ctyun_ecs_snapshot" "test" {
  name = "tf-test-group"
  instance_id = ctyun_ecs.jdxutuzpfr.id
}
