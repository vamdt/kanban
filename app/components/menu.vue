<style>
.v-link-active {
  color: red;
}
[v-cloak] {
  display: none;
}
</style>

<template>
  <div class="pure-menu pure-menu-horizontal">
    <ul class="pure-menu-list">
        <li class="pure-menu-item pure-menu-selected"><a v-link="{ path: '/' }" class="pure-menu-link">Home</a></li>
        <li class="pure-menu-item pure-menu-has-children pure-menu-allow-hover">
            <a class="pure-menu-link">Stock</a>
            <ul class="pure-menu-children">
                <li class="pure-menu-item">
                  <form class="pure-form">
                    <input type="text" @keyup.enter="submit" v-model="sid"
                    placeholder="Search" class="pure-input-rounded"
                    autocomplete="off" autocorrect="off" autocapitalize="off"
                    spellcheck="false">
                  </form>
                </li>
                <li v-for="s in stocks" class="pure-menu-item">
                  <a v-link="{ path: '/s/'+s.sid+'/1' }"
                  class="pure-menu-link">{{s.name ? s.name : s.sid}}</a>
                </li>
            </ul>
        </li>
        <li class="pure-menu-item pure-menu-selected"><a v-link="{ path:
        '/settings' }" class="pure-menu-link">Settings</a></li>
    </ul>
  </div>
</template>

<script lang="coffee">
d3 = require 'd3'
module.exports =
  data: ->
    try
      stocks = JSON.parse localStorage.getItem 'stocks'
    catch

    stocks: stocks || []

  methods:
    lru: (s) ->
      stocks = @stocks || []
      i = stocks.indexOf s
      if i > -1
        stocks.splice(i, 1)
      stocks.unshift(s)
      localStorage.setItem('stocks', JSON.stringify(stocks))
      @stocks = stocks

    show_stock: (to) ->
      @lru to
      @$route.router.go(name: 'stock', params: { sid: to.sid, k: 1})

    submit: (event) ->
      sid = @sid
      event.preventDefault()
      @sid = ''
      for s in @stocks
        if s.sid == sid
          @show_stock(s)
          return

      d3.json '/search?s='+sid, (error, data) =>
        if error
          console.log error
          return
        if data.sid
          @show_stock data

</script>
