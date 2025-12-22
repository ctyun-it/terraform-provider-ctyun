output "cluster_name" {
  value = data.ctyun_hpfs_clusters.test.clusters[0].name
}