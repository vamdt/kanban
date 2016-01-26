<template>
  <form class="pure-form pure-form-aligned">
    <fieldset>
      <legend>Candle</legend>
      <div class="pure-control-group">
        <label for="nc" class="pure-checkbox">
          <input id="nc" v-model="settings.nc" type="checkbox"> no candle
        </label>
      </div>

      <div class="pure-control-group">
        <label for="nmas" class="pure-checkbox">
          <input id="nmas" v-model="settings.nmas" type="checkbox"> no mas
        </label>
      </div>

      <div class="pure-control-group">
        <label for="ocl" class="pure-checkbox">
          <input id="ocl" v-model="settings.ocl" type="checkbox">
          only candle line
        </label>
      </div>

      <div class="pure-control-group">
        <label for="nvolume" class="pure-checkbox">
          <input id="nvolume" v-model="settings.nvolume" type="checkbox">
          no volume
        </label>
      </div>

      <div class="pure-control-group">
        <label for="nmacd" class="pure-checkbox">
          <input id="nmacd" v-model="settings.nmacd" type="checkbox">
          no macd
        </label>
      </div>
    </fieldset>
    <fieldset>
      <legend>Mas</legend>
      <div class="pure-control-group" v-for="mas in settings.mas">
        <label v-bind="{for:'mas_'+mas.interval}">mas_{{mas.interval}}</label>
        <input v-bind="{id:'mas_'+mas.interval}" placeholder="mas" v-model="mas.interval" type="number">
        <input type="color" v-on:click="cur_mas=mas.interval"
        v-model="mas.color" placeholder="color" readonly
        v-bind:style="{ backgroundColor: mas.color }">
        <button class="pure-button" @click="del_mas(mas)">Delete</button>
        <div v-if="color && cur_mas==mas.interval" v-for="cs in color" class="pure-g">
          <div v-for="c in cs" class="pure-u-1-24" @click="cur_mas=mas.color=c" v-bind:style="{ backgroundColor: c }">
          {{c}}
          </div>
        </div>
      </div>
      <div class="pure-control-group">
        <label></label>
        <button class="pure-button" @click="add_mas">Add</button>
        <button class="pure-button" @click="reset_mas">Reset</button>
      </div>
    </fieldset>
  </form>
</template>

<script lang="coffee">
d3 = require 'd3'
config = require './config'
module.exports =
  watch:
    'settings':
      handler: 'submit'
      deep: true
  data: ->
    issafari = false
    ua = navigator.userAgent.toLowerCase()
    if ua.indexOf('safari') != -1
      if ua.indexOf('chrome') > -1
      else
        issafari = true
    color = [[],[],[]]
    for c,j in [d3.scale.category20(), d3.scale.category20b(), d3.scale.category20c()]
      color[j].push c i for i in [0...20]
    unless issafari
      color = off
    color: color
    cur_mas: 0
    settings: {mas:[], fxck: off}
  route:
    data: ->
      @settings = config.load()

  methods:
    submit: (val, oldVal) ->
      return unless val and oldVal
      return if val.hasOwnProperty('fxck')
      return if oldVal.hasOwnProperty('fxck')
      config.save @settings
    add_mas: ->
      unless Array.isArray(@settings.mas)
        @settings.mas = []
      mas = @settings.mas
      n = interval: 5
      if mas.length
        n.interval = +mas[mas.length-1].interval + 1
      mas.push n
    del_mas: (mas) ->
      unless Array.isArray(@settings.mas)
        @settings.mas = []
      mass = @settings.mas
      index = mass.indexOf mas
      return if index is -1
      mass.splice(index, 1)
    reset_mas: ->
      mas = []
      for i in [5, 13, 21, 34, 55, 89, 144, 233]
        mas.push interval: i
      @settings.mas = mas
</script>
