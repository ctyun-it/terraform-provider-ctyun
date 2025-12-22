data "ctyun_vips" "test" {
  region_id = "100054c0416811e9a6690242ac110002"
  project_id = "0"
  page_no = 1
  page_size = 10
}

output "vips" {
  value = data.ctyun_vips.test.vips
}