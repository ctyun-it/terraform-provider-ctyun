terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 可参考 index.md，在环境变量中配置 ak、sk、资源池 ID、可用区名称
provider "ctyun" {
  env = "prod"
}

# 查询函数列表（使用默认配置）
data "ctyun_functions" "all_functions" {
  # region_id 从 provider 或环境变量获取
  # page_index 默认为 1
  # page_size 默认为 10
}

# 分页查询函数列表
data "ctyun_functions" "functions_page1" {
  page_index  = 1
  page_size   = 20
  
  # region_id 可选，指定资源池
  # region_id = "200000002401"
}

# 模糊搜索函数
data "ctyun_functions" "search_functions" {
  search      = "tf-function"
  page_index  = 1
  page_size   = 10
  
  # region_id 可选
  # region_id = "200000002401"
}

# 输出查询结果
output "all_functions_info" {
  value = {
    id         = data.ctyun_functions.all_functions.id
    region_id  = data.ctyun_functions.all_functions.region_id
    page_index = data.ctyun_functions.all_functions.page_index
    page_size  = data.ctyun_functions.all_functions.page_size
    total      = length(data.ctyun_functions.all_functions.functions)
    
    functions = [
      for func in data.ctyun_functions.all_functions.functions : {
        function_id   = func.function_id
        function_name = func.function_name
        description   = func.description
        create_type   = func.create_type
        status        = func.status
        
        runtime = func.runtime != null ? {
          runtime         = func.runtime.runtime
          handler         = func.runtime.handler
          execute_timeout = func.runtime.execute_timeout
        } : null
        
        container = func.container != null ? {
          time_zone   = func.container.time_zone
          memory_size = func.container.memory_size
          cpu         = func.container.cpu
          disk_size   = func.container.disk_size
        } : null
      }
    ]
  }
}

output "search_functions_result" {
  value = {
    id         = data.ctyun_functions.search_functions.id
    region_id  = data.ctyun_functions.search_functions.region_id
    search     = data.ctyun_functions.search_functions.search
    total      = length(data.ctyun_functions.search_functions.functions)
    
    function_names = [
      for func in data.ctyun_functions.search_functions.functions : 
      func.function_name
    ]
  }
}
