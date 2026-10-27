package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"archival-vault/backend/database"
	"archival-vault/backend/middleware"
	"archival-vault/backend/models"
)

func ListRecords(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	search := strings.TrimSpace(q.Get("search"))
	category := strings.TrimSpace(q.Get("category"))
	class := strings.TrimSpace(q.Get("classification"))
	status := strings.TrimSpace(q.Get("status"))

	query := `SELECT id, code, title, classification, category, status, tags, summary, payload, created_by, created_at, updated_at 
	          FROM records WHERE 1=1`
	var args []any

	if search != "" {
		query += ` AND (title LIKE ? OR summary LIKE ? OR tags LIKE ? OR code LIKE ?)`
		like := "%" + search + "%"
		args = append(args, like, like, like, like)
	}
	if category != "" && category != "ALL" {
		query += ` AND category = ?`
		args = append(args, category)
	}
	if class != "" && class != "ALL" {
		query += ` AND classification = ?`
		args = append(args, class)
	}
	if status != "" && status != "ALL" {
		query += ` AND status = ?`
		args = append(args, status)
	}

	query += ` ORDER BY id DESC LIMIT 200`

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		middleware.JSONError(w, "query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	records := make([]models.Record, 0)
	for rows.Next() {
		var rec models.Record
		err := rows.Scan(
			&rec.ID, &rec.Code, &rec.Title, &rec.Classification, &rec.Category,
			&rec.Status, &rec.Tags, &rec.Summary, &rec.Payload, &rec.CreatedBy,
			&rec.CreatedAt, &rec.UpdatedAt,
		)
		if err != nil {
			middleware.JSONError(w, "scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		records = append(records, rec)
	}

	middleware.JSONSuccess(w, http.StatusOK, records, "records retrieved")
}

func GetRecord(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.JSONError(w, "invalid record id", http.StatusBadRequest)
		return
	}

	var rec models.Record
	row := database.DB.QueryRow(
		`SELECT id, code, title, classification, category, status, tags, summary, payload, created_by, created_at, updated_at 
		 FROM records WHERE id = ?`, id,
	)
	err = row.Scan(
		&rec.ID, &rec.Code, &rec.Title, &rec.Classification, &rec.Category,
		&rec.Status, &rec.Tags, &rec.Summary, &rec.Payload, &rec.CreatedBy,
		&rec.CreatedAt, &rec.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			middleware.JSONError(w, "record not found", http.StatusNotFound)
			return
		}
		middleware.JSONError(w, "database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	middleware.JSONSuccess(w, http.StatusOK, rec, "record retrieved")
}

func CreateRecord(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetUserFromContext(r)
	creator := "ANONYMOUS"
	if claims != nil && claims.Username != "" {
		creator = claims.Username
	}

	var req models.RecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.JSONError(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		middleware.JSONError(w, "record title is required", http.StatusBadRequest)
		return
	}

	if req.Classification == "" {
		req.Classification = "INTERNAL"
	}
	if req.Category == "" {
		req.Category = "SPEC"
	}
	if req.Status == "" {
		req.Status = "ACTIVE"
	}

	var maxID sql.NullInt64
	_ = database.DB.QueryRow("SELECT MAX(id) FROM records").Scan(&maxID)
	nextSeq := int64(1)
	if maxID.Valid {
		nextSeq = maxID.Int64 + 1
	}
	code := fmt.Sprintf("ARC-%d-%03d", time.Now().Year(), nextSeq)

	res, err := database.DB.Exec(
		`INSERT INTO records (code, title, classification, category, status, tags, summary, payload, created_by)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		code, req.Title, req.Classification, req.Category, req.Status, req.Tags, req.Summary, req.Payload, creator,
	)
	if err != nil {
		middleware.JSONError(w, "failed to insert record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	created := models.Record{
		ID:             id,
		Code:           code,
		Title:          req.Title,
		Classification: req.Classification,
		Category:       req.Category,
		Status:         req.Status,
		Tags:           req.Tags,
		Summary:        req.Summary,
		Payload:        req.Payload,
		CreatedBy:      creator,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	middleware.JSONSuccess(w, http.StatusCreated, created, "record created")
}

func UpdateRecord(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.JSONError(w, "invalid record id", http.StatusBadRequest)
		return
	}

	var req models.RecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.JSONError(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		middleware.JSONError(w, "record title cannot be empty", http.StatusBadRequest)
		return
	}

	res, err := database.DB.Exec(
		`UPDATE records SET title = ?, classification = ?, category = ?, status = ?, tags = ?, summary = ?, payload = ?
		 WHERE id = ?`,
		req.Title, req.Classification, req.Category, req.Status, req.Tags, req.Summary, req.Payload, id,
	)
	if err != nil {
		middleware.JSONError(w, "failed to update record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAff, _ := res.RowsAffected()
	if rowsAff == 0 {
		var exists bool
		_ = database.DB.QueryRow("SELECT 1 FROM records WHERE id = ?", id).Scan(&exists)
		if !exists {
			middleware.JSONError(w, "record not found", http.StatusNotFound)
			return
		}
	}

	// Fetch updated
	var rec models.Record
	row := database.DB.QueryRow(
		`SELECT id, code, title, classification, category, status, tags, summary, payload, created_by, created_at, updated_at 
		 FROM records WHERE id = ?`, id,
	)
	_ = row.Scan(
		&rec.ID, &rec.Code, &rec.Title, &rec.Classification, &rec.Category,
		&rec.Status, &rec.Tags, &rec.Summary, &rec.Payload, &rec.CreatedBy,
		&rec.CreatedAt, &rec.UpdatedAt,
	)

	middleware.JSONSuccess(w, http.StatusOK, rec, "record updated")
}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.JSONError(w, "invalid record id", http.StatusBadRequest)
		return
	}

	res, err := database.DB.Exec("DELETE FROM records WHERE id = ?", id)
	if err != nil {
		middleware.JSONError(w, "failed to delete record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAff, _ := res.RowsAffected()
	if rowsAff == 0 {
		middleware.JSONError(w, "record not found", http.StatusNotFound)
		return
	}

	middleware.JSONSuccess(w, http.StatusOK, map[string]int64{"deleted_id": id}, "record deleted")
}

func GetStats(w http.ResponseWriter, r *http.Request) {
	var total int
	_ = database.DB.QueryRow("SELECT COUNT(*) FROM records").Scan(&total)

	byClass := make(map[string]int)
	classRows, _ := database.DB.Query("SELECT classification, COUNT(*) FROM records GROUP BY classification")
	if classRows != nil {
		defer classRows.Close()
		for classRows.Next() {
			var k string
			var v int
			if err := classRows.Scan(&k, &v); err == nil {
				byClass[k] = v
			}
		}
	}

	byCat := make(map[string]int)
	catRows, _ := database.DB.Query("SELECT category, COUNT(*) FROM records GROUP BY category")
	if catRows != nil {
		defer catRows.Close()
		for catRows.Next() {
			var k string
			var v int
			if err := catRows.Scan(&k, &v); err == nil {
				byCat[k] = v
			}
		}
	}

	byStatus := make(map[string]int)
	statusRows, _ := database.DB.Query("SELECT status, COUNT(*) FROM records GROUP BY status")
	if statusRows != nil {
		defer statusRows.Close()
		for statusRows.Next() {
			var k string
			var v int
			if err := statusRows.Scan(&k, &v); err == nil {
				byStatus[k] = v
			}
		}
	}

	var latestCode string
	_ = database.DB.QueryRow("SELECT code FROM records ORDER BY id DESC LIMIT 1").Scan(&latestCode)

	stats := models.StatsResponse{
		TotalRecords:  total,
		ByClass:       byClass,
		ByCategory:    byCat,
		ByStatus:      byStatus,
		LatestCode:    latestCode,
		UptimeSeconds: int64(time.Since(database.StartTime).Seconds()),
	}

	middleware.JSONSuccess(w, http.StatusOK, stats, "telemetry stats retrieved")
}
