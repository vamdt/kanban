<template>
  <form class="pure-form pure-form-aligned">
    <fieldset>
      <div class="pure-control-group">
        <label for="nc">candle</label>
        <input id="nc" v-model="nc" type="checkbox"> no candle
      </div>

      <div class="pure-control-group">
        <label for="nmas">mas</label>
        <input id="nmas" v-model="nmas" type="checkbox"> no mas
      </div>

      <div class="pure-controls">
        <button type="submit" @click="submit" class="pure-button pure-button-primary">Submit</button>
      </div>
    </fieldset>
  </form>
</template>

<script lang="coffee">
module.exports =
  route:
    data: ->
      try
        JSON.parse localStorage.getItem 'settings'
      catch
        nc: true
        nmas: true

  methods:
    submit: ->
      settings =
        nc: @nc
        nmas: @nmas
      localStorage.setItem('settings', JSON.stringify(settings))
</script>
