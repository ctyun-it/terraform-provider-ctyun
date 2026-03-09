resource "ctyun_image" "%[1]s" {
  name         = "%[2]s"
  os_distro    = "%[3]s"
  os_version   = "%[4]s"
  file_source = "%[5]s"  # 示例文件源，实际使用时需要替换为有效值
  description ="描述"
  disk_size = 40
  region_id = "bb9fdb42056f11eda1610242ac110002" // 文件只在华东1有
}