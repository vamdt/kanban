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

<script lang="coffee">
d3 = require 'd3'
cmd = require './cmd.vue'
param = (hash, key) -> (hash=hash||{})[key]
module.exports =
  events:
    'show_stock': 'show_stock'
  components:
    cmd: cmd
  watch:
    cur_stock_name: (v) ->
      t = document.title.split('/')
      t[0] = v
      document.title = t.join('/')
  data: ->
    try
      stocks = JSON.parse localStorage.getItem 'stocks'
    catch

    stocks: stocks || []
    cur_stock_name: @stock_name param(@$route.params, 'sid'), stocks

  methods:
    tips: (msg, type) ->
      types = ['success', 'warning', 'danger']
      if -1 is types.indexOf type
        type = 'warning'
      @tip = msg: msg, className: 'tip-'+type

    stock_name: (sid, stocks) ->
      for s in stocks||@stocks||[]
        if s.sid == sid
          return s.name || sid
      sid

    lru: (s) ->
      stocks = @stocks || []
      i = -1
      i = j for ss, j in stocks when ss.sid == s.sid
      if i > -1
        stocks.splice(i, 1)
      stocks.unshift(s)
      localStorage.setItem('stocks', JSON.stringify(stocks))
      @stocks = stocks

    show_stock: (to) ->
      @cur_stock_name = to.name || to.sid
      @lru to
      k = param(@$route.params, 'k') || 1
      @$route.router.go
        name: 'stock'
        params: { sid: to.sid, k: k}
        replace: @$route.name is 'stock'
</script>
