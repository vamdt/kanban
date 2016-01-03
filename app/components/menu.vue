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
                class="pure-menu-link">{{s.name ? s.name : s.sid}}</a>
              </li>
            </ul>
        </li>
        <li class="pure-menu-item pure-menu-selected"><a v-link="{ path:
        '/settings' }" class="pure-menu-link">Settings</a></li>
        <li class="pure-menu-item pure-menu-has-children pure-menu-allow-hover">
          <form class="pure-form" v-on:submit.prevent>
            <input type="text" @keyup.enter.prevent="do_sugg" @keyup.esc="cancle_sugg" v-model="sid"
            placeholder="Search" class="pure-input-rounded"
            autocomplete="off" autocorrect="off" autocapitalize="off"
            spellcheck="false">
          </form>
          <ul v-show="sugg" class="pure-menu-children">
            <li v-for="s in sugg" class="pure-menu-item">
              <a v-link="{ path: '/s/'+s.sid+'/1' }"
              @click.prevent="show_stock(s)"
              class="pure-menu-link">{{s.name}}</a>
            </li>
          </ul>
        </li>
    </ul>
    <div class="pure-menu-item" v-bind:class="tip.className" v-if="tip&&tip.msg">
      {{tip.msg}}
    </div>
  </div>
</template>

<script lang="coffee">
d3 = require 'd3'
param = (hash, key) -> (hash||{})[key]
module.exports =
  data: ->
    try
      stocks = JSON.parse localStorage.getItem 'stocks'
    catch

    stocks: stocks || []
    cur_stock_name: @stock_name param(@$route.params, 'sid'), stocks
    sugg: []

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
      @sugg = off
      @cur_stock_name = to.name || to.sid
      @lru to
      @$route.router.go
        name: 'stock'
        params: { sid: to.sid, k: 1}
        replace: @$route.name is 'stock'

    cancle_sugg: ->
      @sugg = off
    do_sugg: (event) ->
      sid = @sid
      event.preventDefault()
      @sid = ''
      for s in @stocks
        if s.sid == sid
          @show_stock(s)
          return

      d3.text '/search?s='+sid, (error, data) =>
        if error
          console.log error
          return
        info = data.split(';')
        info.forEach (v, i) ->
          v = v.split(',')
          info[i] =
            sid: v[3]
            name: v[4]
        if info.length is 1
          return @show_stock(info[0])
        @sugg = info

</script>
