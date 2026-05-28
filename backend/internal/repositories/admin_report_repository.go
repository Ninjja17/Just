package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type AdminReportRepository struct {
	db *sql.DB
}

func NewAdminReportRepository(db *sql.DB) *AdminReportRepository {
	return &AdminReportRepository{db: db}
}

type UserListRow struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	IsVerified    bool      `json:"is_verified"`
	CreatedAt     string    `json:"created_at"`
	SkillCount    int       `json:"skill_count"`
	SessionCount  int       `json:"session_count"`
	TotalMinutes  int       `json:"total_minutes"`
}

func (r *AdminReportRepository) ListUsers(ctx context.Context, search string, limit, offset int) ([]UserListRow, int, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	whereSQL := ""
	args := []interface{}{}
	if s := strings.TrimSpace(search); s != "" {
		whereSQL = "WHERE u.email ILIKE $1 OR u.name ILIKE $1"
		args = append(args, "%"+s+"%")
	}

	// Total count
	var total int
	countQ := fmt.Sprintf("SELECT COUNT(*) FROM users u %s", whereSQL)
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limitIdx := len(args) + 1
	offsetIdx := len(args) + 2
	args = append(args, limit, offset)

	q := fmt.Sprintf(`
		SELECT u.id, u.email, u.name, u.is_verified, u.created_at,
		       COALESCE((SELECT COUNT(*) FROM skills s WHERE s.user_id = u.id), 0)   AS skill_count,
		       COALESCE((SELECT COUNT(*) FROM sessions se WHERE se.user_id = u.id), 0) AS session_count,
		       COALESCE((SELECT SUM(se.duration_minutes) FROM sessions se WHERE se.user_id = u.id), 0) AS total_minutes
		FROM users u
		%s
		ORDER BY u.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereSQL, limitIdx, offsetIdx)

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := []UserListRow{}
	for rows.Next() {
		var u UserListRow
		var created sql.NullTime
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.IsVerified, &created, &u.SkillCount, &u.SessionCount, &u.TotalMinutes); err != nil {
			return nil, 0, err
		}
		if created.Valid {
			u.CreatedAt = created.Time.UTC().Format("2006-01-02T15:04:05Z")
		}
		out = append(out, u)
	}
	return out, total, rows.Err()
}

type UserDetail struct {
	User         map[string]interface{}   `json:"user"`
	Skills       []map[string]interface{} `json:"skills"`
	RecentSessions []map[string]interface{} `json:"recent_sessions"`
	TotalMinutes int                      `json:"total_minutes"`
	SessionCount int                      `json:"session_count"`
}

func (r *AdminReportRepository) GetUserDetail(ctx context.Context, id uuid.UUID) (*UserDetail, error) {
	u := map[string]interface{}{}
	var email, name string
	var verified bool
	var created sql.NullTime
	err := r.db.QueryRowContext(ctx,
		`SELECT email, name, is_verified, created_at FROM users WHERE id = $1`, id,
	).Scan(&email, &name, &verified, &created)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	u["id"] = id
	u["email"] = email
	u["name"] = name
	u["is_verified"] = verified
	if created.Valid {
		u["created_at"] = created.Time
	}

	skills := []map[string]interface{}{}
	rows, err := r.db.QueryContext(ctx,
		`SELECT s.id, s.name, s.category, s.target_hours,
		        COALESCE(SUM(se.duration_minutes), 0) AS total_minutes
		 FROM skills s
		 LEFT JOIN sessions se ON se.skill_id = s.id
		 WHERE s.user_id = $1
		 GROUP BY s.id
		 ORDER BY s.name`, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var sid uuid.UUID
		var sname, scat string
		var target, mins int
		if err := rows.Scan(&sid, &sname, &scat, &target, &mins); err != nil {
			rows.Close()
			return nil, err
		}
		skills = append(skills, map[string]interface{}{
			"id":            sid,
			"name":          sname,
			"category":      scat,
			"target_hours":  target,
			"total_minutes": mins,
		})
	}
	rows.Close()

	recent := []map[string]interface{}{}
	rows, err = r.db.QueryContext(ctx,
		`SELECT se.id, se.skill_id, sk.name, se.start_time, se.end_time, se.duration_minutes, se.session_type
		 FROM sessions se
		 JOIN skills sk ON sk.id = se.skill_id
		 WHERE se.user_id = $1
		 ORDER BY se.start_time DESC
		 LIMIT 25`, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var sesID, skID uuid.UUID
		var skName, sType string
		var start sql.NullTime
		var end sql.NullTime
		var dur int
		if err := rows.Scan(&sesID, &skID, &skName, &start, &end, &dur, &sType); err != nil {
			rows.Close()
			return nil, err
		}
		row := map[string]interface{}{
			"id":               sesID,
			"skill_id":         skID,
			"skill_name":       skName,
			"duration_minutes": dur,
			"session_type":     sType,
		}
		if start.Valid {
			row["start_time"] = start.Time
		}
		if end.Valid {
			row["end_time"] = end.Time
		}
		recent = append(recent, row)
	}
	rows.Close()

	var totalSessions, totalMinutes int
	_ = r.db.QueryRowContext(ctx,
		`SELECT COUNT(*), COALESCE(SUM(duration_minutes),0) FROM sessions WHERE user_id = $1`, id,
	).Scan(&totalSessions, &totalMinutes)

	return &UserDetail{
		User:           u,
		Skills:         skills,
		RecentSessions: recent,
		TotalMinutes:   totalMinutes,
		SessionCount:   totalSessions,
	}, nil
}

type GlobalStats struct {
	TotalUsers       int     `json:"total_users"`
	VerifiedUsers    int     `json:"verified_users"`
	TotalSkills      int     `json:"total_skills"`
	TotalSessions    int     `json:"total_sessions"`
	TotalMinutes     int     `json:"total_minutes"`
	NewUsers7d       int     `json:"new_users_7d"`
	ActiveUsers7d    int     `json:"active_users_7d"`
	SessionsPerDay   []DayPoint `json:"sessions_per_day"`
	NewUsersPerDay   []DayPoint `json:"new_users_per_day"`
	TopSkills        []SkillCount `json:"top_skills"`
}

type DayPoint struct {
	Date  string `json:"date"`
	Value int    `json:"value"`
}

type SkillCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (r *AdminReportRepository) GlobalStats(ctx context.Context) (*GlobalStats, error) {
	g := &GlobalStats{}

	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&g.TotalUsers)
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE is_verified = true`).Scan(&g.VerifiedUsers)
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM skills`).Scan(&g.TotalSkills)
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*), COALESCE(SUM(duration_minutes),0) FROM sessions`).Scan(&g.TotalSessions, &g.TotalMinutes)
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE created_at > NOW() - INTERVAL '7 days'`).Scan(&g.NewUsers7d)
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(DISTINCT user_id) FROM sessions WHERE start_time > NOW() - INTERVAL '7 days'`).Scan(&g.ActiveUsers7d)

	// Sessions per day, last 30 days
	rows, err := r.db.QueryContext(ctx, `
		SELECT to_char(d::date, 'YYYY-MM-DD') AS date,
		       COALESCE(COUNT(s.id), 0) AS cnt
		FROM generate_series(NOW()::date - INTERVAL '29 days', NOW()::date, INTERVAL '1 day') d
		LEFT JOIN sessions s ON s.start_time::date = d::date
		GROUP BY d
		ORDER BY d
	`)
	if err == nil {
		for rows.Next() {
			var p DayPoint
			if err := rows.Scan(&p.Date, &p.Value); err == nil {
				g.SessionsPerDay = append(g.SessionsPerDay, p)
			}
		}
		rows.Close()
	}

	// New users per day, last 30 days
	rows, err = r.db.QueryContext(ctx, `
		SELECT to_char(d::date, 'YYYY-MM-DD') AS date,
		       COALESCE(COUNT(u.id), 0) AS cnt
		FROM generate_series(NOW()::date - INTERVAL '29 days', NOW()::date, INTERVAL '1 day') d
		LEFT JOIN users u ON u.created_at::date = d::date
		GROUP BY d
		ORDER BY d
	`)
	if err == nil {
		for rows.Next() {
			var p DayPoint
			if err := rows.Scan(&p.Date, &p.Value); err == nil {
				g.NewUsersPerDay = append(g.NewUsersPerDay, p)
			}
		}
		rows.Close()
	}

	// Top 10 skills by occurrence
	rows, err = r.db.QueryContext(ctx, `
		SELECT name, COUNT(*) AS cnt
		FROM skills
		GROUP BY name
		ORDER BY cnt DESC
		LIMIT 10
	`)
	if err == nil {
		for rows.Next() {
			var s SkillCount
			if err := rows.Scan(&s.Name, &s.Count); err == nil {
				g.TopSkills = append(g.TopSkills, s)
			}
		}
		rows.Close()
	}

	return g, nil
}

type ExportUserRow struct {
	ID         uuid.UUID
	Email      string
	Name       string
	IsVerified bool
	CreatedAt  string
	TotalMin   int
	Sessions   int
	Skills     int
}

func (r *AdminReportRepository) ExportUsers(ctx context.Context) ([]ExportUserRow, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT u.id, u.email, u.name, u.is_verified, u.created_at,
		       COALESCE((SELECT SUM(duration_minutes) FROM sessions WHERE user_id = u.id), 0),
		       COALESCE((SELECT COUNT(*) FROM sessions WHERE user_id = u.id), 0),
		       COALESCE((SELECT COUNT(*) FROM skills WHERE user_id = u.id), 0)
		FROM users u
		ORDER BY u.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []ExportUserRow{}
	for rows.Next() {
		var u ExportUserRow
		var created sql.NullTime
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.IsVerified, &created, &u.TotalMin, &u.Sessions, &u.Skills); err != nil {
			return nil, err
		}
		if created.Valid {
			u.CreatedAt = created.Time.UTC().Format("2006-01-02T15:04:05Z")
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

type ExportSessionRow struct {
	ID        uuid.UUID
	UserEmail string
	SkillName string
	StartTime string
	EndTime   string
	Minutes   int
	Type      string
}

func (r *AdminReportRepository) ExportSessions(ctx context.Context, limit int) ([]ExportSessionRow, error) {
	if limit <= 0 {
		limit = 10000
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT se.id, u.email, sk.name, se.start_time, se.end_time, se.duration_minutes, se.session_type
		FROM sessions se
		JOIN users u  ON u.id  = se.user_id
		JOIN skills sk ON sk.id = se.skill_id
		ORDER BY se.start_time DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []ExportSessionRow{}
	for rows.Next() {
		var r ExportSessionRow
		var start, end sql.NullTime
		if err := rows.Scan(&r.ID, &r.UserEmail, &r.SkillName, &start, &end, &r.Minutes, &r.Type); err != nil {
			return nil, err
		}
		if start.Valid {
			r.StartTime = start.Time.UTC().Format("2006-01-02T15:04:05Z")
		}
		if end.Valid {
			r.EndTime = end.Time.UTC().Format("2006-01-02T15:04:05Z")
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// DeleteUser removes a user. All related records cascade via ON DELETE CASCADE in 001_initial_schema.
func (r *AdminReportRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}
