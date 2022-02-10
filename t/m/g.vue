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
    width: 45vw;
    overflow: auto;
    white-space: nowrap;
    padding: 1em;
    border-radius: 6px;
    border: solid 2px grey;
  }
</style>

<script lang="ts">
export default {
  mounted() {
    const id = this.$route.params.id;
    const s = new WebSocket((import.meta.env.VITE_API||'').replace(/https?/, 'ws')+'/api/ws')
    s.addEventListener('open', function (event) {
      s.send(JSON.stringify({
        s: 'join',
        as: 'name',
        id: id
      }));
    });

    s.addEventListener('message', s => {
      console.log(s.data)
    })
  }
}
</script>