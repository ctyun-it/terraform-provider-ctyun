variable "vpc_name" {
  description = "vpc名称"
}

variable "vpc_cidr" {
  default     = "10.0.0.0/8"
  description = "vpc网段"
}