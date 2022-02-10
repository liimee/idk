<script setup lang="ts">
  import Board from './board.vue';
</script>

<template>
  <div class="parent">
  <div class="board" style="">
    <Board />
  </div>
  <div style="padding: 1em">
    <div class="center">You currently have
      <div style="font-size: 1.6em; font-weight: bold">${{(da.find(v => v.Id === id)||{Money: 0}).Money}}</div>
    </div>
  </div>
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

  .parent {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
  }

  .center {
    text-align: center;
  }
</style>

<script lang="ts">
export default {
  data() {
    return {
      id: '',
      da: []
    }
  },
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
      const ws = JSON.parse(s.data)
      switch(ws.S) {
        case 'id':
        this.id = ws.Id
        break;
        case 'data':
        this.da = ws.Data
      }
    })
  }
}
</script>