<template>
  <div class="app-container">
    <div>
      <el-queryform />
    </div>
    <div style="margin-top: 20px; float:right">
      <el-button type="primary" size="mini" @click="toggleSelection([tableData[1], tableData[2]])">全部上架</el-button>
      <el-button type="primary" size="mini" @click="toggleSelection()">全部删除</el-button>
    </div>
    <br>
    <el-table
      v-loading="listLoading"
      :data="list"
      element-loading-text="Loading"
      border
      fit
      highlight-current-row
    >
      <el-table-column type="selection" width="55" />
      <el-table-column type="index" :index="getIndex" label="序号" width="130" align="center" />
      <el-table-column label="姓名" width="130" align="center">
        <template slot-scope="scope">
          {{ scope.row.name }}
        </template>
      </el-table-column>
      <el-table-column label="性别" width="130" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.gender }}</span>
        </template>
      </el-table-column>
      <el-table-column label="民族" width="130" align="center">
        <template slot-scope="scope">
          {{ scope.row.nation }}
        </template>
      </el-table-column>
      <el-table-column label="出生地" width="130" align="center">
        <template slot-scope="scope">
          {{ scope.row.bornin }}
        </template>
      </el-table-column>
      <el-table-column label="出生时间" width="130" align="center">
        <template slot-scope="scope">
          {{ scope.row.bornat }}
        </template>
      </el-table-column>
      <el-table-column label="工作单位" align="center">
        <template slot-scope="scope">
          {{ scope.row.company }}
        </template>
      </el-table-column>
      <el-table-column class-name="status-col" label="操作" width="220" align="center">
        <template slot-scope="scope">
          <el-button
            type="primary"
            size="mini"
            @click="openDialog(scope.row)"
          >
            编辑
          </el-button>
          <el-button
            type="primary"
            size="mini"
            @click="openDialog(scope.row)"
          >
            上架
          </el-button>
          <el-button
            type="primary"
            size="mini"
            @click="openDialog(scope.row)"
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
  </div></template>

<script>
import { getAuthorList } from '@/api/author'
import queryform from '@/components/queryform'

export default {
  name: 'Management',
  components: {
    'el-queryform': queryform
  },
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

    return {
      list: null,
      listLoading: true,
      currentPage: 1, // 默认值
      pageSizes: pageSizes,
      pageSize: pageSize, // 默认值
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
      getAuthorList(params).then(response => {
        this.list = response.data.items
        console.log(this.list)
        this.listLoading = false
        this.totalPage = response.data.total
        this.currentPage = pageNum
      })
    },
    getIndex(index) {
      return (this.currentPage - 1) * this.pageSize + index + 1
    }
  }
}
</script>
