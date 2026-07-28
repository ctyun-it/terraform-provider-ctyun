# testdata/datasource_ctyun_search_instances.tf
data "ctyun_search_instances" "%[1]s" {
  region_id  = "%[2]s"
  page_no    = 1
  page_size  = 10
}
