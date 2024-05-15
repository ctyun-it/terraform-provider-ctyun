variable "security_group_ids" {
  description = "安全组id"
  type        = set(string)
}

variable "ecs_name" {
  description = "ecs名称"
}

variable "vpc_id" {
  description = "vpcid"
}

variable "subnet_id" {
  description = "子网id"
}

variable "ecs_password" {
  description = "密码"
}

variable "ecs_instance_count" {
  default     = 3
  description = "mysql集群大小"
}

variable "ecs_key_pair_name" {
  description = "密钥对名称"
}