<template>
  <form class="pure-form pure-form-aligned">
    <fieldset>
      <div class="pure-control-group">
        <label for="nc" class="pure-checkbox">
          <input id="nc" v-model="nc" @click="submit" type="checkbox"> no candle
        </label>
      </div>

      <div class="pure-control-group">
        <label for="nmas" class="pure-checkbox">
          <input id="nmas" v-model="nmas" @click="submit" type="checkbox"> no mas
        </label>
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
