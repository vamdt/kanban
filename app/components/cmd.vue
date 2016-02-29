<template>
  <form class="pure-form" v-on:submit.prevent>
    <input type="text"
    @keyup.enter.prevent="run"
    @keyup.esc="cancle"
    v-model="cmd"
    placeholder="Search" class="pure-input-rounded"
    autocomplete="off" autocorrect="off" autocapitalize="off"
    spellcheck="false">
  </form>
  <ul v-show="sugg" class="pure-menu-children">
    <li v-for="s in sugg" class="pure-menu-item">
      <a v-link="{ path: '/s/'+s.sid+'/1' }"
      @click.prevent="show_stock(s)"
      class="pure-menu-link">{{s.name}}</a>
    </li>
  </ul>
</template>

<script lang="coffee">
module.exports =
  props:
    stocks:
      type: Array
      twoWay: true
  data: ->
    sugg: []

  mixins: [require './cmd']

  methods:
    show_stock: (to) ->
      @sugg = off
      @$dispatch 'show_stock', to

    cancle: ->
      @sugg = off

    run: ->
      cmd = @cmd
      @cmd = ''
      opt = cmd.split ' '
      return if opt.length < 1
      if opt.length < 2
        cmd = 'sugg'
      else
        cmd = opt.shift()
      @$emit cmd, opt
</script>
