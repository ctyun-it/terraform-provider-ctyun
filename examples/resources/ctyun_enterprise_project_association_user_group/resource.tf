terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

resource "ctyun_enterprise_project_association_user_group" "enterprise_project_association_user_group_test" {
  enterprise_project_id = "38f35a0eeaa549be8b456b8d1c251d11"
  user_group_id         = "6edf8a6a9b09442295206feef0d39132"
  policy_ids = [
    "cf0d8c9024a448aa94e1a6d3ef38bb15",
    "9d096fce7c764908a5d94c55a1bba7e6"
  ]
}