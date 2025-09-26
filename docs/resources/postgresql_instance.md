# ctyun_postgresql_instance (Resource)
-> 详细说明请见文档：https://www.ctyun.cn/document/10034019/10153165



## Example

```terraform
terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 可参考index.md，在环境变量中配置ak、sk、资源池ID、可用区名称
provider "ctyun" {
  env = "prod"
}

variable "password" {
  type      = string
  sensitive = true
}

resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-pgsql"
  cidr        = "192.168.0.0/16"
  description = "terraform-kafka测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-pgsql"
  cidr        = "192.168.1.0/24"
  description = "terraform-kafka测试使用"
  dns = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
}
resource "ctyun_security_group" "sg_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-esc"
  description = "terraform-kafka测试使用"
  lifecycle {
    prevent_destroy = false
  }
}
// 开通样例
resource "ctyun_postgresql_instance" "test" {
  cycle_type            = "on_demand"
  host_type             = "S7"
  prod_id               = "Single1222"
  storage_type          = "SATA"
  storage_space         = 100
  name                  = "pgsql-test"
  password              = var.password
  case_sensitive        = true
  instance_series       = "S"
  prod_performance_spec = "2C4G"
  vpc_id                = ctyun_vpc.vpc_test.id
  subnet_id             = ctyun_subnet.subnet_test.id
  security_group_id     = ctyun_security_group.sg_test.id
  availability_zone_info = [
    { "availability_zone_name" : "cn-gs-qyi2-1a-public-ctcloud", "availability_zone_count" : 1, "node_type" : "master" }
  ] // availability_zone_name值根据情况而定
  backup_storage_type   = "SATA"
  backup_storage_space  = 100
  os_type               = "ctyunos"
  cpu_type              = "Intel"
}

// 升配pgsql--对备用磁盘扩容(在升配主storage时候，确保备用磁盘空间>主磁盘空间)
resource "ctyun_postgresql_instance" "test" {
  cycle_type            = "on_demand"
  host_type             = "S7"
  prod_id               = "Single1222"
  storage_type          = "SATA"
  storage_space         = 100
  name                  = "pgsql-test"
  password              = var.password
  case_sensitive        = true
  instance_series       = "S"
  prod_performance_spec = "2C4G"
  vpc_id                = ctyun_vpc.vpc_test.id
  subnet_id             = ctyun_subnet.subnet_test.id
  security_group_id     = ctyun_security_group.sg_test.id
  availability_zone_info = [
    { "availability_zone_name" : "cn-gs-qyi2-1a-public-ctcloud", "availability_zone_count" : 1, "node_type" : "master" }
  ] // availability_zone_name值根据情况而定
  backup_storage_type   = "SATA"
  backup_storage_space  = 120
  os_type               = "ctyunos"
  cpu_type              = "Intel"
}


// 升配pgsql--对主磁盘扩容(在升配主storage时候，确保备用磁盘空间>主磁盘空间)
resource "ctyun_postgresql_instance" "test" {
  cycle_type            = "on_demand"
  host_type             = "S7"
  prod_id               = "Single1222"
  storage_type          = "SATA"
  storage_space         = 120
  name                  = "pgsql-test"
  password              = var.password
  case_sensitive        = true
  instance_series       = "S"
  prod_performance_spec = "2C4G"
  vpc_id                = ctyun_vpc.vpc_test.id
  subnet_id             = ctyun_subnet.subnet_test.id
  security_group_id     = ctyun_security_group.sg_test.id
  availability_zone_info = [
    { "availability_zone_name" : "cn-gs-qyi2-1a-public-ctcloud", "availability_zone_count" : 1, "node_type" : "master" }
  ] // availability_zone_name值根据情况而定
  backup_storage_type   = "SATA"
  backup_storage_space  = 120
  os_type               = "ctyunos"
  cpu_type              = "Intel"
}


// 升配规格 2C4G->2C8G
resource "ctyun_postgresql_instance" "test" {
  cycle_type            = "on_demand"
  host_type             = "S7"
  prod_id               = "Single1222"
  storage_type          = "SATA"
  storage_space         = 120
  name                  = "pgsql-test"
  password              = var.password
  case_sensitive        = true
  instance_series       = "S"
  prod_performance_spec = "2C8G"
  vpc_id                = ctyun_vpc.vpc_test.id
  subnet_id             = ctyun_subnet.subnet_test.id
  security_group_id     = ctyun_security_group.sg_test.id
  availability_zone_info = [
    { "availability_zone_name" : "cn-gs-qyi2-1a-public-ctcloud", "availability_zone_count" : 1, "node_type" : "master" }
  ] // availability_zone_name值根据情况而定
  backup_storage_type   = "SATA"
  backup_storage_space  = 120
  os_type               = "ctyunos"
  cpu_type              = "Intel"
}


// 升配类型

// 升配规格 单节点->1主2备
resource "ctyun_postgresql_instance" "test" {
  cycle_type            = "on_demand"
  host_type             = "S7"
  prod_id               = "Master2Slave1222"
  storage_type          = "SATA"
  storage_space         = 120
  name                  = "pgsql-test"
  password              = var.password
  case_sensitive        = true
  instance_series       = "S"
  prod_performance_spec = "2C8G"
  vpc_id                = ctyun_vpc.vpc_test.id
  subnet_id             = ctyun_subnet.subnet_test.id
  security_group_id     = ctyun_security_group.sg_test.id
  availability_zone_info = [
    { "availability_zone_name" : "cn-gs-qyi2-1a-public-ctcloud", "availability_zone_count" : 2, "node_type" : "slave" }
  ] // availability_zone_name值根据情况而定
  backup_storage_type   = "SATA"
  backup_storage_space  = 120
  os_type               = "ctyunos"
  cpu_type              = "Intel"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cycle_type` (String) 订购周期类型，取值范围：month：按月，on_demand：按需。当此值为month时，cycle_count为必填
- `flavor_name` (String) 规格名称，形如c7.2xlarge.4，可从data.ctyun_postgresql_specs查询支持的规格，支持更新。
- `name` (String) 实例名称（长度在 4 到 64个字符，必须以字母开头，不区分大小写，可以包含字母、数字、中划线或下划线，不能包含其他特殊字符）。支持更新，但不支持更新为重名实例名称
- `prod_id` (String) 产品ID，支持更新。取值范围包括：Single1222-（单实例12.22版本）, MasterSlave1222（一主一备12.22版本）, Single1417（单实例14.17版本）, MasterSlave1417（一主一备14.17版本）, Single1320（单实例13.20版本）, MasterSlave1320（一主一备13.20版本）, ReadOnly1222（只读实例12.22版本）, ReadOnly1320（只读实例13.20版本）, ReadOnly1417（只读实例14.17版本）, Single1512（单实例15.12版本）, MasterSlave1512（一主一备15.12版本）, ReadOnly1512（只读实例15.12版本）, Master2Slave1222（一主两备12.22版本）, Master2Slave1417（一主两备14.17版本）, Master2Slave1320（一主两备13.20版本）, Master2Slave1512（一主两备15.12版本）, Single168（单实例16.8版本）, MasterSlave168（一主一备16.8版本）, Master2Slave168（一主两备16.8版本）, ReadOnly168（只读实例16.8版本）。注：扩容过程中，不支持磁盘(storage_space, backup_storage_space)、规格(flavor_name)和实例(prod_id)扩容同时进行
- `security_group_id` (String) 安全组Id
- `storage_space` (Number) 主存储空间(单位:G，范围100-32768)。支持更新，扩容过程中不支持磁盘(storage_space, backup_storage_space)、规格(flavor_name)和实例(pord_id)扩容同时进行
- `storage_type` (String) 主存储类型: SSD=超高IO, SATA=普通IO, SAS=高IO, SSD-genric=通用型SSD, FAST-SSD=极速型SSD
- `subnet_id` (String) 子网Id
- `vpc_id` (String) 虚拟私有云Id

### Optional

- `appoint_vip` (String) 指定VIP
- `auto_renew` (Boolean) 是否自动续订，默认非自动续订，当cycle_type不等于on_demand时才可填写，当cycle_count<12，到期自动续订1个月，当cycle_count>=12，到期自动续订12个月
- `availability_zone_info` (Attributes List) pgsql实例节点指定可用区字段，选填，若未填写根据实例节点数分配至各个az。示例：若创建一个一主两备的pgsql，对应的availability_zone_info为：[{availabilityZoneName:cn-huadong1-jsnj1A-public-ctcloud,availabilityZoneCount:1,nodeType:master},{availabilityZoneName:cn-huadong1-jsnj1A-public-ctcloud,availabilityZoneCount:1,nodeType:slave},{availabilityZoneName:cn-huadong1-jsnj1A-public-ctcloud,availabilityZoneCount:1,nodeType:slave}] (see [below for nested schema](#nestedatt--availability_zone_info))
- `backup_storage_space` (Number) 备份存储空间大小。支持更新，主存储空间(storage_space)若备份存储空间(backup_storage_space)同时更新，先更新backup_storage_space。
- `backup_storage_type` (String) 备份存储类型: OS=对象存储, SSD=超高IO, SATA=普通IO, SAS=高IO。注：当填写OS时，无需填写backup_storage_size
- `case_sensitive` (Boolean) 是否区分大小写: true=区分, false=不区分。默认不区分
- `cycle_count` (Number) 订购时长，该参数当且仅当在cycle_type为month时填写，支持传递1-36
- `is_mgr` (Boolean) 是否开启MRG，默认false
- `password` (String, Sensitive) 实例密码，8-32位由大写字母、小写字母、数字、特殊字符中的任意三种组成 特殊字符为!@#$%^&*()_+-=
- `project_id` (String) 企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID
- `region_id` (String) 资源池id,如果不填这默认使用provider ctyun总region_id 或者环境变量
- `running_control` (String) 控制是否暂停，启用和重启实例。支持更新，取值范围：stop, start, restart

### Read-Only

- `alive` (Number) 实例是否存活,0:存活，-1:异常
- `disk_rated` (Number) 磁盘使用率
- `id` (String) postgresql实例id
- `master_order_id` (String) 订单id
- `outer_prod_inst_id` (String) 对外的实例ID，对应PaaS平台
- `prod_db_engine` (String) 数据库实例引擎
- `prod_order_status` (Number) 订单状态，0：正常，1：冻结，2：删除，3：操作中，4：失败,2005:扩容中
- `prod_running_status` (Number) 实例状态
- `prod_type` (Number) 实例部署方式 0：单机部署,1：主备部署
- `read_port` (Number) 读端口
- `tool_type` (Number) 备份工具类型，1：pg_baseback, 2：pgbackrest, 3：s3
- `write_port` (String) 写端口

<a id="nestedatt--availability_zone_info"></a>
### Nested Schema for `availability_zone_info`

Required:

- `availability_zone_count` (Number) 资源池可用区总数
- `availability_zone_name` (String) 资源池可用区名称，可以根据data.ctyun_zones查询
- `node_type` (String) 节点类型(master/slave)