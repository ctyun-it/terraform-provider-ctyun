output "vpc_id" {
  value = ctyun_vpc.vpc_test.id
}

output "subnet_id" {
  value = ctyun_subnet.subnet_test.id
}

output "security_group_id" {
  value = ctyun_security_group.security_group_test.id
}

output "image_id" {
  value = data.ctyun_images.image_test.images[0].id
}

output "flavor_id" {
  value = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].id
}

output "flavor_id2" {
  value = data.ctyun_ecs_flavors.ecs_flavor_test2.flavors[0].id
}

output "affinity_group_id" {
  value = ctyun_ecs_affinity_group.affinity_group_test.id
}

output "key_pair_name" {
  value = ctyun_keypair.keypair_test.name
}

output "key_pair_name2" {
  value = ctyun_keypair.keypair_test2.name
}

output "ecs_id" {
  value = ctyun_ecs.ecs_test.id
}