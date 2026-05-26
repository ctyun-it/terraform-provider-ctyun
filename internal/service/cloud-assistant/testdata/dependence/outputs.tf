output "ecs_id1" {
  value = ctyun_ecs.ecs_test1.id
  # value = "5744fdef-0544-cec2-d06a-978bc8680779"
}
output "ecs_id2" {
  value = ctyun_ecs.ecs_test2.id
  # value = "1de0c62c-9cbf-02a6-1a13-79e2118f0b01"
}

output "command_id" {
  value = ctyun_cloud_assistant_command.command_test.id
  # value = "3050451e-22ae-4b9b-862e-c7c000104b1b"
}