data "ctyun_iam_user_groups" "test" {
  page_no = 1
  page_size = 10
}

data "ctyun_services" "test" {

}

data "ctyun_iam_authorities" "test" {
  service_id = data.ctyun_services.test.services[0].id
}

data "ctyun_iam_policies" "test" {
  page_no = 1
  page_size = 100
}

resource "ctyun_iam_user" "test" {
  email          = "0987qwer@ctyun.cn"
  phone          = "18809871234"
  name           = "tf-user"
  password       = var.password
  description    = "terraform测试"
  user_group_ids = [
    data.ctyun_iam_user_groups.test.groups[0].id
  ]
}

resource "ctyun_iam_policy" "test"  {
  name        = "tf-policy"
  description = "terraform测试"
  range       = "region"
  content     = {
    version   = "1.1"
    statement = [
      {
        effect   = "deny"
        action   = [data.ctyun_iam_authorities.test.authorities[0].code]
        resource = ["*"]
      }
    ]
  }
}

resource "ctyun_iam_policy" "test2"  {
  name        = "tf-policy2"
  description = "terraform测试"
  range       = "global"
  content     = {
    version   = "1.1"
    statement = [
      {
        effect   = "deny"
        action   = [data.ctyun_iam_authorities.test.authorities[1].code]
        resource = ["*"]
      }
    ]
  }
}

variable "password" {
  type        = string
  sensitive = true
}
