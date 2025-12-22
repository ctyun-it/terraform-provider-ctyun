output "vpc_id" {
  value = ctyun_vpc.vpc_test.id
}

output "vpc_id1" {
  value = ctyun_vpc.vpc_test1.id
}

# output "vpc_id2" {
#   value = ctyun_vpc.vpc_test2.id
# }

output "peer_connection_id" {
  value = ctyun_vpc_peer_connection.test.id
}

output "rtb_id" {
  value = data.ctyun_vpc_route_tables.route_table_test.route_tables[0].route_table_id
}