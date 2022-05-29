<template>
  <div class="app-container">
    <el-table
      v-loading="listLoading"
      :data="list"
      element-loading-text="Loading"
      border
      fit
      highlight-current-row
    >
      <!-- <el-table-column align="center" label="序号" width="180">
        <template slot-scope="scope">
          {{ scope.$index }}
        </template>
      </el-table-column> -->
      <el-table-column type="index" :index="getIndex" label="序号" width="160" align="center" />
      <el-table-column label="用户名" width="160" align="center">
        <template slot-scope="scope">
          {{ scope.row.name }}
        </template>
      </el-table-column>
      <el-table-column label="角色" width="160" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.role }}</span>
        </template>
      </el-table-column>
      <el-table-column label="创建人" width="160" align="center">
        <template slot-scope="scope">
          {{ scope.row.creater }}
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="160" align="center">
        <template slot-scope="scope">
          {{ scope.row.createon }}
        </template>
      </el-table-column>
      <el-table-column label="部门" align="center">
        <template slot-scope="scope">
          {{ scope.row.department }}
        </template>
      </el-table-column>
      <el-table-column class-name="status-col" label="操作" width="160" align="center">
        <template slot-scope="scope">
          <el-button
            type="primary"
            size="mini"
            @click="openDialog(scope.row)"
          >
            修改
          </el-button>
          <el-button
            type="primary"
            size="mini"
            danger
            @click="deleteUser(scope.row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <br>
    <div style="float:right">
      <el-pagination
        background
        :current-page="currentPage"
        :page-sizes="pageSizes"
        :page-size="pageSize"
        layout="total, sizes, prev, pager, next, jumper"
        :total="totalPage"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <el-dialog :title="form.username" :visible.sync="dialogFormVisible" center>
      <el-form ref="updateUserForm" :model="form" :rules="rules">
        <el-form-item label="新密码" :label-width="formLabelWidth" prop="password">
          <el-input v-model="form.password" type="password" autocomplete="off" />
        </el-form-item>
        <el-form-item label="确认密码" :label-width="formLabelWidth" prop="checkPassword">
          <el-input v-model="form.checkPassword" type="password" autocomplete="off" />
        </el-form-item>
        <el-form-item label="角色" prop="role" :label-width="formLabelWidth">
          <el-select v-model="form.role" placeholder="请选择角色">
            <div v-for="(val, i) in form.roles" :key="i">
              <el-option :label="val" :value="val" />
            </div>
          </el-select>
          <div style="font-size: 5px;">
            仅可为用户选择权限小于操作者的角色
          </div>
          <!-- <el-input v-model.number="ruleForm.role" /> -->
        </el-form-item>
        <el-form-item label="部门" :label-width="formLabelWidth">
          <el-input v-model="form.department" autocomplete="off" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="updateUser">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { getUserList } from '@/api/user'
import { deleteUser as deleteUserApi } from '@/api/user'
import { updateUser as updateUserApi } from '@/api/user'

export default {
  name: 'List',
  filters: {
    statusFilter(status) {
      const statusMap = {
        published: 'success',
        draft: 'gray',
        deleted: 'danger'
      }
      return statusMap[status]
    }
  },
  data() {
    const pageSizes = [10, 20, 50, 100]
    const pageSize = pageSizes[0]

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
      } else if (value !== this.form.password) {
        callback(new Error('两次输入密码不一致'))
      } else {
        callback()
      }
    }

    return {
      list: null,
      listLoading: true,
      currentPage: 1, // 默认值
      pageSizes: pageSizes,
      pageSize: pageSize, // 默认值
      totalPage: null,

      dialogFormVisible: false,
      form: {
        id: -1,
        username: '',
        password: '',
        checkPassword: '',
        role: '',
        roles: ['super', 'admin', 'normal', 'guest'],
        department: ''
      },
      formLabelWidth: '120px',
      rules: {
        password: [
          { validator: validatePass, trigger: 'blur' }
        ],
        checkPassword: [
          { validator: validatePass2, trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    const pageNum = this.currentPage
    this.fetchData(pageNum, this.pageSize)
  },
  methods: {
    handleSizeChange(val) {
      this.pageSize = val
      this.fetchData(1, val)
    },
    handleCurrentChange(val) {
      this.fetchData(val, this.pageSize)
    },
    fetchData(pageNum, pageSize) {
      var params = new URLSearchParams()
      params.append('pageNum', pageNum)
      params.append('pageSize', pageSize)

      this.listLoading = true
      getUserList(params).then(response => {
        this.list = response.data.items
        this.listLoading = false
        this.totalPage = response.data.total
        this.currentPage = pageNum
      })
    },
    getIndex(index) {
      return (this.currentPage - 1) * this.pageSize + index + 1
    },
    deleteConfirm(n) {
      return this.$confirm(`此操作将删除 ${n}, 是否继续？`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
    },
    deleteUser(row) {
      this.deleteConfirm(row.name).then(() => {
        deleteUserApi(row.id).then(() => {
          this.$message({
            type: 'success',
            message: '删除成功.'
          })
          this.fetchData(this.currentPage, this.pageSize)
        })
      })
    },
    resetForm(n) {
      this.$refs[n].resetFields()
    },
    openDialog(row) {
      // this.resetForm('updateUserForm')
      this.dialogFormVisible = true
      this.form.id = row.id
      this.form.username = row.name
    },
    updateUser(row) {
      const data = {
        id: this.form.id,
        username: this.form.username,
        password: this.form.pass,
        role: this.form.role,
        department: this.form.department
      }
      updateUserApi(data).then(() => this.fetchData(this.currentPage, this.pageSize))
      this.dialogFormVisible = false
    }
  }
}
</script>
