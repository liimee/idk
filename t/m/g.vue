<script setup lang="ts">
  import Board from './board.vue';
</script>

<template>
  <div class="board" style="">
    <Board />
  </div>
</template>

<style scoped>
  .board {
    width: 40em;
    overflow: auto;
    white-space: nowrap;
    padding: 1em;
    border-radius: 6px;
    border: solid 2px grey;
  }
</style>