resource "ctyun_search_instance" "%[1]s" {
  name      = "%[2]s"
  region_id         = "%[3]s"
  zone_list = ["%[4]s"]
  vpc_id            = "%[5]s"
  subnet_id         = "%[6]s"
  security_group_id = "%[7]s"
  enable_ipv6       = %[8]t
  password     = "%[9]s"
  cluster_type      = %[10]d
  os_type           = "%[11]s"
  enable_https      = "%[12]s"
  cycle_count         = %[13]d
  cycle_type        = "%[14]s"

  node_details = [
    {
    host_num          = %[15]d
    storage_type           = "%[16]s"
    storage_space            = %[17]d
    flavor_name = "%[18]s"
      node_group_type   = "%[19]s"
    },
    {
    host_num          = 3
    storage_type           = "%[16]s"
    storage_space            = %[17]d
    flavor_name = "%[18]s"
    node_group_type   = "%[20]s"
    },
    {
    host_num          = %[15]d
    storage_type           = "%[16]s"
    storage_space            = %[17]d
    flavor_name = "%[18]s"
    node_group_type   = "%[21]s"
    },
    {
    host_num          = %[15]d
    storage_type           = "%[16]s"
    storage_space            = %[17]d
    flavor_name = "%[18]s"
    node_group_type   = "%[22]s"
    }
  ]

}
