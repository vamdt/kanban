<template>
  <div v-if="show" v-for="cs in color" class="pure-g">
    <div v-for="c in cs" class="pure-u-1-24"
    @click="sel(c)"
    :style="{
      backgroundColor: c,
      border: (input.value==c?3:0)+'px solid'
      }">
    {{c}}
    </div>
  </div>
</template>

<script lang="coffee">
d3 = require 'd3'

module.exports =
  props: ['for']
  data: ->
    issafari = false
    ua = navigator.userAgent.toLowerCase()
    if ua.indexOf('safari') != -1
      if ua.indexOf('chrome') > -1
      else
        issafari = true
    unless issafari
      return color: off, show: off
    @$nextTick =>
      input = @input = document.getElementById(@for)
      return unless @input
      @$nextTick => input.style.backgroundColor = input.value
      input.addEventListener 'focus', =>
        @show = on
      input.addEventListener 'blur', =>
        delay = => @show = off
        setTimeout delay, 200
    color = [[],[],[]]
    for c,j in [d3.scale.category20(), d3.scale.category20b(), d3.scale.category20c()]
      color[j].push c i for i in [0...20]
    color: color
    show: off

  methods:
    sel: (c) ->
      return unless @input
      @input.value = c
      @input.style.backgroundColor = c
      @input.dispatchEvent(new Event('change'))
      @show = off
</script>
