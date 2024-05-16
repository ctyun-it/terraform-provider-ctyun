terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

data "ctyun_iam_authorities" "iam_authorities_test" {
  service_id = 108
}

output "ctyun_policies_test" {
  value = data.ctyun_iam_authorities.iam_authorities_test.authorities
}
