<style>
.view {
  transition: all .5s ease;
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
    <a v-link="{ path: '/s/sh600570/1' }">sh600570</a>
    <a v-link="{ path: '/s/sh601005/1' }">sh601005</a>
    <router-view class="view" transition="test" transition-mode="out-in" keep-alive></router-view>
  </div>
</template>

<script lang="coffee">
module.exports = {}
</script>
