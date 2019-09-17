import Vue from 'vue'
import Router from 'vue-router'
// import HelloWorld from '@/components/HelloWorld'
import tableList from '@/views/tableList'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'tableList',
      component: tableList
    }
  ]
})
