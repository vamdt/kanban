<template>
  <div></div>
</template>

<script>
import lucky from './cmd/lucky';
export default {
  route: {
    data() {
      lucky.call(this);
    },
  },
};
</script>
