package repositories

import (
	"database/sql"
	"sync"

	"github.com/EasyDarwin/EasyDarwin/models"
)

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

func (r *RecordRepository) GetRecords() ([]models.Record, error) {
	rows, err := r.db.Query("SELECT id, live_id, live_name, live_url, file_mp4, hls_url,file_record," +
		" create_time, update_time FROM t_records")
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
