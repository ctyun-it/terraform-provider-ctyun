data "ctyun_ccse_template_market" "%[1]s" {
  page_no = 1
  page_size = 10
}

data "ctyun_ccse_template_market" "%[2]s" {
  tpl_name = data.ctyun_ccse_template_market.%[1]s.records[0].tpl_name
}

data "ctyun_ccse_template_market" "%[3]s" {
  tpl_name = data.ctyun_ccse_template_market.%[1]s.records[0].tpl_name
  tpl_version = data.ctyun_ccse_template_market.%[2]s.versions[0].tpl_version
  values_type = "%[4]s"
}