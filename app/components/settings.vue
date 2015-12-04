<template>
  <form class="pure-form pure-form-aligned">
    <fieldset>
      <div class="pure-control-group">
        <label for="nc">candle</label>
        <input id="nc" type="checkbox"> no candle
      </div>

      <div class="pure-control-group">
        <label for="nmas">mas</label>
        <input id="nmas" type="checkbox"> no mas
      </div>

      <div class="pure-control-group">
        <label for="nmas">mas</label>
        <input id="nmas" type="checkbox"> no mas
      </div>

      <div class="pure-controls">
        <button type="submit" class="pure-button pure-button-primary">Submit</button>
      </div>
    </fieldset>
  </form>
</template>

<script lang="coffee">
module.exports =
  ready: ->
    localStorage.clear()
</script>
