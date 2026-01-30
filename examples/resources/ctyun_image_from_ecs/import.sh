#!/bin/bash
# 导入ctyun_image_from_ecs资源的示例脚本
# 使用方法: 
# 1. 将此文件保存为 import.sh
# 2. 替换下方的 <image-id> 和 <region-id> 为实际值
# 3. 运行命令: bash import.sh

# 示例:
# terraform import ctyun_image_from_ecs.example <image-id>,<region-id>

# 注意：请将下面的占位符替换为实际值
# <image-id> - 镜像的实际ID
# <region-id> - 区域ID

terraform import ctyun_image_from_ecs.example <image-id>,<region-id>