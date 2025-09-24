resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-tongyiwangguan"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
}

data "ctyun_vpc_route_tables" "rtest" {
  vpc_id = ctyun_vpc.vpc_test.id
}

locals {
  default_route_table_id = data.ctyun_vpc_route_tables.rtest.route_tables[0].route_table_id
}

resource "ctyun_vpc_route_table_rule" "rule_test"{
  destination = "0.0.0.0/0"
  next_hop_id = ctyun_nat.nat_test.id
  next_hop_type = "natgw"
  route_table_id = local.default_route_table_id
  ip_version = 4
}

# 创建子网
resource "ctyun_subnet" "subnet_test" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-tongyiwangguan"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
}

# 创建eip
resource "ctyun_eip" "eip_test" {
  name                = "tf-eip-test"
  bandwidth           = 5
  cycle_type          = "on_demand"
  demand_billing_type = "upflowc"
}

# 创建nat网关
resource "ctyun_nat" "nat_test"{
  vpc_id = ctyun_vpc.vpc_test.id
  spec = 1
  name = "tf-nat-unified-gateway"
  cycle_type    = "month"
  cycle_count   = 1
}

resource "ctyun_nat_snat" "snat_test"{
  nat_gateway_id = ctyun_nat.nat_test.id
  source_subnet_id = ctyun_subnet.subnet_test.id
  snat_ips = [ctyun_eip.eip_test.id]
}

# 创建云主机
data "ctyun_images" "image_test" {
  name       = "CentOS Linux 8.4"
  visibility = "public"
  page_no = 1
  page_size = 10
}

data "ctyun_ecs_flavors" "ecs_flavor_test" {
  cpu    = 2
  ram    = 4
  arch   = "x86"
  series = "C"
  type   = "CPU_C7"
}

resource "ctyun_keypair" "keypair_test" {
  name       = "tf-keypair-for-123"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjUnAnTid4wmVtajSmElMtH03OvOyY81ybfswbUu9Gt83DVVzDnwb3rcQW1us8SeKm/gRINkgdrRAgfXAmTKR7AorYtWWc/tzb6kcDpL2E8Qk+n6cyFAxXNoX2vXBr4kC9wz1uwjGyxoSlpHLIpscfI0Ef652gMlSyfODehAJHj3JPMr8pvtPIUqsZI3JOGTUzxaA2JVC0LxQegphYYf2TxGd9GLRUv1p/0BUAPCMg1NaITXNVEj3A11hk1nrFoJMmvIwIUkLmRuQcxuNAdxeLB7GXXVjKpnKIJL4L64dyA9GWa3Gb7gCJyRaBc5UhK4hT57wmukCrldHHtdF1IJr"
}

resource "ctyun_ecs" "ecs_test" {
  instance_name       = "tf-ecs-for-0627-1"
  display_name        = "tf-ecs-for-0627-1"
  flavor_id           = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].id
  image_id            = data.ctyun_images.image_test.images[0].id
  system_disk_type    = "sata"
  system_disk_size    = 40
  vpc_id = ctyun_vpc.vpc_test.id
  key_pair_name       = ctyun_keypair.keypair_test.name
  cycle_type          = "on_demand"
  subnet_id = ctyun_subnet.subnet_test.id
}

resource "ctyun_ecs" "ecs_test2" {
  instance_name       = "tf-ecs-for-0627-2"
  display_name        = "tf-ecs-for-0627-2"
  flavor_id           = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].id
  image_id            = data.ctyun_images.image_test.images[0].id
  system_disk_type    = "sata"
  system_disk_size    = 40
  vpc_id = ctyun_vpc.vpc_test.id
  key_pair_name       = ctyun_keypair.keypair_test.name
  cycle_type          = "on_demand"
  subnet_id = ctyun_subnet.subnet_test.id
}
