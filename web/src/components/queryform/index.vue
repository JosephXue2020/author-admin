<template>
  <div class="queryform">
    <div class="demo-input-suffix">
      姓名：
      <el-input
        v-model="input2"
        placeholder="请输入内容"
        prefix-icon="el-icon-search"
      />
    </div>
    <div class="demo-input-suffix">
      毕业院校：
      <el-input
        v-model="input4"
        placeholder="请输入内容"
      >
        <i slot="prefix" class="el-input__icon el-icon-search" />
      </el-input>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Queryform',
  data() {
    return {
      input1: '',
      input2: '',
      input3: '',
      input4: ''
    }
  }
}
</script>

<style scoped>
/* .queryform {
  width: 100%;
  height: 60px;
} */
</style>
