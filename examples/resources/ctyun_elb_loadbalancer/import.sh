#!/bin/bash
# 导入ctyun_elb_loadbalancer资源的示例脚本
# 使用方法: 
# 1. 将此文件保存为 import.sh
# 2. 替换下方的 <loadbalancer-id> 和 <region-id> 为实际值
# 3. 运行命令: bash import.sh

# 示例:
# terraform import ctyun_elb_loadbalancer.example <loadbalancer-id>,<region-id>

terraform import ctyun_elb_loadbalancer.example <loadbalancer-id>,<region-id>