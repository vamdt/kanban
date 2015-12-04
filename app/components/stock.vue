<template>
  <div id="container" v-kanpan="opt"></div>
</template>

<script lang="coffee">
Vue = require 'vue'
KLine = require '../stock'

Vue.directive 'kanpan',
  deep: true
  update: (value, oldValue) ->
    return unless value
    if @kl
      @kl.stop()
    else
      @kl = new KLine(container: @el)
    kl = @kl
    kl.param JSON.parse(JSON.stringify(value))

    setTimeout ->
      kl.init()
      kl.start()
    , 500
  unbind: ->
    console.log 'unbind'

module.exports =
  props:
    opt:
      s:
        type: String
        required: true
      k: Number
      fq:
        type: String
        default: ''
      nc:
        type: Number
        default: 1
      nmas:
        type: Number
        default: 1
  data: ->
    opt: @param()

  route:
    data: (transition) ->
      @opt = @param()

  methods:
    param: ->
      s: @$route.params.sid
      k: @$route.params.k
</script>
