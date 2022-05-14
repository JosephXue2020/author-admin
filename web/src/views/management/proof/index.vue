<template>
  <div class="management-container">
    <div class="management-text">专家审核</div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'Proof',
  computed: {
    ...mapGetters([
      'name'
    ])
  }
}
</script>

<style lang="scss" scoped>
.management {
  &-container {
    margin: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}
</style>
