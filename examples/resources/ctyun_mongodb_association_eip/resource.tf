terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 可参考index.md，在环境变量中配置ak、sk、资源池ID、可用区名称
provider "ctyun" {

}

resource "ctyun_mongodb_association_eip" "test" {
  eip_id = "eip-xjw2ndksn3"
  inst_id = "85d8cb5914ae4b6b9852635f3bc43023"
  host_ip = "192.168.1.2"
}
