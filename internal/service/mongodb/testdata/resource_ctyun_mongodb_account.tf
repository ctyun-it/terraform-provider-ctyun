resource "ctyun_mongodb_account" "%[1]s" {
  instance_id = "%[2]s"
  name        = "%[3]s"
  password    = "%[6]s"
  database    = "%[4]s"
  roles  = [
    {
      db="admin"
      role="%[5]s"
    }
  ]
}

