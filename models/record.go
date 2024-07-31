package models

type Record struct {
	ID         string `json:"id"`
	LiveID     string `json:"live_id"`
	LiveName   string `json:"live_name"`
	LiveUrl    string `json:"-"`
	FileMp4    string `json:"-"`
	HlsUrl     string `json:"hls_url"`     //m3u8的ownload地址
	FileRecord string `json:"-"`           //文件的保存位置
	CreateTime string `json:"create_time"` //"2024-07-31 03:19:26"
	UpdateTime string `json:"update_time"` //"2024-07-31 03:19:26"
}
