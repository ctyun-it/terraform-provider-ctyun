# ctyun_bandwidth_association_eip (Resource)
**详细说明请见文档：https://www.ctyun.cn/document/10026761/10030030**



## 样例

```terraform
terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

resource "ctyun_bandwidth_association_eip" "bandwidth_association_eip_test" {
  bandwidth_id = "bandwidth-at2yy664m5"
  eip_id       = "eip-p9qvl63yt6"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `bandwidth_id` (String) 共享带宽id
- `eip_id` (String) 弹性ip的id

### Optional

- `project_id` (String) 企业项目id，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID
- `region_id` (String) 资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID