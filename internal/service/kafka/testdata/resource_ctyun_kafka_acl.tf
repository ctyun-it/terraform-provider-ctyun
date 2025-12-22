resource "ctyun_kafka_acl" "%[1]s" {
  name = "%[2]s"
  instance_id = "%[3]s"
  use_new_topic = %[4]t
  rules = [{
           permission:"ALLOW",
           user_name:"%[5]s"
           operation:"READ"
    }]
}
