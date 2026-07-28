resource "ctyun_ccse_namespace" "%[1]s" {
  cluster_id = "%[2]s"
  values_yaml = <<EOF
%[3]sEOF
}
