<template>
  <div class="contact-container">
    <div class="contact-text">
      <el-card class="box-card">
        <div slot="header" class="clearfix">
          <span>导入</span>
          <el-button style="float: right; padding: 3px 0" type="text">
            <a href="/static/api/v1/author/template.xlsx">
              模板下载
            </a>
          </el-button>
        </div>
        <div class="text item" style="text-align: center">
          <el-upload
            class="upload-demo"
            drag
            action="/api/v1/author/upload"
            multiple
          >
            <i class="el-icon-upload" />
            <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
            <div slot="tip" class="el-upload__tip">只能上传jpg/png文件，且不超过500kb</div>
          </el-upload>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script>

export default {
  name: 'Upload'
}
</script>

<style lang="scss" scoped>
.contact {
  &-container {
    margin: 30px;
  }
  &-text {
    padding: 18px;
    width: 50%;
  }
}
.text {
    font-size: 16px;
  }

.item {
  margin-bottom: 18px;
}

.clearfix:before,
.clearfix:after {
  display: table;
  content: "";
}
.clearfix:after {
  clear: both
}

.box-card {
  width: 100%;
}
</style>
