package main

import (
	"math"
	"time"
)

// CalculateDistance 计算两个坐标之间的距离（单位：公里）
func CalculateDistance(coord1, coord2 Coordinate) float64 {
	const R = 6371.0 // 地球半径
	lat1 := coord1.Latitude * math.Pi / 180.0
	lon1 := coord1.Longitude * math.Pi / 180.0
	lat2 := coord2.Latitude * math.Pi / 180.0
	lon2 := coord2.Longitude * math.Pi / 180.0
	
	dlat := lat2 - lat1
	dlon := lon2 - lon1
	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := R * c
	
	return distance
}

// ReadSchedules 读取班次数据（简化版，直接返回预设数据）
func ReadSchedules() []Schedule {
	schedules := []Schedule{}
	
	// 创建预设班次数据
	sch1 := Schedule{
		ID:             "SCH001",
		City:           "北京",
		TotalSeats:     50,
		RemainingSeats: 30,
		DepartureTime:  "2024-01-01 08:00:00",
		Stations: []Station{
			{
				Name:       "北京站",
				Coord:      Coordinate{Latitude: 39.9042, Longitude: 116.4074},
				ArriveTime: "2024-01-01 07:45:00",
				DepartTime: "2024-01-01 08:00:00",
			},
			{
				Name:       "北京南站",
				Coord:      Coordinate{Latitude: 39.8651, Longitude: 116.3786},
				ArriveTime: "2024-01-01 08:30:00",
				DepartTime: "2024-01-01 08:40:00",
			},
			{
				Name:       "天津站",
				Coord:      Coordinate{Latitude: 39.0842, Longitude: 117.2009},
				ArriveTime: "2024-01-01 09:30:00",
				DepartTime: "2024-01-01 09:40:00",
			},
		},
	}
	schedules = append(schedules, sch1)
	
	sch2 := Schedule{
		ID:             "SCH002",
		City:           "北京",
		TotalSeats:     40,
		RemainingSeats: 5,
		DepartureTime:  "2024-01-01 09:00:00",
		Stations: []Station{
			{
				Name:       "北京站",
				Coord:      Coordinate{Latitude: 39.9042, Longitude: 116.4074},
				ArriveTime: "2024-01-01 08:45:00",
				DepartTime: "2024-01-01 09:00:00",
			},
			{
				Name:       "北京西站",
				Coord:      Coordinate{Latitude: 39.8586, Longitude: 116.3286},
				ArriveTime: "2024-01-01 09:20:00",
				DepartTime: "2024-01-01 09:30:00",
			},
			{
				Name:       "石家庄站",
				Coord:      Coordinate{Latitude: 38.0457, Longitude: 114.4995},
				ArriveTime: "2024-01-01 11:00:00",
				DepartTime: "2024-01-01 11:10:00",
			},
		},
	}
	schedules = append(schedules, sch2)
	
	sch3 := Schedule{
		ID:             "SCH003",
		City:           "北京",
		TotalSeats:     60,
		RemainingSeats: 40,
		DepartureTime:  "2024-01-01 10:00:00",
		Stations: []Station{
			{
				Name:       "北京南站",
				Coord:      Coordinate{Latitude: 39.8651, Longitude: 116.3786},
				ArriveTime: "2024-01-01 09:45:00",
				DepartTime: "2024-01-01 10:00:00",
			},
			{
				Name:       "天津站",
				Coord:      Coordinate{Latitude: 39.0842, Longitude: 117.2009},
				ArriveTime: "2024-01-01 10:50:00",
				DepartTime: "2024-01-01 11:00:00",
			},
			{
				Name:       "济南站",
				Coord:      Coordinate{Latitude: 36.6683, Longitude: 117.0207},
				ArriveTime: "2024-01-01 13:00:00",
				DepartTime: "2024-01-01 13:10:00",
			},
		},
	}
	schedules = append(schedules, sch3)
	
	sch4 := Schedule{
		ID:             "SCH004",
		City:           "上海",
		TotalSeats:     55,
		RemainingSeats: 25,
		DepartureTime:  "2024-01-01 08:30:00",
		Stations: []Station{
			{
				Name:       "上海站",
				Coord:      Coordinate{Latitude: 31.2304, Longitude: 121.4737},
				ArriveTime: "2024-01-01 08:15:00",
				DepartTime: "2024-01-01 08:30:00",
			},
			{
				Name:       "上海南站",
				Coord:      Coordinate{Latitude: 31.1629, Longitude: 121.4366},
				ArriveTime: "2024-01-01 08:50:00",
				DepartTime: "2024-01-01 09:00:00",
			},
			{
				Name:       "杭州站",
				Coord:      Coordinate{Latitude: 30.2741, Longitude: 120.1551},
				ArriveTime: "2024-01-01 10:30:00",
				DepartTime: "2024-01-01 10:40:00",
			},
		},
	}
	schedules = append(schedules, sch4)
	
	sch5 := Schedule{
		ID:             "SCH005",
		City:           "上海",
		TotalSeats:     45,
		RemainingSeats: 10,
		DepartureTime:  "2024-01-01 09:30:00",
		Stations: []Station{
			{
				Name:       "上海站",
				Coord:      Coordinate{Latitude: 31.2304, Longitude: 121.4737},
				ArriveTime: "2024-01-01 09:15:00",
				DepartTime: "2024-01-01 09:30:00",
			},
			{
				Name:       "南京站",
				Coord:      Coordinate{Latitude: 32.0603, Longitude: 118.7969},
				ArriveTime: "2024-01-01 11:30:00",
				DepartTime: "2024-01-01 11:40:00",
			},
		},
	}
	schedules = append(schedules, sch5)
	
	sch6 := Schedule{
		ID:             "SCH006",
		City:           "广州",
		TotalSeats:     50,
		RemainingSeats: 35,
		DepartureTime:  "2024-01-01 10:00:00",
		Stations: []Station{
			{
				Name:       "广州站",
				Coord:      Coordinate{Latitude: 23.1291, Longitude: 113.2644},
				ArriveTime: "2024-01-01 09:45:00",
				DepartTime: "2024-01-01 10:00:00",
			},
			{
				Name:       "深圳站",
				Coord:      Coordinate{Latitude: 22.5431, Longitude: 114.0579},
				ArriveTime: "2024-01-01 11:30:00",
				DepartTime: "2024-01-01 11:40:00",
			},
		},
	}
	schedules = append(schedules, sch6)
	
	return schedules
}

// ReadRequests 读取用户请求数据（简化版，直接返回预设数据）
func ReadRequests() []UserRequest {
	requests := []UserRequest{}
	
	// 创建预设请求数据
	req1 := UserRequest{
		ID:            "REQ001",
		RequiredSeats: 2,
		DepartureTime: "2024-01-01 08:00:00",
		City:          "北京",
		StartCoord:    Coordinate{Latitude: 39.9042, Longitude: 116.4074},
		EndCoord:      Coordinate{Latitude: 39.0842, Longitude: 117.2009},
		PageSize:      10,
		PageNum:       1,
	}
	requests = append(requests, req1)
	
	req2 := UserRequest{
		ID:            "REQ002",
		RequiredSeats: 5,
		DepartureTime: "2024-01-01 09:00:00",
		City:          "北京",
		StartCoord:    Coordinate{Latitude: 39.8586, Longitude: 116.3286},
		EndCoord:      Coordinate{Latitude: 38.0457, Longitude: 114.4995},
		PageSize:      10,
		PageNum:       1,
	}
	requests = append(requests, req2)
	
	req3 := UserRequest{
		ID:            "REQ003",
		RequiredSeats: 3,
		DepartureTime: "2024-01-01 10:00:00",
		City:          "上海",
		StartCoord:    Coordinate{Latitude: 31.2304, Longitude: 121.4737},
		EndCoord:      Coordinate{Latitude: 30.2741, Longitude: 120.1551},
		PageSize:      10,
		PageNum:       1,
	}
	requests = append(requests, req3)
	
	req4 := UserRequest{
		ID:            "REQ004",
		RequiredSeats: 10,
		DepartureTime: "2024-01-01 08:30:00",
		City:          "上海",
		StartCoord:    Coordinate{Latitude: 31.1629, Longitude: 121.4366},
		EndCoord:      Coordinate{Latitude: 30.2741, Longitude: 120.1551},
		PageSize:      10,
		PageNum:       1,
	}
	requests = append(requests, req4)
	
	req5 := UserRequest{
		ID:            "REQ005",
		RequiredSeats: 2,
		DepartureTime: "2024-01-01 10:00:00",
		City:          "广州",
		StartCoord:    Coordinate{Latitude: 23.1291, Longitude: 113.2644},
		EndCoord:      Coordinate{Latitude: 22.5431, Longitude: 114.0579},
		PageSize:      10,
		PageNum:       1,
	}
	requests = append(requests, req5)
	
	return requests
}

// FormatTime 格式化时间字符串为时间戳
func FormatTime(timeStr string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return time.Time{}
	}
	return t
}
