<style>
#level a.v-link-active {
  background: #fcebbd;
}
</style>

<template>
<div>
  <div id="level" class="pure-menu pure-menu-horizontal">
    <ul class="pure-menu-list">
        <li class="pure-menu-item" v-for="k in levels">
          <a class="pure-menu-link" v-link="{name:'stock', params: {sid:
          opt.s, k: k}, replace: true}">{{k}}</a>
        </li>
        <li v-if="starred" class="pure-menu-item">
          <button class="pure-button" @click="unstar(opt.s)">Unstar</button>
        </li>
        <li v-else class="pure-menu-item">
          <button class="pure-button" @click="star(opt.s)">Star</button>
        </li>
        <li class="pure-menu-item">
          <a v-link="{ path: '/lucky/'+opt.s }" class="pure-menu-link">Lucky</a>
        </li>
    </ul>
  </div>
  <div id="container" v-kanpan="opt"></div>
</div>
</template>

<script>
import config from './config';
import { default as star, isStar } from './cmd/star';
import unstar from './cmd/unstar';
import './kanpan';

export default {
  events: {
    param_change: 'param_change',
  },
  watch: {
    'opt.s': 'checkStar',
    sname: 'updateTitle',
    'opt.k': 'updateTitle',
  },

  data() {
    return {
      levels: ['1', '5', '30', 'day', 'week', 'month'],
      opt: {},
      sname: '',
      starred: false,
    };
  },
  route: {
    data() {
      this.opt = {
        s: this.$route.params.sid,
        k: this.$route.params.k,
        v: +(new Date()),
      };
    },
  },
  methods: {
    param_change(opts) {
      config.update(opts);
      for (const k in opts) {
        if (opts.hasOwnProperty(k)) {
          this.opt[k] = opts[k];
        }
      }
      this.opt.v = +(new Date());
    },
    star,
    unstar,
    checkStar(v) {
      if (!v || v.length < 1) {
        this.starred = false;
        return;
      }
      isStar(v, (e, res) => {
        this.starred = res && res.star;
      });
    },
    updateTitle() {
      document.title = [this.sname, this.opt.s, this.opt.k].join('/');
      this.$root.$broadcast('stock_change', {
        sid: this.opt.s,
        name: this.sname,
      });
    },
  },
};
</script>
