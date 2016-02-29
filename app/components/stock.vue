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

<script lang="coffee">
require './kanpan'
module.exports =
  watch:
    'opt.k': (v) ->
      return unless v
      document.title = document.title.split('/')[0] + '/' + v
  data: ->
    levels: ['1', '5', '30', 'day', 'week', 'month']
    opt:{}
  route:
    data: (transition) ->
      @opt =
        s: @$route.params.sid
        k: @$route.params.k
        v: +(new Date())
</script>
