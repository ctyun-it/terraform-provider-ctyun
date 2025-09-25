variable "additional_routes" {
  description = "需要添加到路由表的路由规则"
  type = list(object({
    description = string
    cidr_block  = string
    gateway_id  = optional(string)
  }))
  default = [
    {
      description = "资源池CN2 1107, 10.11.0.0/16"
      cidr_block  = "10.11.0.0/16"
    },
    {
      description = "资源池CN2 1107, 10.12.0.0/16"
      cidr_block  = "10.12.0.0/16"
    },
    {
      description = "资源池CN2 1107, 10.13.0.0/16"
      cidr_block  = "10.13.0.0/16"
    },
    {
      description = "资源池CN2 1107, 10.14.0.0/16"
      cidr_block  = "10.14.0.0/16"
    },
    {
      description = "资源池CN2 1107, 10.16.0.0/16"
      cidr_block  = "10.16.0.0/16"
    },
    {
      description = "资源池CN2 1107, 10.17.0.0/16"
      cidr_block  = "10.17.0.0/16"
    },
    {
      description = "资源池CN2 1107, 10.18.0.0/16"
      cidr_block  = "10.18.0.0/16"
    },
    {
      description = "资源池CN2 1107, 10.19.0.0/16"
      cidr_block  = "10.19.0.0/16"
    },
    {
      description = "资源池CN2 1107, 10.20.0.0/16"
      cidr_block  = "10.20.0.0/16"
    },
    {
      description = "资源池CN2 1107, 10.245.0.0/16"
      cidr_block  = "10.245.0.0/16"
    },
    {
      description = "资源池CN2 1107, 10.249.0.0/16"
      cidr_block  = "10.249.0.0/16"
    },
    {
      description = "资源池CN2 1107, 10.251.0.0/16"
      cidr_block  = "10.251.0.0/16"
    },
    {
      description = "资源池CN2 1107, 10.252.0.0/16"
      cidr_block  = "10.252.0.0/16"
    },
    {
      description = "资源池CN2 1107, 10.253.0.0/16"
      cidr_block  = "10.253.0.0/16"
    },
    {
      description = "资源池CN2 1107, 10.254.0.0/16"
      cidr_block  = "10.254.0.0/16"
    },
    {
      description = "资源池CN2 1107, 10.255.0.0/16"
      cidr_block  = "10.255.0.0/16"
    }
  ]
}

