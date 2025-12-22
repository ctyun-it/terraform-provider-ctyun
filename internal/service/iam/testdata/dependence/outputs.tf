output "user_id" {
  value = ctyun_iam_user.test.id
}

output "group_id" {
  value = data.ctyun_iam_user_groups.test.groups[0].id
}

output "group_id2" {
  value = data.ctyun_iam_user_groups.test.groups[1].id
}

output "auth_code" {
  value = data.ctyun_iam_authorities.test.authorities[0].code
}

output "auth_code2" {
  value = data.ctyun_iam_authorities.test.authorities[1].code
}

output "policy_id" {
  value = ctyun_iam_policy.test.id
}

output "policy_id2" {
  value = ctyun_iam_policy.test2.id
}