data "ctyun_iam_policies" "%[1]s" {
  policy_id = %[2]s
}

data "ctyun_iam_policies" "test" {
  page_no = 1
  page_size = 100
}

data "ctyun_services" "test" {

}

data "ctyun_iam_authorities" "test" {
  service_id = data.ctyun_services.test.services[0].id
}

