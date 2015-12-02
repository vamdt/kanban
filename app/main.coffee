require('purecss/build/pure-min.css')
require('purecss/build/grids-responsive-min.css')
require('purecss/build/buttons-min.css')

Vue = require('vue')
VueRouter = require('vue-router')
Vue.use(VueRouter)

router = new VueRouter()

App = require('./components/app.vue')
Stock = require('./components/stock.vue')

router.map
  '/s/:sid/:k':
    component: Stock

router.start(App, '#app')
