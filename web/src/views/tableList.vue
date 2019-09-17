<template>
  <div class="container">
    <el-card class="box-card">
      <el-row>
        <el-select v-model="select_value" placeholder="请选择">
          <el-option
            v-for="item in options"
            :key="item.value"
            :label="item.label"
            :value="item.value">
          </el-option>
        </el-select>
        &nbsp;&nbsp;
        <el-button type="primary" plain @click="AllConnect" :disabled="connect_but">全部连接</el-button>
        <el-button type="danger" plain @click="AllDisConnect" :disabled="connect_but">全部断开</el-button>
<!--        <el-button type="warning" plain @click="Pressure" :disabled="connect_but">开始压测</el-button>-->
        <el-button type="warning" class="button" plain @click="Pressure" :disabled="connect_but">
        {{ content }}
        </el-button>
      </el-row>
      <br><br>
      <el-table
        v-loading="loading"
        element-loading-text="拼命加载中"
        element-loading-spinner="el-icon-loading"
        element-loading-background="rgba(0, 0, 0, 0.8)"
        :data="tableData"
        style="width: 100%">
        <el-table-column
          label="IP"
          width="180">
          <template slot-scope="scope">
            <span style="margin-left: 10px">{{ scope.row.ip }}</span>
          </template>
        </el-table-column>
        <el-table-column
          label="状态"
          width="180">
          <template slot-scope="scope">
            <el-popover>
              <div slot="reference" class="name-wrapper">
                <el-tag size="medium" :type="scope.row.state !== 'idle'  ? 'primary' : 'danger'">{{ scope.row.state }}</el-tag>
              </div>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template slot-scope="scope">
            <el-button
              type="primary"
              size="mini"
              @click="handleEdit(scope.$index, scope.row)"
            :disabled="connect_but">连接
            </el-button>
            <el-button
              size="mini"
              type="danger"
              @click="handleDelete(scope.$index, scope.row)"
            :disabled="connect_but">断开
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'tableList',
  data () {
    return {
      timer: '',
      value: 0,
      loading: true,
      content: '开始压测',
      totalTime: 600,
      connect_but: false,
      tableData: [
        {
          ip: '192.168.100.2',
          state: 'cisco',
          value: '2'
        },
        {
          ip: '192.168.100.3',
          state: 'idle',
          value: '3'
        },
        {
          ip: '192.168.100.4',
          state: 'idle',
          value: '4'
        },
        {
          ip: '192.168.100.5',
          state: 'idle',
          value: '5'
        },
        {
          ip: '192.168.100.6',
          state: 'idle',
          value: '6'
        },
        {
          ip: '192.168.100.7',
          state: 'idle',
          value: '7'
        },
        {
          ip: '192.168.100.8',
          state: 'idle',
          value: '8'
        },
        {
          ip: '192.168.100.9',
          state: 'idle',
          value: '9'
        },
        {
          ip: '192.168.100.10',
          state: 'idle',
          value: '10'
        },
        {
          ip: '192.168.100.11',
          state: 'idle',
          value: '11'
        },
        {
          ip: '192.168.100.12',
          state: 'idle',
          value: '12'
        },
        {
          ip: '192.168.100.13',
          state: 'idle',
          value: '13'
        },
        {
          ip: '192.168.100.14',
          state: 'idle',
          value: '14'
        },
        {
          ip: '192.168.100.15',
          state: 'idle',
          value: '15'
        },
        {
          ip: '192.168.100.16',
          state: 'idle',
          value: '16'
        },
        {
          ip: '192.168.100.17',
          state: 'idle',
          value: '17'
        },
        {
          ip: '192.168.100.18',
          state: 'idle',
          value: '18'
        },
        {
          ip: '192.168.100.19',
          state: 'idle',
          value: '19'
        },
        {
          ip: '192.168.100.20',
          state: 'idle',
          value: '20'
        },
        {
          ip: '192.168.100.21',
          state: 'idle',
          value: '21'
        },
        {
          ip: '192.168.100.22',
          state: 'idle',
          value: '22'
        },
        {
          ip: '192.168.100.23',
          state: 'idle',
          value: '23'
        },
        {
          ip: '192.168.100.24',
          state: 'idle',
          value: '24'
        },
        {
          ip: '192.168.100.25',
          state: 'idle',
          value: '25'
        },
        {
          ip: '192.168.100.26',
          state: 'idle',
          value: '26'
        },
        {
          ip: '192.168.100.27',
          state: 'idle',
          value: '27'
        },
        {
          ip: '192.168.100.28',
          state: 'idle',
          value: '28'
        },
        {
          ip: '192.168.100.29',
          state: 'idle',
          value: '29'
        },
        {
          ip: '192.168.100.30',
          state: 'idle',
          value: '30'
        }
      ],
      options: [
        {
          value: 'cisco',
          label: 'cisco'
        }, {
          value: 'aruba',
          label: 'aruba'
        }, {
          value: 'arista',
          label: 'arista'
        }
      ],
      select_value: 'cisco'
    }
  },
  methods: {
    countDown () {
      this.connect_but = true
      this.content = '倒计时：' + this.totalTime + 's'
      let clock = window.setInterval(() => {
        this.totalTime--
        this.content = '倒计时：' + this.totalTime + 's'
        if (this.totalTime <= 0) {
          window.clearInterval(clock)
          this.content = '开始压测'
          this.totalTime = 600
          this.connect_but = false
        }
      }, 1000)
    },
    Pressure () {
      const that = this
      var state = false
      that.connect_but = true
      this.tableData.forEach(function (info, index) {
        if (info.state === 'idle') {
          console.log(info.ip + ':' + info.state)
          state = true
        }
      })
      if (state) {
        that.$message({
          message: '有服务器为idle状态，请先连接wifi，再进行压测！',
          type: 'error'
        })
        that.connect_but = false
        return false
      }
      this.tableData.forEach(function (info, index) {
        const url = '/trigger/' + info.value + '/'
        // console.log(url)
        console.log('pusser')
        axios.get(url)
          .then(response => {
            console.log(response)
          })
          .catch(function (error) {
            console.log('error')
            console.log(error)
          })
      })
      this.countDown()
    },
    AllConnect () {
      const that = this
      that.connect_but = true
      this.tableData.forEach(function (info, index) {
        const url = '/connect/' + info.value + '/?name=' + that.select_value
        axios.get(url)
          .then(response => {
            const data = response.data
            info.state = data
            // console.log(response)
            that.connect_but = false
          })
          .catch(function (error) {
            console.log('error')
            console.log(error)
            that.connect_but = false
          })
      })
    },
    AllDisConnect () {
      const that = this
      that.connect_but = true
      this.tableData.forEach(function (info) {
        const url = '/connect/' + info.value + '/?name=idle'
        axios.get(url)
          .then(response => {
            const data = response.data
            info.state = data
            that.connect_but = false
          })
          .catch(function (error) {
            console.log('error')
            console.log(error)
            that.connect_but = false
          })
      })
    },
    handleEdit (index, row) {
      console.log(index, row)
      this.connect_but = true
      const url = '/connect/' + row.value + '/?name=' + this.select_value
      console.log(url)
      const that = this
      axios.get(url)
        .then(response => {
          const data = response.data
          row.state = data
          console.log(response)
          that.connect_but = false
        })
        .catch(function (error) {
          console.log('error')
          console.log(error)
          that.connect_but = false
        })
    },
    handleDelete (index, row) {
      console.log(index, row)
      const that = this
      this.connect_but = true
      axios.get('/connect/' + row.value + '/?name=idle')
        .then(response => {
          const data = response.data
          row.state = data
          console.log(response)
          that.connect_but = false
        })
        .catch(function (error) {
          console.log('error')
          console.log(error)
          that.connect_but = false
        })
    },
    GetAllState () {
      this.tableData.forEach(function (info) {
        const url = '/state/' + info.value + '/'
        axios.get(url)
          .then(response => {
            const data = response.data
            info.state = data
            // console.log(response)
          })
          .catch(function (error) {
            console.log('error')
            console.log(error)
          })
      })
      window.setInterval(() => {
        this.loading = false
      }, 3000)
    },
    get () {
      this.tableData.forEach(function (info) {
        const url = '/state/' + info.value + '/'
        axios.get(url)
          .then(response => {
            const data = response.data
            info.state = data
          })
          .catch(function (error) {
            console.log('error')
            console.log(error)
          })
      })
      this.value = this.value + 1
      console.log(this.value)
    }
  },
  created () {
    this.GetAllState()
  },
  mounted () {
    this.timer = setInterval(this.get, 30000)
  }
}
</script>
<style>
  .container {
    width: 50%;
    text-align: center;
    padding-left: 20%;
  }
</style>
