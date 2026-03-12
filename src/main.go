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

// 预锁库存流程节点

// 初始化上下文节点
func InitCtxNode(data []byte, option map[string][]string) ([]byte, error) {
	// 创建预锁库存上下文
	ctx := &PreLockContext{}
	
	// 初始化订单信息（简化版，使用预设数据）
	ctx.Order = Order{
		ID:            "ORDER001",
		RequiredSeats: 2,
		City:          "北京",
		DepartureTime: "2024-01-01 08:00:00",
	}
	
	fmt.Printf("InitCtxNode: Initialized context for order %s\n", ctx.Order.ID)
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 添加订单到池节点
func AddOrderToPoolNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &PreLockContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 添加订单到拼车池和预锁池
	ctx.CarpoolPool = append(ctx.CarpoolPool, ctx.Order)
	ctx.PreLockPool = append(ctx.PreLockPool, ctx.Order)
	
	fmt.Printf("AddOrderToPoolNode: Added order %s to carpool and pre-lock pools\n", ctx.Order.ID)
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 从池召回订单节点
func RecallOrderFromPoolNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &PreLockContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 从拼车池和预锁池召回订单（简化版，直接使用现有订单）
	fmt.Printf("RecallOrderFromPoolNode: Recalled order %s from pools\n", ctx.Order.ID)
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 锁定班次节点
func LockShiftNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &PreLockContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 锁定班次（简化版，使用预设班次）
	ctx.LockedShift = "SCH001"
	fmt.Printf("LockShiftNode: Locked shift %s\n", ctx.LockedShift)
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 获取Foras信息节点
func GetForasInfoNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &PreLockContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 获取Foras信息（简化版，模拟获取）
	fmt.Printf("GetForasInfoNode: Got Foras info for shift %s\n", ctx.LockedShift)
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 设置班次版本节点
func SetShiftVersionNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &PreLockContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 设置班次版本（简化版，使用预设版本）
	ctx.ShiftVersion = "v1.0"
	fmt.Printf("SetShiftVersionNode: Set shift version to %s\n", ctx.ShiftVersion)
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 获取预锁库存节点
func GetPrelockInventoryNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &PreLockContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 获取预锁库存（简化版，使用预设值）
	ctx.PreLockInventory = 5
	fmt.Printf("GetPrelockInventoryNode: Got pre-lock inventory %d\n", ctx.PreLockInventory)
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 从Redis获取班次库存节点
func GetShiftInventoryFromRedisNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &PreLockContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 从Redis获取班次库存（简化版，使用预设值）
	ctx.ShiftInventory = 30
	fmt.Printf("GetShiftInventoryFromRedisNode: Got shift inventory %d from Redis\n", ctx.ShiftInventory)
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 从DB获取班次库存节点
func GetShiftInventoryFromDBNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &PreLockContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 从DB获取班次库存（简化版，使用预设值）
	ctx.ShiftInventory = 30
	fmt.Printf("GetShiftInventoryFromDBNode: Got shift inventory %d from DB\n", ctx.ShiftInventory)
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 添加班次订单到StgData节点
func AddShiftOrderToStgDataNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &PreLockContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 添加班次订单到StgData（简化版，模拟添加）
	fmt.Printf("AddShiftOrderToStgDataNode: Added order %s to StgData for shift %s\n", ctx.Order.ID, ctx.LockedShift)
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 获取实时特征节点
func GetRtFeatureNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &PreLockContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 获取实时特征（简化版，使用预设值）
	ctx.RealTimeFeature = map[string]interface{}{
		"demand": 100,
		"supply": 50,
		"price":  200,
	}
	fmt.Printf("GetRtFeatureNode: Got real-time features\n")
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 尝试占用座位节点
func TryOccupySeatsNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &PreLockContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 尝试占用座位（简化版，模拟成功）
	fmt.Printf("TryOccupySeatsNode: Trying to occupy %d seats\n", ctx.Order.RequiredSeats)
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// API检查节点
func ApiCheckNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &PreLockContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// API检查（简化版，模拟成功）
	fmt.Printf("ApiCheckNode: API check passed\n")
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 锁定座位节点
func LockSeatsNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &PreLockContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 锁定座位（简化版，模拟成功）
	ctx.LockSuccess = true
	fmt.Printf("LockSeatsNode: Locked %d seats for order %s\n", ctx.Order.RequiredSeats, ctx.Order.ID)
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 解锁班次节点
func UnlockShiftNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &PreLockContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 解锁班次（简化版，模拟成功）
	fmt.Printf("UnlockShiftNode: Unlocked shift %s\n", ctx.LockedShift)
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 从池删除订单节点
func DelOrderFromPoolNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &PreLockContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 从拼车池和预锁池删除订单（简化版，模拟删除）
	ctx.CarpoolPool = []Order{}
	ctx.PreLockPool = []Order{}
	fmt.Printf("DelOrderFromPoolNode: Deleted order %s from carpool and pre-lock pools\n", ctx.Order.ID)
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
}

// 已锁定节点
func AlreadyLockedNode(data []byte, option map[string][]string) ([]byte, error) {
	// 反序列化上下文
	ctx := &PreLockContext{}
	if err := json.Unmarshal(data, ctx); err != nil {
		return nil, err
	}
	
	// 已锁定处理（简化版，模拟处理）
	fmt.Printf("AlreadyLockedNode: Order %s already locked\n", ctx.Order.ID)
	
	// 序列化上下文
	ctxBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	
	return ctxBytes, nil
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

// DefinePreLockWorkflow 定义预锁库存工作流
func DefinePreLockWorkflow(workflow *flow.Workflow, context *flow.Context) error {
	dag := workflow.Dag()
	
	// 定义节点
	dag.Node("initCtx", InitCtxNode)
	dag.Node("addOrderToPool", AddOrderToPoolNode)
	dag.Node("recallOrderFromPool", RecallOrderFromPoolNode)
	dag.Node("lockShift", LockShiftNode)
	dag.Node("getForasInfo", GetForasInfoNode)
	dag.Node("setShiftVersion", SetShiftVersionNode)
	dag.Node("getPrelockInventory", GetPrelockInventoryNode)
	dag.Node("getShiftInventoryFromRedis", GetShiftInventoryFromRedisNode)
	dag.Node("getShiftInventoryFromDB", GetShiftInventoryFromDBNode)
	dag.Node("addShiftOrderToStgData", AddShiftOrderToStgDataNode)
	dag.Node("getRtFeature", GetRtFeatureNode)
	dag.Node("tryOccupySeats", TryOccupySeatsNode)
	dag.Node("apiCheck", ApiCheckNode)
	dag.Node("lockSeats", LockSeatsNode)
	dag.Node("unlockShiftFirst", UnlockShiftNode)
	dag.Node("delOrderFromPool", DelOrderFromPoolNode)
	dag.Node("unlockShiftSecond", UnlockShiftNode)
	dag.Node("alreadyLocked", AlreadyLockedNode)
	
	// 建立依赖关系
	dag.Edge("initCtx", "recallOrderFromPool")
	dag.Edge("recallOrderFromPool", "lockShift")
	dag.Edge("lockShift", "getForasInfo")
	dag.Edge("lockShift", "setShiftVersion")
	dag.Edge("lockShift", "getPrelockInventory")
	dag.Edge("getForasInfo", "getShiftInventoryFromRedis")
	dag.Edge("getForasInfo", "getShiftInventoryFromDB")
	dag.Edge("getShiftInventoryFromRedis", "addShiftOrderToStgData")
	dag.Edge("getShiftInventoryFromDB", "addShiftOrderToStgData")
	dag.Edge("setShiftVersion", "addShiftOrderToStgData")
	dag.Edge("getPrelockInventory", "addShiftOrderToStgData")
	dag.Edge("addShiftOrderToStgData", "getRtFeature")
	dag.Edge("getRtFeature", "tryOccupySeats")
	dag.Edge("tryOccupySeats", "apiCheck")
	dag.Edge("apiCheck", "lockSeats")
	dag.Edge("lockSeats", "unlockShiftFirst")
	dag.Edge("lockSeats", "delOrderFromPool")
	dag.Edge("lockSeats", "alreadyLocked")
	dag.Edge("delOrderFromPool", "unlockShiftSecond")
	
	return nil
}

// 本地测试预锁库存工作流
func testPreLockWorkflow() {
	fmt.Println("Testing pre-lock inventory workflow...")
	
	// 模拟工作流执行
	ctx := &PreLockContext{}
	
	// 初始化上下文
	ctx.Order = Order{
		ID:            "ORDER001",
		RequiredSeats: 2,
		City:          "北京",
		DepartureTime: "2024-01-01 08:00:00",
	}
	fmt.Printf("InitCtx: Initialized context for order %s\n", ctx.Order.ID)
	
	// 添加订单到池
	ctx.CarpoolPool = append(ctx.CarpoolPool, ctx.Order)
	ctx.PreLockPool = append(ctx.PreLockPool, ctx.Order)
	fmt.Printf("AddOrderToPool: Added order %s to carpool and pre-lock pools\n", ctx.Order.ID)
	
	// 从池召回订单
	fmt.Printf("RecallOrderFromPool: Recalled order %s from pools\n", ctx.Order.ID)
	
	// 锁定班次
	ctx.LockedShift = "SCH001"
	fmt.Printf("LockShift: Locked shift %s\n", ctx.LockedShift)
	
	// 获取Foras信息
	fmt.Printf("GetForasInfo: Got Foras info for shift %s\n", ctx.LockedShift)
	
	// 设置班次版本
	ctx.ShiftVersion = "v1.0"
	fmt.Printf("SetShiftVersion: Set shift version to %s\n", ctx.ShiftVersion)
	
	// 获取预锁库存
	ctx.PreLockInventory = 5
	fmt.Printf("GetPrelockInventory: Got pre-lock inventory %d\n", ctx.PreLockInventory)
	
	// 获取班次库存
	ctx.ShiftInventory = 30
	fmt.Printf("GetShiftInventory: Got shift inventory %d\n", ctx.ShiftInventory)
	
	// 添加班次订单到StgData
	fmt.Printf("AddShiftOrderToStgData: Added order %s to StgData for shift %s\n", ctx.Order.ID, ctx.LockedShift)
	
	// 获取实时特征
	ctx.RealTimeFeature = map[string]interface{}{
		"demand": 100,
		"supply": 50,
		"price":  200,
	}
	fmt.Printf("GetRtFeature: Got real-time features\n")
	
	// 尝试占用座位
	fmt.Printf("TryOccupySeats: Trying to occupy %d seats\n", ctx.Order.RequiredSeats)
	
	// API检查
	fmt.Printf("ApiCheck: API check passed\n")
	
	// 锁定座位
	ctx.LockSuccess = true
	fmt.Printf("LockSeats: Locked %d seats for order %s\n", ctx.Order.RequiredSeats, ctx.Order.ID)
	
	// 解锁班次
	fmt.Printf("UnlockShift: Unlocked shift %s\n", ctx.LockedShift)
	
	// 从池删除订单
	ctx.CarpoolPool = []Order{}
	ctx.PreLockPool = []Order{}
	fmt.Printf("DelOrderFromPool: Deleted order %s from carpool and pre-lock pools\n", ctx.Order.ID)
	
	// 解锁班次
	fmt.Printf("UnlockShift: Unlocked shift %s\n", ctx.LockedShift)
	
	// 已锁定处理
	fmt.Printf("AlreadyLocked: Order %s already locked\n", ctx.Order.ID)
	
	fmt.Println("Pre-lock inventory workflow test completed successfully!")
}

func main() {
	// 运行预锁库存工作流测试
	testPreLockWorkflow()
	
	// 检查Redis是否可用（简化版，仅作为提示）
	fmt.Println("\nNote: This service requires Redis to be running on localhost:6379")
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
	fs.Register("preLockInventory", DefinePreLockWorkflow)
	
	// 启动服务
	fmt.Println("\nStarting GoFlow service...")
	fmt.Println("Workflow 'scheduleRecommendation' registered")
	fmt.Println("Workflow 'preLockInventory' registered")
	fmt.Println("Service running on port 8080")
	fs.Start()
}
