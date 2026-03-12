package main

import (
	"encoding/json"
	"fmt"
	"sort"

	flow "github.com/s8sg/goflow/flow/v1"
	goflow "github.com/s8sg/goflow/v1"
)

// 数据加载节点
func DataLoadNode(data []byte, option map[string][]string) ([]byte, error) {
	// 创建上下文
	ctx := &ScheduleContext{}
	
	// 读取班次数据
	ctx.AllSchedules = ReadSchedules()
	fmt.Printf("DataLoadNode: Loaded %d schedules\n", len(ctx.AllSchedules))
	
	// 读取请求数据（这里简化处理，只取第一个请求）
	requests := ReadRequests()
	if len(requests) > 0 {
		ctx.Request = requests[0]
	}
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 召回节点
func RecallNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &ScheduleContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 根据城市召回班次
	city := ctx.Request.City
	for _, schedule := range ctx.AllSchedules {
		if schedule.City == city {
			ctx.RecalledSchedules = append(ctx.RecalledSchedules, schedule)
		}
	}
	
	fmt.Printf("RecallNode: Recalled %d schedules for city %s\n", len(ctx.RecalledSchedules), city)
	ctx.HasSchedules = len(ctx.RecalledSchedules) > 0
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 座位过滤节点
func SeatFilterNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &ScheduleContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 过滤座位数
	requiredSeats := ctx.Request.RequiredSeats
	for _, schedule := range ctx.RecalledSchedules {
		if schedule.RemainingSeats >= requiredSeats {
			ctx.FilteredSchedules = append(ctx.FilteredSchedules, schedule)
		}
	}
	
	fmt.Printf("SeatFilterNode: Filtered to %d schedules with at least %d seats\n", len(ctx.FilteredSchedules), requiredSeats)
	ctx.HasSchedules = len(ctx.FilteredSchedules) > 0
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 时间过滤节点
func TimeFilterNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &ScheduleContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 过滤时间
	requestTime := FormatTime(ctx.Request.DepartureTime)
	var tempSchedules []Schedule
	
	for _, schedule := range ctx.FilteredSchedules {
		scheduleTime := FormatTime(schedule.DepartureTime)
		// 只保留发车时间大于等于请求时间的班次
		if !scheduleTime.Before(requestTime) {
			tempSchedules = append(tempSchedules, schedule)
		}
	}
	
	ctx.FilteredSchedules = tempSchedules
	fmt.Printf("TimeFilterNode: Filtered to %d schedules with departure time >= %s\n", len(ctx.FilteredSchedules), ctx.Request.DepartureTime)
	ctx.HasSchedules = len(ctx.FilteredSchedules) > 0
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 距离过滤节点
func DistanceFilterNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &ScheduleContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 计算距离并过滤
	startCoord := ctx.Request.StartCoord
	endCoord := ctx.Request.EndCoord
	var tempSchedules []Schedule
	
	for _, schedule := range ctx.FilteredSchedules {
		// 计算与起始站点的距离
		minStartDistance := 100000.0
		minEndDistance := 100000.0
		
		for _, station := range schedule.Stations {
			startDistance := CalculateDistance(startCoord, station.Coord)
			endDistance := CalculateDistance(endCoord, station.Coord)
			if startDistance < minStartDistance {
				minStartDistance = startDistance
			}
			if endDistance < minEndDistance {
				minEndDistance = endDistance
			}
		}
		
		// 计算总距离
		schedule.Distance = minStartDistance + minEndDistance
		
		// 过滤距离阈值（这里设置为50公里）
		if minStartDistance <= 50 && minEndDistance <= 50 {
			tempSchedules = append(tempSchedules, schedule)
		}
	}
	
	ctx.FilteredSchedules = tempSchedules
	fmt.Printf("DistanceFilterNode: Filtered to %d schedules within distance threshold\n", len(ctx.FilteredSchedules))
	ctx.HasSchedules = len(ctx.FilteredSchedules) > 0
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 排序节点
func SortNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &ScheduleContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 根据距离排序
	ctx.SortedSchedules = make([]Schedule, len(ctx.FilteredSchedules))
	copy(ctx.SortedSchedules, ctx.FilteredSchedules)
	sort.Slice(ctx.SortedSchedules, func(i, j int) bool {
		return ctx.SortedSchedules[i].Distance < ctx.SortedSchedules[j].Distance
	})
	
	fmt.Printf("SortNode: Sorted %d schedules by distance\n", len(ctx.SortedSchedules))
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 分页和结果处理节点
func PaginateAndResultNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &ScheduleContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 分页处理
	pageSize := ctx.Request.PageSize
	pageNum := ctx.Request.PageNum
	startIdx := (pageNum - 1) * pageSize
	endIdx := startIdx + pageSize
	
	if startIdx < len(ctx.SortedSchedules) {
		if endIdx > len(ctx.SortedSchedules) {
			endIdx = len(ctx.SortedSchedules)
		}
		ctx.PaginatedSchedules = ctx.SortedSchedules[startIdx:endIdx]
	}
	
	fmt.Printf("PaginateNode: Paginated to %d schedules (page %d, size %d)\n", len(ctx.PaginatedSchedules), pageNum, pageSize)
	
	// 处理结果
	if len(ctx.PaginatedSchedules) > 0 {
		// 有结果
		ctx.Response.Success = true
		ctx.Response.Message = "Success"
		ctx.Response.Schedules = ctx.PaginatedSchedules
		ctx.Response.Total = len(ctx.SortedSchedules)
		ctx.Response.PageSize = ctx.Request.PageSize
		ctx.Response.PageNum = ctx.Request.PageNum
		
		// 输出结果
		fmt.Printf("ResultNode: Processed result for request %s\n", ctx.Request.ID)
		fmt.Printf("Total schedules: %d\n", ctx.Response.Total)
		fmt.Printf("Returned schedules: %d\n", len(ctx.Response.Schedules))
		for _, schedule := range ctx.Response.Schedules {
			fmt.Printf("  Schedule: %s, Distance: %.2fkm, Departure: %s\n", schedule.ID, schedule.Distance, schedule.DepartureTime)
		}
	} else {
		// 无结果
		ctx.Response.Success = false
		ctx.Response.Message = "No schedules found"
		ctx.Response.Total = 0
		ctx.Response.PageSize = ctx.Request.PageSize
		ctx.Response.PageNum = ctx.Request.PageNum
		
		fmt.Printf("NoResultNode: No schedules found for request %s\n", ctx.Request.ID)
	}
	
	// 序列化响应
	responseBytes, err := json.Marshal(ctx.Response)
	if err != nil {
		return nil, err
	}
	
	return responseBytes, nil
}

// DefineWorkflow 定义工作流
func DefineWorkflow(workflow *flow.Workflow, context *flow.Context) error {
	dag := workflow.Dag()
	
	// 定义节点
	dag.Node("dataLoad", DataLoadNode)
	dag.Node("recall", RecallNode)
	dag.Node("seatFilter", SeatFilterNode)
	dag.Node("timeFilter", TimeFilterNode)
	dag.Node("distanceFilter", DistanceFilterNode)
	dag.Node("sort", SortNode)
	dag.Node("paginateAndResult", PaginateAndResultNode)
	
	// 建立依赖关系
	dag.Edge("dataLoad", "recall")
	dag.Edge("recall", "seatFilter")
	dag.Edge("seatFilter", "timeFilter")
	dag.Edge("timeFilter", "distanceFilter")
	dag.Edge("distanceFilter", "sort")
	dag.Edge("sort", "paginateAndResult")
	
	return nil
}

func main() {
	// 检查Redis是否可用（简化版，仅作为提示）
	fmt.Println("Note: This service requires Redis to be running on localhost:6379")
	fmt.Println("If Redis is not available, the service may not start properly")
	fmt.Println("Please start Redis before running this service")
	
	fs := &goflow.FlowService{
		Port:              8080,
		RedisURL:          "localhost:6379",
		OpenTraceUrl:      "localhost:5775",
		WorkerConcurrency: 5,
		EnableMonitoring:  true,
	}
	
	// 注册工作流
	fs.Register("scheduleRecommendation", DefineWorkflow)
	
	// 启动服务
	fmt.Println("Starting GoFlow service...")
	fmt.Println("Workflow 'scheduleRecommendation' registered")
	fmt.Println("Service running on port 8080")
	fs.Start()
}
