package models

import "time"

type Trend struct {
	Value   float64 `json:"value"`
	Percent float64 `json:"percent"`
}

type AlertStats struct {
	Total    int `json:"total"`
	Critical int `json:"critical"`
}

type DashboardStats struct {
	Requests  Trend `json:"requests"`
	ErrorRate Trend `json:"errorRate"`
	AuthFail  Trend `json:"authFail"`

	ActiveAlerts AlertStats `json:"activeAlerts"`
}

type StatsFilter struct {
	AgentID *string `form:"agent_id"`
	Period  string  `form:"period"` // "1h", "24h", "7d" - default "1h"
}

type LogVolumeStats struct {
	Time      time.Time `json:"time"`
	NginxLogs int       `json:"nginx_logs"`
	AppLogs   int       `json:"app_logs"`
	ErrorLogs int       `json:"error_count"`
	HasError  bool      `json:"has_error"`
}
