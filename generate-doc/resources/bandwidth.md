# ctyun_bandwidth (Resource)
**详细说明请见文档：https://www.ctyun.cn/document/10026761**



## 样例

```terraform
terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 创建一个按需，大小为5Mbit/s的带宽
resource "ctyun_bandwidth" "bandwidth_test1" {
  name       = "bandwidth-test1"
  cycle_type = "on_demand"
  bandwidth  = 5
}

# 创建一个包年，大小为10Mbit/s的带宽
resource "ctyun_bandwidth" "bandwidth_test2" {
  name        = "bandwidth-test2"
  cycle_type  = "year"
  bandwidth   = 10
  cycle_count = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `bandwidth` (Number) 共享带宽的带宽峰值（Mbit/s），必须大于等于5
- `name` (String) 共享带宽命名，单账户单资源池下，命名需唯一，长度为2-63个字符，只能由数字、字母、-组成，不能以数字、-开头，且不能以-结尾

### Optional

- `cycle_count` (Number) 订购时长, 该参数在cycle_type为month或year时才生效，当cycleType=month，支持续订1-11个月；当cycleType=year，支持续订1-3年
- `cycle_type` (String) 订购周期类型，取值范围：month：按月，year：按年、on_demand：按需。当此值为month或者year时，cycle_count为必填
- `project_id` (String) 企业项目id，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID
- `region_id` (String) 资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID

### Read-Only

- `id` (String) 共享带宽id
- `status` (String) 共享带宽状态: active：有效，expired：已过期，freezing：冻结