<template>
  <div id="container"></div>
</template>

<script lang="coffee">
KLine = require '../stock'
module.exports =
  route:
    data: (transition) ->
      setTimeout(=>@kinit())

  methods:
    kinit: ->
      params = =>
        s: @$route.params.sid
        k: @$route.params.k
        fq: ''
        nc: 1
        nmas: 1

      if @kl
        console.log('stop first', @kl.param())
        @kl.stop()
      else
        @kl = new KLine(container: '#container')
      kl = @kl
      kl.param params()
      console.log('init', @kl.param())

      kl.init()
      kl.start()

  destroyed: ->
    console.log 'kl stop'
    if @kl
      @kl.stop()
</script>
