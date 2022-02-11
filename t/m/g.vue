<script setup lang="ts">
  import Board from './board.vue';
</script>

<template>
  <div class="parent">
  <div class="board" style="">
    <Board :dt="da" />
  </div>
  <div style="padding: 1em">
    <div class="center">You currently have
      <div style="font-size: 1.6em; font-weight: bold">${{(da.find(v => v.Id === id)||{Money: 0}).Money}}</div>
    </div>
    <div v-if="!start">
      <button class="b" @click="st">Start Game</button>
    </div>
    <div v-if="mt">
      <button v-if="!rolled" class="a" @click="roll">Roll</button>
      <button v-if="rolled" class="b" @click="endTurn">End Turn</button>
      <button v-if="mt && rolled && da.every(v => !v.Owns.includes(da.find(v => v.Id === id).Pos))" class="a" @click="buy">Buy</button>
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
let s: WebSocket;
export default {
  data() {
    return {
      id: '',
      da: [],
      start: false,
      mt: false,
      rolled: false
    }
  },
  mounted() {
    const id = this.$route.params.id;
    s = new WebSocket((import.meta.env.VITE_API||'').replace(/https?/, 'ws')+'/api/ws')
    s.addEventListener('open', function () {
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
        console.log(ws.Data)
        break;
        case 'start':
        this.start = ws.Start
        break;
        case 'turn':
        this.mt = true;
      }
    })
  },
  unmounted() {
    s.close();
  },
  methods: {
    st() {
      s.send(JSON.stringify({
        s: 'start'
      }))
    },
    roll() {
      this.rolled = true;
      s.send(JSON.stringify({
        s: 'roll'
      }))
    },
    endTurn() {
      this.mt = false;
      this.rolled = false;
      s.send(JSON.stringify({
        s: 'endturn'
      }))
    },
    buy() {
      s.send(JSON.stringify({
        s: 'buy'
      }))
    }
  }
}
</script>