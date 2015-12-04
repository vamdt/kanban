<style>
.v-link-active {
  color: red;
}
[v-cloak] {
  display: none;
}
</style>

<template>
  <div class="pure-menu pure-menu-horizontal">
    <ul class="pure-menu-list">
        <li class="pure-menu-item pure-menu-selected"><a v-link="{ path: '/' }" class="pure-menu-link">Home</a></li>
        <li class="pure-menu-item pure-menu-has-children pure-menu-allow-hover">
            <a href="#" id="menuLink1" class="pure-menu-link">Stock</a>
            <ul class="pure-menu-children">
                <li class="pure-menu-item"><a v-link="{ path: '/s/sh600570/1' }" class="pure-menu-link">sh600570</a></li>
                <li class="pure-menu-item"><a v-link="{ path: '/s/sh601005/1' }" class="pure-menu-link">sh601005</a></li>
            </ul>
        </li>
        <li class="pure-menu-item pure-menu-selected"><a v-link="{ path:
        '/settings' }" class="pure-menu-link">Settings</a></li>
    </ul>
  </div>
</template>

<script lang="coffee">
module.exports = {}
</script>
