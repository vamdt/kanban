<template>
  <form class="pure-form" v-on:submit.prevent>
    <input type="text"
    @keyup.enter.prevent="run"
    @keyup.esc="cancle"
    v-model="cmd"
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

  ready: ->
    events = [
      'sugg'
      'unwatch'
      'watch'
    ]
    for e in events
      @$on e, 'do_' + e

  methods:
    show_stock: (to) ->
      @sugg = off
      @$dispatch 'show_stock', to

    cancle: ->
      @sugg = off

    do_watch: (opt) ->
      unless Array.isArray opt
        opt = [opt]

      num = 0
      for o in opt
        s = sid:o, name:o
        i = -1
        i = j for ss, j in @stocks when ss.sid == o
        if i > -1
          s.name = ss.name
          @stocks.splice(i, 1)
        @stocks.unshift(s)
        num++
      if num
        localStorage.setItem('stocks', JSON.stringify(@stocks))

    do_unwatch: (opt) ->
      unless Array.isArray opt
        opt = [opt]

      num = 0
      for o in opt
        i = -1
        i = j for ss, j in @stocks when ss.sid == o
        if i > -1
          @stocks.splice(i, 1)
          num++
      if num
        localStorage.setItem('stocks', JSON.stringify(@stocks))

    do_sugg: (sid) ->
      if sid.length < 1
        return
      if Array.isArray sid
        sid = sid[0]
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

    run: ->
      cmd = @cmd
      @cmd = ''
      opt = cmd.split ' '
      return if opt.length < 1
      if opt.length < 2
        cmd = 'sugg'
      else
        cmd = opt.shift()
      @$emit cmd, opt

</script>
