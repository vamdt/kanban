<style>
#level a.v-link-active {
  background: #fcebbd;
}
</style>

<template>
  <div id="level" class="pure-menu pure-menu-horizontal">
    <ul class="pure-menu-list">
        <li class="pure-menu-item" v-for="k in levels">
          <a class="pure-menu-link" v-link="{name:'stock', params: {sid:
          opt.s, k: k}, replace: true}">{{k}}</a>
        </li>
    </ul>
  </div>
  <div id="container" v-kanpan="opt"></div>
</template>

<script>
import config from './config';
import './kanpan';

export default {
  events: {
    param_change: 'param_change',
  },
  watch: {
    'opt.k': (v) => {
      if (!v) {
        return;
      }
      document.title = [document.title.split('/')[0], v].join('/');
    },
  },

  data() {
    return {
      levels: ['1', '5', '30', 'day', 'week', 'month'],
      opt: {},
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
  },
};
</script>
