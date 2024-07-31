package routers

import (
	"bytes"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"path"

	"github.com/EasyDarwin/EasyDarwin/models"
	"github.com/EasyDarwin/EasyDarwin/repositories"
	"github.com/gin-gonic/gin"
	"github.com/penggy/EasyGoLib/db"
	"github.com/penggy/EasyGoLib/utils"
)

/**
 * @apiDefine record 录像
 */

/**
 * @apiDefine fileInfo
 * @apiSuccess (200) {String} duration	格式化好的录像时长
 * @apiSuccess (200) {Number} durationMillis	录像时长，毫秒为单位
 * @apiSuccess (200) {String} path 录像文件的相对路径,其绝对路径为：http[s]://host:port/record/[path]。
 * @apiSuccess (200) {String} folder 录像文件夹，录像文件夹以推流路径命名。
 */

/**
 * @api {get} /api/v1/record/folders 获取所有录像文件夹
 * @apiGroup record
 * @apiName RecordFolders
 * @apiParam {Number} [start] 分页开始,从零开始
 * @apiParam {Number} [limit] 分页大小
 * @apiParam {String} [sort] 排序字段
 * @apiParam {String=ascending,descending} [order] 排序顺序
 * @apiParam {String} [q] 查询参数
 * @apiSuccess (200) {Number} total 总数
 * @apiSuccess (200) {Array} rows 文件夹列表
 * @apiSuccess (200) {String} rows.folder	录像文件夹名称
 */
func (h *APIHandler) RecordFolders(c *gin.Context) {
	mp4Path := utils.Conf().Section("rtsp").Key("m3u8_dir_path").MustString("")
	form := utils.NewPageForm()
	if err := c.Bind(form); err != nil {
		log.Printf("record folder bind err:%v", err)
		return
	}
	var files = make([]interface{}, 0)
	if mp4Path != "" {
		visit := func(files *[]interface{}) filepath.WalkFunc {
			return func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if path == mp4Path {
					return nil
				}
				if !info.IsDir() {
					return nil
				}
				*files = append(*files, map[string]interface{}{"folder": info.Name()})
				return filepath.SkipDir
			}
		}
		err := filepath.Walk(mp4Path, visit(&files))
		if err != nil {
			log.Printf("Query RecordFolders err:%v", err)
		}
	}
	pr := utils.NewPageResult(files)
	if form.Sort != "" {
		pr.Sort(form.Sort, form.Order)
	}
	pr.Slice(form.Start, form.Limit)
	c.IndentedJSON(200, pr)

}

/**
 * @api {get} /api/v1/record/files 获取所有录像文件
 * @apiGroup record
 * @apiName RecordFiles
 * @apiParam {Number} folder 录像文件所在的文件夹
 * @apiParam {Number} [start] 分页开始,从零开始
 * @apiParam {Number} [limit] 分页大小
 * @apiParam {String} [sort] 排序字段
 * @apiParam {String=ascending,descending} [order] 排序顺序
 * @apiParam {String} [q] 查询参数
 * @apiSuccess (200) {Number} total 总数
 * @apiSuccess (200) {Array} rows 文件列表
 * @apiSuccess (200) {String} rows.duration	格式化好的录像时长
 * @apiSuccess (200) {Number} rows.durationMillis	录像时长，毫秒为单位
 * @apiSuccess (200) {String} rows.path 录像文件的相对路径,录像文件为m3u8格式，将其放到video标签中便可直接播放。其绝对路径为：http[s]://host:port/record/[path]。
 */
func (h *APIHandler) RecordFiles(c *gin.Context) {
	type Form struct {
		utils.PageForm
		Folder  string `form:"folder" binding:"required"`
		StartAt int    `form:"beginUTCSecond"`
		StopAt  int    `form:"endUTCSecond"`
	}
	var form = Form{}
	form.Limit = math.MaxUint32
	err := c.Bind(&form)
	if err != nil {
		log.Printf("record file bind err:%v", err)
		return
	}

	files := make([]interface{}, 0)
	mp4Path := utils.Conf().Section("rtsp").Key("m3u8_dir_path").MustString("")
	if mp4Path != "" {
		ffmpeg_path := utils.Conf().Section("rtsp").Key("ffmpeg_path").MustString("")
		ffmpeg_folder, executable := filepath.Split(ffmpeg_path)
		split := strings.Split(executable, ".")
		suffix := ""
		if len(split) > 1 {
			suffix = split[1]
		}
		ffprobe := ffmpeg_folder + "ffprobe" + suffix
		folder := filepath.Join(mp4Path, form.Folder)
		visit := func(files *[]interface{}) filepath.WalkFunc {
			return func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if path == folder {
					return nil
				}
				if info.IsDir() {
					return nil
				}
				if info.Size() == 0 {
					return nil
				}
				if info.Name() == ".DS_Store" {
					return nil
				}
				if !strings.HasSuffix(strings.ToLower(info.Name()), ".m3u8") && !strings.HasSuffix(strings.ToLower(info.Name()), ".ts") {
					return nil
				}
				cmd := exec.Command(ffprobe, "-i", path)
				cmdOutput := &bytes.Buffer{}
				//cmd.Stdout = cmdOutput
				cmd.Stderr = cmdOutput
				err = cmd.Run()
				bytes := cmdOutput.Bytes()
				output := string(bytes)
				//log.Printf("%v result:%v", cmd, output)
				var average = regexp.MustCompile(`Duration: ((\d+):(\d+):(\d+).(\d+))`)
				result := average.FindStringSubmatch(output)
				duration := time.Duration(0)
				durationStr := ""
				if len(result) > 0 {
					durationStr = result[1]
					h, _ := strconv.Atoi(result[2])
					duration += time.Duration(h) * time.Hour
					m, _ := strconv.Atoi(result[3])
					duration += time.Duration(m) * time.Minute
					s, _ := strconv.Atoi(result[4])
					duration += time.Duration(s) * time.Second
					millis, _ := strconv.Atoi(result[5])
					duration += time.Duration(millis) * time.Millisecond
				}

				*files = append(*files, map[string]interface{}{
					"path":           path[len(mp4Path):],
					"durationMillis": duration / time.Millisecond,
					"duration":       durationStr})
				return nil
			}
		}
		err = filepath.Walk(folder, visit(&files))
		if err != nil {
			log.Printf("Query RecordFolders err:%v", err)
		}
	}

	pr := utils.NewPageResult(files)
	if form.Sort != "" {
		pr.Sort(form.Sort, form.Order)
	}
	pr.Slice(form.Start, form.Limit)
	c.IndentedJSON(200, pr)
}

/**
 * @api {get} /record/download/*anyPath 录像m3u8合成mp4文件,并下载
 * @apiGroup record
 * @apiName RecordVie0s
 */
func (h *APIHandler) RecordDownload(c *gin.Context) {
	ffmpeg := utils.Conf().Section("rtsp").Key("ffmpeg_path").MustString("")
	m3u8_dir_path := utils.Conf().Section("rtsp").Key("m3u8_dir_path").MustString("")
	// 提取 /user 之后的部分
	// userPath := requestPath[len("/record/download/"):]
	// userPath := strings.Replace(c.Request.URL.Path, "/record/download/", "", 1)
	userPath := c.Param("anyPath")
	pathslice := strings.Split(userPath, "/")
	var fileName string
	if len(pathslice[len(pathslice)-1]) == 0 {
		fileName = pathslice[len(pathslice)-2]
	} else {
		fileName = pathslice[len(pathslice)-1]
	}
	// fileMp4 := m3u8_dir_path + userPath + fileName + ".mp4"
	fileMp4 := path.Join(m3u8_dir_path, userPath, fileName+".mp4")
	fileM3u8 := path.Join(m3u8_dir_path, userPath, "record.m3u8")
	//ffmpeg -i input.m3u8 -c copy output.mp4
	params := []string{"-i", fileM3u8, "-c", "copy", fileMp4}
	cmd := exec.Command(ffmpeg, params...)
	err := cmd.Start()
	if err != nil {
		log.Printf("Start ffmpeg err:%v", err)
	}

	//发送mp4文件给客户端
	file, err := os.Open(fileMp4)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("Open mp4 err:%v", err)
		return
	}
	defer file.Close()

	// Get the file size and content type
	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fileSize := fileInfo.Size()
	contentType := "video/mp4"

	// Set the response headers
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(fileMp4))
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", strconv.FormatInt(fileSize, 10))

	// Send the file content
	if _, err := io.Copy(c.Writer, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

/**
 * @api {get} /record/query/:id 录像m3u8合成mp4文件,并下载
 * @apiGroup record
 * @apiName RecordVie0s
 */
func (h *APIHandler) RecordQuery(c *gin.Context) {
	liveID := c.Param("liveID")
	repositories := repositories.GetUserRepository(db.SQLite.DB())
	if len(liveID) == 0 {
		if records, err := repositories.GetRecords(); err == nil {
			response := struct {
				List []models.Record `json:"list"`
			}{
				List: records,
			}
			c.JSON(http.StatusOK, response)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if records, err := repositories.GetRecordsByLiveID(liveID); err == nil {
		// c.JSON(http.StatusOK, record)
		response := struct {
			liveId string          `json:"liveID"`
			List   []models.Record `json:"list"`
		}{
			liveId: liveID,
			List:   records,
		}
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
