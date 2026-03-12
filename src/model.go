package main

// Coordinate 坐标结构
type Coordinate struct {
	Latitude  float64 `json:"latitude"`  // 纬度
	Longitude float64 `json:"longitude"` // 经度
}

// Station 站点信息
type Station struct {
	Name        string     `json:"name"`        // 站点名称
	Coord       Coordinate `json:"coord"`       // 站点坐标
	ArriveTime  string     `json:"arrive_time"`  // 预计到达时间
	DepartTime  string     `json:"depart_time"`  // 预计出发时间
}

// Schedule 班次信息
type Schedule struct {
	ID              string     `json:"id"`              // 班次ID
	City            string     `json:"city"`            // 城市
	TotalSeats      int        `json:"total_seats"`      // 库存总数
	RemainingSeats  int        `json:"remaining_seats"`  // 剩余库存数
	DepartureTime   string     `json:"departure_time"`   // 发车时间
	Stations        []Station  `json:"stations"`        // 站点列表
	Distance        float64    `json:"distance"`        // 与用户的距离（计算后填充）
}

// UserRequest 用户请求信息
type UserRequest struct {
	ID              string     `json:"id"`              // 请求ID
	RequiredSeats   int        `json:"required_seats"`   // 座位数
	DepartureTime   string     `json:"departure_time"`   // 乘车时间
	City            string     `json:"city"`            // 城市
	StartCoord      Coordinate `json:"start_coord"`      // 上车坐标
	EndCoord        Coordinate `json:"end_coord"`        // 下车坐标
	PageSize        int        `json:"page_size"`        // 分页大小
	PageNum         int        `json:"page_num"`         // 页码
}

// Response 响应结果
type Response struct {
	Success     bool       `json:"success"`     // 是否成功
	Message     string     `json:"message"`     // 消息
	Schedules   []Schedule `json:"schedules"`   // 推荐班次列表
	Total       int        `json:"total"`       // 总数
	PageSize    int        `json:"page_size"`    // 分页大小
	PageNum     int        `json:"page_num"`     // 页码
}

// ScheduleContext 上下文结构
type ScheduleContext struct {
	AllSchedules     []Schedule
	RecalledSchedules []Schedule
	FilteredSchedules []Schedule
	SortedSchedules   []Schedule
	PaginatedSchedules []Schedule
	Request          UserRequest
	Response         Response
	HasSchedules     bool
}
