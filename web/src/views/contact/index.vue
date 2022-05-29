<template>
  <div class="contact-container">
    <div class="contact-text">
      <el-card class="box-card">
        <div slot="header" class="clearfix">
          <span>管理员</span>
        <!-- <el-button style="float: right; padding: 3px 0" type="text">添加</el-button> -->
        </div>
        <div class="text item">
          <p>联系人：Andrew NG</p>
          <p>电  话：88390603</p>
          <p>Email：AndrewNG@qq.com</p>
        </div>
      </el-card>
    </div>
    <div class="contact-text">
      <el-card class="box-card">
        <div slot="header" class="clearfix">
          <span>管理员</span>
        <!-- <el-button style="float: right; padding: 3px 0" type="text">添加</el-button> -->
        </div>
        <div class="text item">
          <p>联系人：Hinton</p>
          <p>电  话：88390603</p>
          <p>Email：Hinton@qq.com</p>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script>

export default {
  name: 'Contact'
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
