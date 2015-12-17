<template>
  <form class="pure-form pure-form-aligned">
    <fieldset>
      <div class="pure-control-group">
        <label for="nc" class="pure-checkbox">
          <input id="nc" v-model="settings.nc" @click="submit" type="checkbox"> no candle
        </label>
      </div>

      <div class="pure-control-group">
        <label for="nmas" class="pure-checkbox">
          <input id="nmas" v-model="settings.nmas" @click="submit" type="checkbox"> no mas
        </label>
      </div>

      <div class="pure-control-group">
        <label for="ocl" class="pure-checkbox">
          <input id="ocl" v-model="settings.ocl" @click="submit" type="checkbox">
          only candle line
        </label>
      </div>
    </fieldset>
  </form>
</template>

<script lang="coffee">
module.exports =
  route:
    data: ->
      settings: try
          JSON.parse localStorage.getItem 'settings'
        catch
          nc: true
          nmas: true
          ocl: true

  methods:
    submit: ->
      localStorage.setItem('settings', JSON.stringify(@settings))
</script>
