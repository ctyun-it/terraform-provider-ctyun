variable "name" {
  description = "eip名称"
}

variable "cycle_type" {
  description = "订购周期类型，与cycle_count配合使用，month：按月，year：按年，on_demand：按需计费类型。当此值为month或者year时，cycle_count为必填"
}

variable "cycle_count" {
  description = "包周期数"
  nullable = true
  type = number
}

variable "bandwidth" {
  description = "带宽大小"
  type = number
}

variable "eip_instance_count" {
  default     = 3
  description = "eip数量"
}

variable "eip_association_count" {
  default     = 3
  description = "eip挂载实例数量"
}

variable "instance_ids" {
  description = "eip绑定实例id"
  type = list(string)
}