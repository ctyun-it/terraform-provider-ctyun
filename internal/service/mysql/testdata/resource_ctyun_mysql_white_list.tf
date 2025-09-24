resource "ctyun_mysql_white_list" "%[1]s" {
  prod_inst_id = "%[2]s"
  group_name = "%[3]s"
  group_white_list = [%[4]s]
}


