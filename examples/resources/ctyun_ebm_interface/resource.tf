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

resource "ctyun_ebm_interface" "test" {
  security_group_ids = ["sg-t0ae11aig1"]
  instance_id = "ss-uadmwtxinfp4tkbhvwp52vnzl2kn"
  ipv4 = "192.168.0.13"
  subnet_id = "subnet-43z7cqmjlp"
}
