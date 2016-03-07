<style>
.view {
  transition: all .2s ease;
}
.test-enter, .test-leave {
  opacity: 0;
  transform: translate3d(10px, 0, 0);
}
.v-link-active {
  color: red;
}
[v-cloak] {
  display: none;
}
</style>

<template>
  <div>
    <menu></menu>
    <router-view class="view" transition="test" transition-mode="out-in" keep-alive></router-view>
  </div>
</template>

<script>
import menu from './menu.vue';
export default {
  components: {
    menu,
  },
};
</script>
