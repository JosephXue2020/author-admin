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
      <el-table-column type="index" :index="getIndex" label="序号" width="180" align="center" />
      <el-table-column label="用户名" width="180" align="center">
        <template slot-scope="scope">
          {{ scope.row.name }}
        </template>
      </el-table-column>
      <el-table-column label="角色" width="180" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.role }}</span>
        </template>
      </el-table-column>
      <el-table-column label="创建人" width="180" align="center">
        <template slot-scope="scope">
          {{ scope.row.remark }}
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="180" align="center">
        <template slot-scope="scope">
          {{ scope.row.remark }}
        </template>
      </el-table-column>
      <el-table-column class-name="status-col" label="操作" align="center">
        <span>
          <template>
            <el-tag> 修改 </el-tag>
          </template>
          &nbsp; &nbsp; &nbsp; &nbsp;
          <template>
            <el-tag> 删除 </el-tag>
          </template>
        </span>
      </el-table-column>
    </el-table>
    <br>
    <div style="float:right">
      <el-pagination
        background
        :current-page="currentPage"
        :page-sizes="[10, 20, 50, 1]"
        :page-size="pageSize"
        layout="total, sizes, prev, pager, next, jumper"
        :total="totalPage"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script>
import { getUserList } from '@/api/user'

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
    return {
      list: null,
      listLoading: true,
      currentPage: 1, // 默认值
      pageSize: 10, // 默认值
      totalPage: null
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
      return (this.currentPage - 1) * this.pageSize + index
    }
  }
}
</script>
