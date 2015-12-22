<style>
#level li.pure-menu-selected {
  background: #fcebbd;
}
</style>

<template>
  <div id="level" class="pure-menu pure-menu-horizontal">
    <ul class="pure-menu-list">
        <li class="pure-menu-item" v-bind:class="{'pure-menu-selected':
        k==opt.k}" v-for="k in levels">
          <a class="pure-menu-link" @click="level(k)">{{k}}</a>
        </li>
    </ul>
  </div>
  <div id="container" v-kanpan="opt"></div>
</template>

<script lang="coffee">
Vue = require 'vue'
KLine = require '../stock'

Vue.directive 'kanpan',
  deep: true
  bind: ->
    window.addEventListener 'resize', =>
      @kl.resize() if @kl
  update: (value, oldValue) ->
    return unless value
    return unless value.s
    settings = {}
    try
      settings = JSON.parse localStorage.getItem 'settings'
    catch
    settings = settings || {}
    params = JSON.parse(JSON.stringify(value))
    for k,v of params
      settings[k] = v

    if @kl
      if value.s is @kl.param 's'
        return @kl.param settings
      @kl.stop()
    else
      @kl = new KLine(container: @el)
    kl = @kl
    kl.param settings

    setTimeout ->
      kl.init()
      kl.start()
    , 500

module.exports =
  data: ->
    levels: ['1', '5', '30', 'day', 'week', 'month']
    opt:{}
  route:
    data: (transition) ->
      @opt =
        s: @$route.params.sid
        k: @$route.params.k

  methods:
    level: (k) ->
      @opt.k = k

</script>
