<template>
  <form class="pure-form" v-on:submit.prevent>
    <input type="text"
    @keyup.enter.prevent="run"
    @keyup.esc="cancle"
    v-model="cmd"
    id="cmd"
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

<script>
import Cmd from './cmd';
export default {
  props: {
    stocks: {
      type: Array,
      twoWay: true,
    },
  },
  data() {
    return {
      sugg: [],
    };
  },

  ready() {
    window.addEventListener('keyup', e => {
      if (e.target.tagName === 'INPUT') {
        return;
      }
      function dofocus() {
        e.preventDefault();
        document.getElementById('cmd').focus();
      }
      if (e.keyCode === 80 && e.ctrlKey && e.shiftKey) {
        dofocus();
        return;
      }
      if (e.keyCode === 186) {
        dofocus();
        return;
      }
    });
  },

  mixins: [Cmd],

  methods: {
    show_stock(to) {
      this.sugg = false;
      this.$dispatch('show_stock', to);
    },

    cancle() {
      this.sugg = false;
    },

    run() {
      let cmd = this.cmd;
      this.cmd = '';
      const opt = cmd.split(' ');
      if (opt.length < 1) {
        return;
      }
      if (opt.length < 2) {
        cmd = 'sugg';
      } else {
        cmd = opt.shift();
      }
      this.$emit(cmd, ...opt);
    },
  },
};
</script>
