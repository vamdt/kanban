<template>
<div class="pure-menu pure-menu-horizontal">
    <ul class="pure-menu-list">
        <li v-for="p in plate" class="pure-menu-item pure-menu-has-children pure-menu-allow-hover">
          <a v-link="{ path: '/plate/' + p.Id + '/0'}"
          class="pure-menu-link">{{p.Name}} {{p.Factor}}</a>
          <ul class="pure-menu-children">
              <li v-for="s in p.Sub" class="pure-menu-item">
                <a @click="show(s, $event)" v-link="{ path: '/plate/' + s.Pid + '/' + s.Id}"
                  class="pure-menu-link">{{s.Name}} {{s.Factor}}</a>
              </li>
          </ul>
        </li>
    </ul>
</div>
<div class="pure-g">
  <div v-for="i in stocks" class="pure-u-1-5">
    <a v-if="i.Leaf" v-link="{ path: '/s/' + i.Name + '/1'}" class="pure-menu-link">{{i.Name}} {{i.Factor}}</a>
  </div>
</div>
</template>

<script lang="coffee">
d3 = require 'd3'
param = (hash, key) -> (hash=hash||{})[key]
module.exports =
  data: ->
    plate: []
    stocks: []
  route:
    data: ->
      pid = param(@$route.params, 'pid') || 0
      id = param(@$route.params, 'id') || 0
      @rdata 0, =>
        @rdata pid, =>
          @rdata id, (s) =>
            @stocks = s if s

  methods:
    show: (plate, e) ->
      return unless plate.Sub
      return unless plate.Sub.length
      @stocks = plate.Sub

    rdata: (pid, cb) ->
      pid = +pid
      cb = cb || ->
      if pid == 0 and @plate.length > 0
        return cb()
      d3.json '/plate?pid=' + pid, (error, data) =>
        data = data.sort (a, b) -> b.Factor - a.Factor
        for d in data
          unless d.Sub
            d.Sub = []
        for p in @plate
          if +p.Id == pid
            p.Sub = data
            return cb(data)
          continue unless p.Sub
          for s in p.Sub
            if +s.Id == pid
              s.Sub = data
              return cb(data)
        @plate = data
        cb(data)
</script>
