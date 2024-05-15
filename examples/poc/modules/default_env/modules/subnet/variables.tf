variable "vpc_id" {
  description = "vpcid"
}

variable "subnet_name" {
  description = "子网名称"
}

variable "subnet_cidr" {
  default     = "10.0.0.0/8"
  description = "子网网段"
}

variable "subnet_dns" {
  default = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
  description = "子网DNS"
}