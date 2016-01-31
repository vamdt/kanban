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
      <legend>Color</legend>
      <div class="pure-control-group" v-for="dir in ['up', 'down', 'eq']">
        <label :for="dir">{{dir}}</label>
        <input type="color" placeholder="color" readonly
        v-model="settings.color[dir]" :id="dir">
        <colorpicker :for="dir"></colorpicker>
      </div>
    </fieldset>
    <fieldset>
      <legend>Mas</legend>
      <div class="pure-control-group" v-for="mas in settings.mas">
        <label v-bind="{for:'mas_'+mas.interval}">mas_{{mas.interval}}</label>
        <input v-bind="{id:'mas_'+mas.interval}" placeholder="mas" v-model="mas.interval" type="number">
        <input type="color" placeholder="color" readonly
        v-model="mas.color" :id="'cp_mas_'+mas.interval">
        <button class="pure-button" @click="del_mas(mas)">Delete</button>
        <colorpicker :for="'cp_mas_'+mas.interval"></colorpicker>
      </div>
      <div class="pure-control-group">
        <label></label>
        <button class="pure-button" @click="add_mas">Add</button>
        <button class="pure-button" @click="reset_mas">Reset</button>
      </div>
    </fieldset>

    <fieldset>
      <legend>Size</legend>
      <div class="pure-control-group">
        <label for="typing_circle_size">
        typing_circle_size
        </label>
        <input id="typing_circle_size" v-model="settings.typing_circle_size" placeholder="1" type="number">
      </div>
      <div class="pure-control-group">
        <label for="segment_circle_size">
        segment_circle_size
        </label>
        <input id="segment_circle_size" v-model="settings.segment_circle_size" placeholder="3" type="number">
      </div>
    </fieldset>
  </form>
</template>

<script lang="coffee">
config = require './config'
colorpicker = require './colorpicker.vue'
module.exports =
  components:
    colorpicker: colorpicker
  watch:
    'settings':
      handler: 'submit'
      deep: true
  data: ->
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
