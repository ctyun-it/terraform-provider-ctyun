resource "ctyun_eip" "%[1]s" {
  name        = "%[2]s"
  bandwidth   = %[3]s
  cycle_type = "on_demand"
  demand_billing_type = "upflowc"
}
