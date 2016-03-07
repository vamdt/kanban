<style>
.v-link-active {
  color: red;
}
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
            <a class="pure-menu-link">{{cur_stock_name ? cur_stock_name : 'Stock'}}</a>
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
  },
  components: {
    cmd,
  },
  watch: {
    cur_stock_name: v => {
      const t = document.title.split('/');
      t[0] = v;
      document.title = t.join('/');
    },
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
      cur_stock_name: this.stock_name(param(this.$route.params, 'sid'), stocks),
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

    stock_name(sid, stocks) {
      const all = stocks || this.stocks || [];
      for (const s in all) {
        if (s.sid === sid) {
          return s.name || sid;
        }
      }
      return sid;
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
      this.cur_stock_name = to.name || to.sid;
      this.lru(to);
      const k = param(this.$route.params, 'k') || 1;
      this.$route.router.go({
        name: 'stock',
        params: { sid: to.sid, k },
        replace: this.$route.name === 'stock',
      });
    },
  },
};
</script>
