require('purecss/build/pure-min.css')
require('purecss/build/grids-responsive-min.css')
require('purecss/build/buttons-min.css')

Vue = require('vue')
VueRouter = require('vue-router')
Vue.use(VueRouter)

router = new VueRouter()

App = require('./components/app.vue')
Stock = require('./components/stock.vue')
Settings = require('./components/settings.vue')
Plate = require('./components/plate.vue')

router.map
  '/s/:sid/:k':
    name: 'stock'
    component: Stock
  '/settings':
    component: Settings
  '/plate':
    component: Plate
    subRoutes:
      '/:pid':
        component: Plate

router.start(App, '#app')
