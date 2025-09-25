data "ctyun_ccse_plugin_market" "%[1]s" {
  page_no = 1
  page_size = 10
}

data "ctyun_ccse_plugin_market" "%[2]s" {
  chart_name = data.ctyun_ccse_plugin_market.%[1]s.records[0].chart_name
}

data "ctyun_ccse_plugin_market" "%[3]s" {
  chart_name = data.ctyun_ccse_plugin_market.%[1]s.records[0].chart_name
  chart_version = data.ctyun_ccse_plugin_market.%[2]s.versions[0].chart_version
  values_type = "%[4]s"
}