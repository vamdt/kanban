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
Vue = require 'vue'
KLine = require '../stock/webpack'

Vue.directive 'kanpan',
  deep: true
  bind: ->
    window.addEventListener 'resize', =>
      @kl.resize() if @kl
    window.addEventListener 'keyup', (e) =>
      return unless e.target.tagName == 'BODY'
      return unless @kl
      handles =
        49: 'nmas'
        50: 'nc'
        51: 'nvolume'
      name = handles[e.keyCode] || off
      if name
        param = {}
        param[name] = not @kl.param name
        @kl.param param
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
      console.log 'has kl'
      if value.s is @kl.param 's'
        return @kl.param settings
      @kl.stop()
    else
      console.log 'new kl'
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
        v: +(new Date())
</script>
