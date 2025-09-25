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

resource "ctyun_kafka_instance" "tbidgqvfbs" {
  instance_name = "tf-kafka-34ywerkb"
  engine_version = "3.6"
  spec_name = "kafka.4u8g.cluster"
  node_num = 5
  zone_list = ["cn-huadong1-jsnj1A-public-ctcloud"]
  disk_type = "SSD"
  disk_size = 300
  vpc_id = "vpc-ewivt5nhiz"
  subnet_id = "subnet-vhyywu7mfe"
  security_group_id = "sg-ed9i3c98t2"
  cycle_type = "on_demand"
  retention_hours = 55
}
