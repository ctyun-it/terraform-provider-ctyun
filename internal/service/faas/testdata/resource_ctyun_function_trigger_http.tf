resource "ctyun_function_trigger" "%[1]s" {
  function_name = "%[2]s"
  trigger_name  = "%[3]s"
  trigger_type  = "http"

  event_data = jsonencode({
    methods      = ["GET"]
    authType     = "anonymous"
    authConfig   = {
      tokenConfig = []
      claimTrans  = []
      matchMode = {
        mode = "All"
        path = [""]
      }
    }
    internetUrl  = "func-hycjxokfed-rmjdenonei-2c05d.www.zhangzhh13.xyz"
    protocol     = "HTTP"
    certConfig   = {
      certName    = ""
      certificate = ""
      privateKey  = ""
    }
  })

  enable = true
}
