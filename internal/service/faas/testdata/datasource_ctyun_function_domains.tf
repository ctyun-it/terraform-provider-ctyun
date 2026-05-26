
data "ctyun_function_domains" "%[1]s" {
  page_index = 1
  page_size  = 10
  search_key = "%[2]s"
}

