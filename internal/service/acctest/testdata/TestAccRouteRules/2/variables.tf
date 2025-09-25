variable "additional_routes" {
  description = "需要添加到路由表的路由规则"
  type = list(object({
    description = string
    cidr_block  = string
    gateway_id  = optional(string)
  }))
  default = [
    {
      description = "资源池CN2 1107, 30.0.0.0/16"
      cidr_block  = "30.0.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.100.0.0/16"
      cidr_block  = "30.100.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.112.0.0/16"
      cidr_block  = "30.112.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.115.0.0/16"
      cidr_block  = "30.115.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.116.0.0/16"
      cidr_block  = "30.116.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.117.0.0/16"
      cidr_block  = "30.117.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.120.0.0/16"
      cidr_block  = "30.120.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.124.0.0/16"
      cidr_block  = "30.124.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.128.0.0/16"
      cidr_block  = "30.128.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.132.0.0/16"
      cidr_block  = "30.132.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.136.0.0/16"
      cidr_block  = "30.136.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.144.0.0/16"
      cidr_block  = "30.144.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.16.0.0/16"
      cidr_block  = "30.16.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.2.0.0/16"
      cidr_block  = "30.2.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.20.0.0/16"
      cidr_block  = "30.20.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.24.138.0/24"
      cidr_block  = "30.24.138.0/24"
    },
    {
      description = "资源池CN2 1107, 30.24.146.0/24"
      cidr_block  = "30.24.146.0/24"
    },
    {
      description = "资源池CN2 1107, 30.24.233.0/24"
      cidr_block  = "30.24.233.0/24"
    },
    {
      description = "资源池CN2 1107, 30.28.0.0/16"
      cidr_block  = "30.28.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.31.0.0/16"
      cidr_block  = "30.31.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.32.0.0/16"
      cidr_block  = "30.32.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.36.0.0/16"
      cidr_block  = "30.36.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.39.0.0/16"
      cidr_block  = "30.39.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.40.0.0/16"
      cidr_block  = "30.40.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.44.0.0/16"
      cidr_block  = "30.44.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.48.0.0/16"
      cidr_block  = "30.48.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.55.0.0/16"
      cidr_block  = "30.55.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.59.0.0/16"
      cidr_block  = "30.59.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.60.0.0/16"
      cidr_block  = "30.60.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.64.0.0/16"
      cidr_block  = "30.64.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.68.0.0/16"
      cidr_block  = "30.68.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.71.0.0/16"
      cidr_block  = "30.71.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.72.0.0/16"
      cidr_block  = "30.72.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.75.0.0/16"
      cidr_block  = "30.75.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.76.0.0/16"
      cidr_block  = "30.76.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.8.6.0/24"
      cidr_block  = "30.8.6.0/24"
    },
    {
      description = "资源池CN2 1107, 30.8.13.0/24"
      cidr_block  = "30.8.13.0/24"
    },
    {
      description = "资源池CN2 1107, 30.8.28.0/22"
      cidr_block  = "30.8.28.0/22"
    },
    {
      description = "资源池CN2 1107, 30.8.45.0/24"
      cidr_block  = "30.8.45.0/24"
    },
    {
      description = "资源池CN2 1107, 30.8.51.0/24"
      cidr_block  = "30.8.51.0/24"
    },
    {
      description = "资源池CN2 1107, 30.8.201.0/24"
      cidr_block  = "30.8.201.0/24"
    },
    {
      description = "资源池CN2 1107, 30.80.0.0/16"
      cidr_block  = "30.80.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.81.0.0/16"
      cidr_block  = "30.81.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.82.0.0/16"
      cidr_block  = "30.82.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.83.0.0/16"
      cidr_block  = "30.83.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.87.0.0/16"
      cidr_block  = "30.87.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.88.0.0/16"
      cidr_block  = "30.88.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.89.0.0/16"
      cidr_block  = "30.89.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.90.0.0/16"
      cidr_block  = "30.90.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.91.0.0/16"
      cidr_block  = "30.91.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.92.0.0/16"
      cidr_block  = "30.92.0.0/16"
    },
    {
      description = "资源池CN2 1107, 30.96.0.0/16"
      cidr_block  = "30.96.0.0/16"
    }
  ]
}

