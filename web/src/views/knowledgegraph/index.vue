<template>
  <div class="kg-container">
    <div class="kg-text">专家知识图谱</div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'KnowledgeGraph',
  computed: {
    ...mapGetters([
      'name'
    ])
  }
}
</script>

<style lang="scss" scoped>
.kg {
  &-container {
    margin: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}
</style>
