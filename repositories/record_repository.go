package repositories

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/EasyDarwin/EasyDarwin/models"
)

var allowedSortFields = map[string]bool{
	"id":          true,
	"live_id":     true,
	"live_name":   true,
	"live_url":    true,
	"file_mp4":    true,
	"hls_url":     true,
	"file_record": true,
	"create_time": true,
	"update_time": true,
}

type RecordRepository struct {
	db *sql.DB
}

var (
	instance *RecordRepository
	once     sync.Once
)

func GetUserRepository(db *sql.DB) *RecordRepository {
	once.Do(func() {
		instance = &RecordRepository{db}
	})
	return instance
}

func NewRecordRepository(db *sql.DB) *RecordRepository {
	return &RecordRepository{db: db}
}

func (r *RecordRepository) CreateRecord(record models.Record) error {
	stmt, err := r.db.Prepare("INSERT INTO t_records ( live_id, live_name, live_url, file_mp4, hls_url," +
		"file_record, create_time, update_time) VALUES ( ?, ? , ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(record.LiveID, record.LiveName, record.LiveUrl, record.FileMp4, record.HlsUrl,
		record.FileRecord, record.CreateTime, record.UpdateTime)
	if err != nil {
		return err
	}
	return nil
}

func (r *RecordRepository) GetRecords(sortBy string, order string, pageSize int, offset int, q string) ([]models.Record, int, error) {

	// Validate sortBy and order parameters
	if sortBy != "" && !allowedSortFields[sortBy] {
		return nil, 0, fmt.Errorf("invalid sort field: %s", sortBy)
	}
	if order != "asc" && order != "desc" {
		order = "asc" // default order
	}

	// 初始化基本查询部件
	countQuery := "SELECT COUNT(*) FROM t_records"
	baseQuery := "SELECT id, live_id, live_name, live_url, file_mp4, hls_url, file_record, create_time, update_time FROM t_records"
	whereClause := ""
	orderClause := ""
	limitClause := ""
	args := []interface{}{}

	// 如果提供了q，则构造where子句
	if q != "" {
		whereClause = " WHERE live_name LIKE ? OR live_id LIKE ?"
		args = append(args, "%"+q+"%", "%"+q+"%")
	}

	// 如果提供sortBy，则添加排序
	if sortBy != "" {
		orderClause = fmt.Sprintf(" ORDER BY %s %s", sortBy, order)
	}

	// 如果提供了页面大小和偏移量，则添加分页
	if pageSize > 0 && offset >= 0 {
		limitClause = fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)
	}

	// 获取总计数
	if whereClause != "" {
		countQuery += whereClause
	}
	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// 构建最终查询
	finalQuery := baseQuery + whereClause + orderClause + limitClause

	rows, err := r.db.Query(finalQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var records []models.Record
	for rows.Next() {
		var record models.Record
		err := rows.Scan(&record.ID, &record.LiveID, &record.LiveName, &record.LiveUrl, &record.FileMp4, &record.HlsUrl, &record.FileRecord, &record.CreateTime, &record.UpdateTime)
		if err != nil {
			return nil, 0, err
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (r *RecordRepository) GetRecordsByLiveID(liveID string) ([]models.Record, error) {
	rows, err := r.db.Query("SELECT id, live_id, live_name, live_url, file_mp4, hls_url,file_record,"+
		" create_time, update_time FROM t_records WHERE live_id = ?", liveID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []models.Record{}
	for rows.Next() {
		var id string
		var liveID string
		var liveName string
		var liveUrl string
		var fileMp4 string
		var hlsUrl string
		var fileRecord string
		var createTime string
		var updateTime string

		err = rows.Scan(&id, &liveID, &liveName, &liveUrl, &fileMp4, &hlsUrl, &fileRecord, &createTime, &updateTime)
		if err != nil {
			return nil, err
		}
		records = append(records, models.Record{
			ID:         id,
			LiveID:     liveID,
			LiveName:   liveName,
			LiveUrl:    liveUrl,
			FileMp4:    fileMp4,
			HlsUrl:     hlsUrl,
			FileRecord: fileRecord,
			CreateTime: createTime,
			UpdateTime: updateTime,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (r *RecordRepository) GetRecordsById(id string) (models.Record, error) {
	row := r.db.QueryRow("SELECT id, live_id, live_name, live_url, file_mp4, hls_url,file_record,create_time, "+
		"update_time FROM t_records WHERE id = ?", id)
	var record models.Record
	err := row.Scan(&record.ID, &record.LiveID, &record.LiveName, &record.LiveUrl, &record.FileMp4, &record.HlsUrl,
		&record.FileRecord, &record.CreateTime, &record.UpdateTime)
	if err != nil {
		return models.Record{}, err
	}
	return record, nil
}

func (r *RecordRepository) DeleteRecord(id string) error {
	stmt, err := r.db.Prepare("DELETE FROM t_records WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
