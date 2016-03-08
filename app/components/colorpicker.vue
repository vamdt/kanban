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
  <div v-if="show" class="pure-g">
    <button class="pure-button" @click="sel('')">Delete</button>
  </div>
</template>

<script>
import d3 from 'd3';

export default {
  props: ['for'],
  data() {
    const ua = navigator.userAgent.toLowerCase();
    const issafari = (ua.indexOf('safari') !== -1) && (ua.indexOf('chrome') === -1);
    if (!issafari) {
      return {
        color: false,
        show: false,
      };
    }
    this.$nextTick(() => {
      const input = this.input = document.getElementById(this.for);
      if (!this.input) {
        return;
      }
      this.$nextTick(() => {
        input.style.backgroundColor = input.value;
      });
      input.addEventListener('focus', () => {
        this.show = true;
      });
      input.addEventListener('blur', () => {
        setTimeout(() => {
          this.show = false;
        }, 200);
      });
    });
    const color = [[], [], []];
    const colors = [d3.scale.category20(), d3.scale.category20b(),
    d3.scale.category20c()];
    colors.forEach((c, j) => {
      for (let i = 0; i < 20; i++) {
        color[j].push(c(i));
      }
    });
    return {
      color,
      show: false,
    };
  },

  methods: {
    sel(c) {
      if (!this.input) {
        return;
      }
      this.input.value = c;
      this.input.style.backgroundColor = c;
      this.input.dispatchEvent(new Event('change'));
      this.show = false;
    },
  },
};
</script>
