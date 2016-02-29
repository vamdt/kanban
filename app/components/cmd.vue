<template>
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
</template>

<script lang="coffee">
d3 = require 'd3'
module.exports =
  props:
    stocks:
      type: Array
      twoWay: true
  data: ->
    sugg: []

  methods:
    show_stock: (to) ->
      @sugg = off
      @$dispatch 'show_stock', to

    cancle_sugg: ->
      @sugg = off
    do_sugg: ->
      sid = @sid
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
