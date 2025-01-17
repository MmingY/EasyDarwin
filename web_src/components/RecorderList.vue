<template>
  <div class="container-fluid no-padding">
    <div class="box box-success">
      <div class="box-header">
        <h4 class="text-success text-center">录像列表</h4>
        <form class="form-inline">
          <div class="form-group pull-right">
            <div class="input-group">
              <input type="text" class="form-control" placeholder="搜索" v-model.trim="q" @keydown.enter.prevent ref="q">
              <div class="input-group-btn">
                <button type="button" class="btn btn-default" @click.prevent="doSearch">
                  <i class="fa fa-search"></i>
                </button>
              </div>
            </div>
          </div>
        </form>
      </div>
      <div class="box-body">
        <el-table :data="records" stripe class="view-list" :default-sort="{ prop: 'create_time', order: 'descending' }"
          @sort-change="sortChange" :header-cell-style="{ 'text-align': 'center' }"
          :cell-style="{ 'text-align': 'center' }">
          <el-table-column prop="live_id" label="ID" min-width="120"></el-table-column>
          <el-table-column prop="hls_url" label="播放地址" min-width="240" show-overflow-tooltip>
            <template slot-scope="scope">
              <span>
                <i class="fa fa-copy" role="button" v-clipboard="scope.row.hls_url" title="点击拷贝"
                  @success="$message({ type: 'success', message: '成功拷贝到粘贴板' })"></i>

                <span @click="playVideo(scope.row.hls_url)" style="color:#0cbb92" class="play-url"> {{
                  scope.row.hls_url }}</span>
              </span>
            </template>
          </el-table-column>
          <el-table-column label="传输方式" min-width="100"></el-table-column>
          <!-- <el-table-column prop="inBytes" label="上行流量" min-width="120" :formatter="formatBytes" sortable="custom"></el-table-column> -->
          <!-- <el-table-column label="下行流量" min-width="120" :formatter="formatBytes" sortable="custom"></el-table-column> -->
          <el-table-column prop="create_time" label="创建时间" min-width="200" sortable="custom"></el-table-column>

          <el-table-column label="操作" min-width="150">
            <template slot-scope="scope">
              <div class="operation-buttons">
                <el-button @click="playVideo(scope.row.hls_url)" size="mini" type="primary">播放</el-button>
                <el-button @click="deleteRecord(scope.row.id)" size="mini" type="danger">删除</el-button>
                <el-button @click="downloadRecord(scope.row.hls_url)" size="mini" type="warning">下载</el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <div class="box-footer clearfix" v-if="total > 0">
        <el-pagination layout="prev,pager,next" class="pull-right" :total="total" :page-size.sync="pageSize"
          :current-page.sync="currentPage"></el-pagination>
      </div>
    </div>

    <el-dialog custom-class="my-dialog" title="录像回放" :visible.sync="dialogVisible" width="50%" height="40%"
      :before-close="handleBeforeClose" style="font-weight:bold;">
      <video ref="videoPlayer" v-if="currentVideo" :src="currentVideo" controls autoplay
        style="width: 100%;height:100%"></video>
    </el-dialog>
  </div>
</template>

<script>
import prettyBytes from "pretty-bytes";
import Hls from 'hls.js';
import { MessageBox, Message } from 'element-ui';

import _ from "lodash";
export default {
  props: {},
  data() {
    return {
      total: 0,          // 用于分页的总条目数
      pageSize: 10,
      code: "",
      msg: "OK",
      data: [],
      dialogVisible: false,
      currentVideo: "",
      records: [],
      currentPage: 1,
      q: "",
      sort: "create_time",
      order: "descending",
    };
  },
  beforeDestroy() {
    if (this.timer) {
      clearInterval(this.timer);
      this.timer = 0;
    }
  },
  mounted() {
    this.$refs["q"].focus();
    this.timer = setInterval(() => {
      this.getRecords();
    }, 3000);
  },
  watch: {
    q: function (newVal, oldVal) {
      this.doDelaySearch();
    },
    currentPage: function (newVal, oldVal) {
      this.doSearch(newVal);
    }
  },
  methods: {
    playVideo(path) {
      this.currentVideo = path;
      this.dialogVisible = true;

      this.$nextTick(() => {
        const video = this.$refs.videoPlayer;
        if (Hls.isSupported()) {
          const hls = new Hls();
          hls.loadSource(this.currentVideo);
          hls.attachMedia(video);
          hls.on(Hls.Events.MANIFEST_PARSED, () => {
            video.play();
          });
        } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
          video.src = this.currentVideo;
          video.addEventListener('loadedmetadata', () => {
            video.play();
          });
        }
      });
    },
    getRecords() {
      //  $.get("/api/v1/record").then(data => {
      $.get("/record/query", {
        q: this.q,
        start: (this.currentPage - 1) * this.pageSize,
        limit: this.pageSize,
        sort: this.sort,
        order: this.order
      }).then(data => {
        this.records = data.list;
        this.total = data.total;
        /*  if (0 === data.code) {
           this.records = data.data;
         } */
      });
    },
    downloadRecord(hslPath) {
      //"/record/ccc/20240802/20240802090544/record.m3u8"
      console.log("downloadRecord", hslPath);
      // Step 1: 去掉文件名和扩展名 "/record/ccc/20240802/20240802090544"
      let pathWithoutFile = hslPath.slice(0, hslPath.lastIndexOf('/'));

      // Step 2: 在路径中插入 "/download"
      let finalPath = pathWithoutFile.replace("/record", "/record/download");
      // 实现下载视频的逻辑
      const link = document.createElement('a');
      link.href = finalPath;
      link.download = 'video.mp4'; // 可以根据实际情况设置下载文件的名称
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    },
    deleteRecord(recordId) {
      console.log("deleteRecord", recordId);
      // 弹出确认对话框
      MessageBox.confirm('是否确定删除该录像?', '警告', {
        confirmButtonText: '是',
        cancelButtonText: '否',
        type: 'warning'
      }).then(() => {
        $.ajax({
          url: '/record/remove/' + recordId, // 替换为你的 API 端点
          type: 'DELETE', // 请求方法
          success: function (response) {
            // 从记录中移除被删除的数据
            this.records = this.records.filter(item => item.id !== recordId);
            // 请求成功
            Message({
              type: 'success',
              message: '录像被删除成功'
            });
          }.bind(this),
          error: function (xhr, status, error) {
            // 显示错误消息
            Message({
              type: 'error',
              message: 'Failed to delete resource.${error}'
            });
          }
        });
      }).catch(() => { });
    },
    handleBeforeClose(done) {
      // 在弹窗关闭前执行的逻辑
      const video = this.$refs.videoPlayer;
      video.pause(); // 停止视频播放
      video.currentTime = 0; // 将播放位置重置为开始位置
      video.src = ''; // 清除视频源
      video.load(); // 重新加载视频（确保完全清除视频流）
      done(); // 继续关闭弹窗
    },
    doSearch(page = 1) {
      var query = {};
      if (this.q) query["q"] = this.q;
      this.$router.replace({
        path: `/recorders/${page}`,
        query: query
      });
    },
    doDelaySearch: _.debounce(function () {
      this.doSearch();
    }, 500),
    sortChange(data) {
      this.sort = data.prop;
      this.order = data.order;
      this.getRecords();
    },
    formatBytes(row, col, val) {
      if (val == undefined) return "-";
      return prettyBytes(val);
    }
  },
  beforeRouteEnter(to, from, next) {
    next(vm => {
      vm.q = to.query.q || "";
      vm.currentPage = parseInt(to.params.page) || 1;
    });
  },
  beforeRouteUpdate(to, from, next) {
    next();
    this.$nextTick(() => {
      this.q = to.query.q || "";
      this.currentPage = parseInt(to.params.page) || 1;
      this.records = [];
      this.getRecords();
    });
  }
};
</script>

<style lang="css" scoped="true">
.play-url {
  color: #0cbb92;
  /* max-width: 400px; */
  /* 设置最大宽度 */
  /* display: inline-block; */
  /* overflow: hidden; */
  /* text-overflow: ellipsis; */
  /* 溢出文本显示省略号 */
  /* white-space: nowrap; */
  /* 防止文本换行 */
}

.my-dialog .el-dialog__header {
  text-align: center;
}

/* .center-column .cell {
  text-align: center;
}

el-table-column {
  text-align: center;
} */
</style>