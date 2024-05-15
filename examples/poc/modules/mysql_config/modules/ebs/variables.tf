variable "name" {
  description = "磁盘命名，单账户单资源池下，命名需唯一"
}

variable "cycle_type" {
  description = "订购周期类型，与cycle_count配合使用，month：按月，year：按年，on_demand：按需计费类型。当此值为month或者year时，cycle_count为必填"
}

variable "cycle_count" {
  description = "包周期数。周期最大长度不能超过5年"
  nullable = true
  type = number
}

variable "mode" {
  description = "磁盘模式，VBD/ISCSI/FCSAN"
}

variable "type" {
  description = "磁盘类型，SATA：普通IO，SAS：高IO，SSD-genric：通用型SSD，SSD：超高IO，FAST-SSD：极速型SSD"
}

variable "size" {
  description = "磁盘大小，单位GB，取值范围[10, 32768]"
  type = number
}

variable "ebs_instance_count" {
  default     = 3
  description = "ebs数量"
}

variable "ebs_association_ecs_count" {
  description = "主机数量"
  type = number
}

variable "instance_ids" {
  description = "挂载实例id"
  type = list(string)
}