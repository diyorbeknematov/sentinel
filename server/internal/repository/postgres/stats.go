package postgres

import (
	"time"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type StatsRepo struct {
	db *sqlx.DB
}

func NewStatsRepo(db *sqlx.DB) *StatsRepo {
	return &StatsRepo{db: db}
}

func (r *StatsRepo) GetDashboardStats(filter models.StatsFilter, duration time.Duration) (models.DashboardStats, error) {
	stats := models.DashboardStats{}

	now := time.Now()
	currentStart := now.Add(-duration)
	previousStart := now.Add(-2 * duration)

	calcPercent := func(curr, prev float64) float64 {
		if prev == 0 {
			return 100
		}
		return ((curr - prev) / prev) * 100
	}

	// Current period args:  [$1=currentStart,  optionally $2=agentID]
	// Previous period args: [$1=previousStart, $2=currentStart, optionally $3=agentID]
	currAgentCond := ""
	prevAgentCond := ""

	currArgs := []interface{}{currentStart}
	prevArgs := []interface{}{previousStart, currentStart}

	if filter.AgentID != nil {
		currAgentCond = " AND agent_id = $2" // $2 = agentID
		prevAgentCond = " AND agent_id = $3" // $3 = agentID (chunki $1,$2 allaqachon band)
		currArgs = append(currArgs, *filter.AgentID)
		prevArgs = append(prevArgs, *filter.AgentID)
	}

	// 1. REQUESTS
	var currReq, prevReq float64

	err := r.db.Get(&currReq, `
		SELECT COUNT(*) FROM applogs
		WHERE recorded_at >= $1`+currAgentCond,
		currArgs...)
	if err != nil {
		return stats, err
	}

	err = r.db.Get(&prevReq, `
		SELECT COUNT(*) FROM applogs
		WHERE recorded_at BETWEEN $1 AND $2`+prevAgentCond,
		prevArgs...)
	if err != nil {
		return stats, err
	}

	stats.Requests = models.Trend{
		Value:   currReq,
		Percent: calcPercent(currReq, prevReq),
	}

	// 2. ERROR RATE
	var currErr, prevErr float64

	err = r.db.Get(&currErr, `
		SELECT COALESCE(
			100.0 * COUNT(CASE WHEN level IN ('error','fatal') THEN 1 END) / NULLIF(COUNT(*), 0),
			0
		)
		FROM applogs
		WHERE recorded_at >= $1`+currAgentCond,
		currArgs...)
	if err != nil {
		return stats, err
	}

	err = r.db.Get(&prevErr, `
		SELECT COALESCE(
			100.0 * COUNT(CASE WHEN level IN ('error','fatal') THEN 1 END) / NULLIF(COUNT(*), 0),
			0
		)
		FROM applogs
		WHERE recorded_at BETWEEN $1 AND $2`+prevAgentCond,
		prevArgs...)
	if err != nil {
		return stats, err
	}

	stats.ErrorRate = models.Trend{
		Value:   currErr,
		Percent: calcPercent(currErr, prevErr),
	}

	// 3. AUTH FAILURES
	var currAuth, prevAuth float64

	err = r.db.Get(&currAuth, `
		SELECT COUNT(*) FROM applogs
		WHERE event = 'auth_failure'
		  AND level IN ('error', 'warn')
		  AND recorded_at >= $1`+currAgentCond,
		currArgs...)
	if err != nil {
		return stats, err
	}

	err = r.db.Get(&prevAuth, `
		SELECT COUNT(*) FROM applogs
		WHERE event = 'auth_failure'
		  AND level IN ('error', 'warn')
		  AND recorded_at BETWEEN $1 AND $2`+prevAgentCond,
		prevArgs...)
	if err != nil {
		return stats, err
	}

	stats.AuthFail = models.Trend{
		Value:   currAuth,
		Percent: calcPercent(currAuth, prevAuth),
	}

	// 4. ACTIVE ALERTS
	var totalAlerts, criticalAlerts int

	alertArgs := []interface{}{}
	alertCond := ""
	if filter.AgentID != nil {
		alertCond = " WHERE agent_id = $1"
		alertArgs = append(alertArgs, *filter.AgentID)
	}

	err = r.db.QueryRow(`
		SELECT
			COUNT(*) FILTER (WHERE is_read = false),
			COUNT(*) FILTER (WHERE is_read = false AND severity = 'critical')
		FROM alerts`+alertCond,
		alertArgs...).Scan(&totalAlerts, &criticalAlerts)
	if err != nil {
		return stats, err
	}

	stats.ActiveAlerts = models.AlertStats{
		Total:    totalAlerts,
		Critical: criticalAlerts,
	}

	return stats, nil
}

func (r *StatsRepo) GetLogVolume(id uuid.UUID) ([]models.LogVolumeStats, error) {
	volumes := []models.LogVolumeStats{}

	query := `
	SELECT 
		DATE_TRUNC('hour', log_time) AS time,
		COUNT(*) FILTER (WHERE source = 'nginx') AS nginx_logs,
		COUNT(*) FILTER (WHERE source = 'app') AS app_logs,
		COUNT(*) FILTER (WHERE level = 'error' OR status >= 400) AS error_count,
		(COUNT(*) FILTER (WHERE level = 'error' OR status >= 400) > 0) AS has_error
	FROM (
		-- Nginx logs
		SELECT 
			log_time,
			agent_id,
			'nginx' AS source,
			NULL AS level,
			status
		FROM nginxlogs
		
		UNION ALL
		
		-- App logs
		SELECT 
			log_time,
			agent_id,
			'app' AS source,
			level,
			NULL AS status
		FROM applogs
	) AS combined_logs

	WHERE agent_id = $1 
		AND log_time > NOW() - INTERVAL '24 hours'
	GROUP BY DATE_TRUNC('hour', log_time)
	ORDER BY time DESC
	LIMIT 24
	`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return volumes, err
	}
	defer rows.Close()

	for rows.Next() {
		var v models.LogVolumeStats
		if err := rows.Scan(
			&v.Time,
			&v.NginxLogs,
			&v.AppLogs,
			&v.ErrorLogs,
			&v.HasError,
		); err != nil {
			return volumes, err
		}
		volumes = append(volumes, v)
	}

	return volumes, nil
}
