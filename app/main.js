import 'purecss/build/pure-min.css';
import 'purecss/build/grids-responsive-min.css';
import 'purecss/build/buttons-min.css';
import Vue from 'vue';
import VueRouter from 'vue-router';
import App from './components/app.vue';
import Stock from './components/stock.vue';
import Settings from './components/settings.vue';
import Plate from './components/plate.vue';
import Lucky from './components/lucky.vue';

Vue.use(VueRouter);

const router = new VueRouter();

router.map({
  '/s/:sid/:k': {
    name: 'stock',
    component: Stock,
  },
  '/settings': {
    component: Settings,
  },
  '/plate/:pid/:id': {
    component: Plate,
  },
  '/lucky/:sid': {
    component: Lucky,
  },
});

router.redirect({
  '/lucky': '/lucky/sh000001',
  '/s/:sid': '/s/:sid/30',
});

router.start(App, '#app');
