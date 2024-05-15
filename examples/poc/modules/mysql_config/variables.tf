variable "mysql_vpc_id" {
  description = "使用的vpc"
}

variable "mysql_security_group_name" {
  description = "mysql默认加入安全组的名称"
}

variable "mysql_default_security_group_ids" {
  type        = set(string)
  description = "mysql默认加入安全组列表"
}

variable "mysql_ecs_name" {
  description = "ecs名称"
}

variable "mysql_subnet_id" {
  description = "子网id"
}

variable "mysql_ecs_password" {
  default     = "MySqlP@ssw0rd__"
  description = "密码"
}

variable "mysql_ecs_count" {
  default     = 3
  description = "mysql集群大小"
}

variable "mysql_eip_name" {
  description = "EIP名称"
}

variable "mysql_eip_cycle_type" {
  description = "EIP订购周期类型，与cycle_count配合使用，month：按月，year：按年，on_demand：按需计费类型。当此值为month或者year时，cycle_count为必填"
}

variable "mysql_eip_cycle_count" {
  description = "EIP包周期数"
  nullable    = true
}

variable "mysql_eip_bandwidth" {
  description = "EIP带宽大小"
}

variable "mysql_bandwidth_name" {
  description = "共享带宽名称"
}

variable "mysql_bandwidth_cycle_type" {
  description = "共享带宽订购周期类型，与cycle_count配合使用，month：按月，year：按年，on_demand：按需计费类型。当此值为month或者year时，cycle_count为必填"
}

variable "mysql_bandwidth_cycle_count" {
  description = "共享带宽包周期数"
  nullable    = true
  type        = number
}

variable "mysql_bandwidth_bandwidth" {
  description = "共享带宽大小"
}

variable "mysql_ebs_cycle_count" {
  description = "磁盘包周期数。周期最大长度不能超过5年"
  nullable    = true
  type        = number
}

variable "mysql_ebs_cycle_type" {
  description = "磁盘订购周期类型，与cycle_count配合使用，month：按月，year：按年，on_demand：按需计费类型"
}

variable "mysql_ebs_name" {
  description = "磁盘命名，单账户单资源池下，命名需唯一"
}

variable "mysql_ebs_mode" {
  description = "磁盘模式，VBD/ISCSI/FCSAN"
}

variable "mysql_ebs_size" {
  description = "磁盘大小，单位GB，取值范围[10, 32768]"
  nullable    = true
  type        = number
}

variable "mysql_ebs_type" {
  description = "磁盘类型，SATA：普通IO，SAS：高IO，SSD-genric：通用型SSD，SSD：超高IO，FAST-SSD：极速型SSD"
}

variable "mysql_key_pair_name" {
  description = "mysql公钥"
}
