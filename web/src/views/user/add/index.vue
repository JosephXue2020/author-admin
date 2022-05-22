<template>
  <div class="adduser-container">
    <div class="adduser-text">
      <el-form ref="ruleForm" :model="ruleForm" status-icon :rules="rules" label-width="100px" class="demo-ruleForm">
        <el-form-item label="用户名" prop="name">
          <el-input v-model="ruleForm.name" type="text" autocomplete="off" />
        </el-form-item>
        <el-form-item label="密码" prop="pass">
          <el-input v-model="ruleForm.pass" type="password" autocomplete="off" />
        </el-form-item>
        <el-form-item label="确认密码" prop="checkPass">
          <el-input v-model="ruleForm.checkPass" type="password" autocomplete="off" />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="ruleForm.role" placeholder="请选择角色">
            <div v-for="(val, i) in ruleForm.roles" :key="i">
              <el-option :label="val" :value="val" />
            </div>
          </el-select>
          <div style="font-size: 5px;">
            仅可创建权限小于操作者的用户
          </div>
          <!-- <el-input v-model.number="ruleForm.role" /> -->
        </el-form-item>
        <el-form-item label="部门" prop="department">
          <el-input v-model="ruleForm.department" />
        </el-form-item>
        <el-form-item>
          <div style="text-align: center;">
            <el-button type="primary" @click="submitForm('ruleForm')">提交</el-button>
            <el-button @click="resetForm('ruleForm')">重置</el-button>
          </div>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>
import { addUser } from '@/api/user'

export default {
  data() {
    var checkName = (rule, value, callback) => {
      if (!value) {
        return callback(new Error('用户名不能为空'))
      }
      setTimeout(() => {
        if (value.length > 50) {
          return callback(new Error('用户名长度不能大于50字符'))
        } else {
          callback()
        }
      }, 1000)
    }
    var validatePass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入密码'))
      } else if (value.length < 6) {
        callback(new Error('密码不能小于6位'))
      } else {
        if (this.ruleForm.checkPass !== '') {
          this.$refs.ruleForm.validateField('checkPass')
        }
        callback()
      }
    }
    var validatePass2 = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请再次输入密码'))
      } else if (value !== this.ruleForm.pass) {
        callback(new Error('两次输入密码不一致'))
      } else {
        callback()
      }
    }
    return {
      ruleForm: {
        name: '',
        pass: '',
        checkPass: '',
        role: '',
        roles: ['super', 'admin', 'normal', 'guest'],
        department: ''
      },
      rules: {
        name: [
          { validator: checkName, trigger: 'blur' }
        ],
        pass: [
          { validator: validatePass, trigger: 'blur' }
        ],
        checkPass: [
          { validator: validatePass2, trigger: 'blur' }
        ]
      }
    }
  },
  methods: {
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          const data = {
            username: this.ruleForm.name,
            password: this.ruleForm.pass,
            role: this.ruleForm.role,
            department: this.ruleForm.department
          }
          addUser(data)
          alert('提交成功')
          this.resetForm(formName)
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    resetForm(formName) {
      this.$refs[formName].resetFields()
    }
  }
}
</script>

<style lang="scss" scoped>
.adduser {
  &-container {
    margin: 30px;
  }
  &-text {
    margin: auto;
    width: 50%;
    font-size: 30px;
    line-height: 46px;
  }
}
</style>
