output "vpc_id" {
  value = local.real_vpc_id
}

output "subnet_id" {
  value =  local.real_subnet_id
}

output "security_group_id" {
  value = local.real_security_group_id
}

output "eip_address" {
  value = ctyun_eip.eip_test.address
}

output "eip_id" {
  value = ctyun_eip.eip_test.id
}

output "redis_version" {
  value = local.spec.version
}

output "redis_engine_edition" {
  value = local.spec.series_code
}

output "redis_instance_id" {
  value = ctyun_redis_instance.test_redis_instance.id
}
output "redis_instance2_id" {
  value = ctyun_redis_instance.test_redis_instance2.id
}
output "redis_address" {
  value = ctyun_redis_instance.test_redis_instance.connection_address
}
output "redis2_address" {
  value = ctyun_redis_instance.test_redis_instance2.connection_address
}

output "instance_account_name" {
  value = ctyun_redis_account.test_instance1_account.name
}
output "instance_account_pswd" {
  value = ctyun_redis_account.test_instance1_account.password
  sensitive = true
}

output "instance2_account_name" {
  value = ctyun_redis_account.test_instance2_account.name
}

output "instance2_account_pswd" {
  value = ctyun_redis_account.test_instance2_account.password
  sensitive = true
}