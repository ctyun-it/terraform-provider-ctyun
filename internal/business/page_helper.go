package business

var NoNextRollingPage = RollingPageResult{
	TotalCount:       0,
	CurrentPageCount: 0,
	ExecuteContinue:  false,
}

type RollingPageResult struct {
	TotalCount       int  // 总数
	CurrentPageCount int  // 当前页面获取到的数量
	ExecuteContinue  bool // 是否继续执行
}

type RollingPageParam struct {
	CurrentPage int
}

// RollingExecution 实际分页执行的动作
type RollingExecution func(RollingPageParam) RollingPageResult

// RollingPage 执行分页动作
func RollingPage(r RollingExecution) {
	currentPage := 1
	totalCount := 0
	accumulativeCount := 0
	param := RollingPageParam{}

	for accumulativeCount < totalCount || currentPage == 1 {
		param.CurrentPage = currentPage
		result := r(param)
		if !result.ExecuteContinue || result.CurrentPageCount == 0 || result.TotalCount == 0 {
			break
		}
		currentPage++
		totalCount = result.TotalCount
		accumulativeCount += result.CurrentPageCount
	}
}
