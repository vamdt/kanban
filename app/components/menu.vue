<style>
[v-cloak] {
  display: none;
}
#menu li.pure-menu-selected a:hover,
#menu li.pure-menu-selected a:focus {
  background: #1f8dd6;
}
.tip-success {
  background: #dff0d8
}
.tip-warning {
  background: #fcf8e3
}
.tip-danger{
  background: #f2dede
}
</style>

<template>
  <div class="pure-menu pure-menu-horizontal">
    <ul class="pure-menu-list">
        <li class="pure-menu-item pure-menu-selected"><a v-link="{ path: '/' }" class="pure-menu-link">Home</a></li>
        <li class="pure-menu-item pure-menu-has-children pure-menu-allow-hover">
            <a v-link="{ path: '/s/'+cur_stock.sid+'/1' }" class="pure-menu-link">{{cur_stock.name}}</a>
            <ul class="pure-menu-children">
              <li v-for="s in stocks" class="pure-menu-item">
                <a v-link="{ path: '/s/'+s.sid+'/1' }"
                @click.prevent="show_stock(s)"
                class="pure-menu-link">{{s.name}}/{{s.sid}}</a>
              </li>
            </ul>
        </li>
        <li class="pure-menu-item pure-menu-selected">
          <a v-link="{ path: '/plate/0/0' }" class="pure-menu-link">Plate</a>
        </li>
        <li class="pure-menu-item pure-menu-selected"><a v-link="{ path:
        '/settings' }" class="pure-menu-link">Settings</a></li>
        <li class="pure-menu-item pure-menu-has-children pure-menu-allow-hover">
          <cmd :stocks.sync='stocks'></cmd>
        </li>
    </ul>
    <div class="pure-menu-item" v-bind:class="tip.className" v-if="tip&&tip.msg">
      {{tip.msg}}
    </div>
  </div>
</template>

<script>
import cmd from './cmd.vue';
function param(hash = {}, key) {
  return hash[key];
}

export default {
  events: {
    show_stock: 'show_stock',
    stock_change: 'stock_change',
  },
  components: {
    cmd,
  },

  data() {
    let stocks = [];
    try {
      stocks = JSON.parse(localStorage.getItem('stocks'));
    } catch (e) {
      stocks = [];
    }
    return {
      stocks: stocks || [],
      cur_stock: { name: 'Stock', sid: '' },
    };
  },

  methods: {
    tips(msg, type) {
      let t = type;
      const types = ['success', 'warning', 'danger'];
      if (types.indexOf(type) === -1) {
        t = 'warning';
      }
      this.tip = {
        msg,
        className: `tip-${t}`,
      };
    },

    lru(s) {
      const stocks = this.stocks || [];
      const i = stocks.findIndex((e) => e.sid === s.sid);
      if (i > -1) {
        stocks.splice(i, 1);
      }
      stocks.unshift(s);
      localStorage.setItem('stocks', JSON.stringify(stocks));
      this.stocks = stocks;
    },

    show_stock(to) {
      this.lru(to);
      const k = param(this.$route.params, 'k') || 1;
      this.$route.router.go({
        name: 'stock',
        params: { sid: to.sid, k },
        replace: this.$route.name === 'stock',
      });
    },

    stock_change(v) {
      this.cur_stock = v;
    },
  },
};
</script>
