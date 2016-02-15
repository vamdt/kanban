<template>
<div class="pure-menu pure-menu-horizontal pure-menu-scrollable">
    <a href="javascript:{}" class="pure-menu-link pure-menu-heading">Plate</a>
    <ul class="pure-menu-list">
        <li v-for="p in plate0" class="pure-menu-item">
          <a v-link="{ path: '/plate/'+p.Id }"
          class="pure-menu-link">{{p.Name}} {{p.Factor}}</a>
        </li>
    </ul>
</div>
<router-view></router-view>
</template>

<script lang="coffee">
d3 = require 'd3'
param = (hash, key) -> (hash=hash||{})[key]
module.exports =
  data: ->
    plate0:[]
    plate1:[]
  route:
    data: ->
      d3.json '/plate?pid=0' + pid, (error, data) =>
        @plate0 = data
      pid = +param(@$route.params, 'pid') || 0
      if pid > 0
        d3.json '/plate?pid=' + pid, (error, data) =>
          @plate1 = data
</script>
