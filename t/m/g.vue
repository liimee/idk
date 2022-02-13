<script setup lang="ts">
  import Board from './board.vue';
</script>

<template>
  <div v-if="!str">
    <form @submit="join">
    <label for="joinas">Join as...</label>
    <input type="text" placeholder="Hm?" id="joinas" v-model="as" required>
    <input type="submit" class="a" value="Join">
  </form>
  </div>
  <div class="parent" v-if="str">
  <div class="board" style="">
    <Board :dt="da" ref="board" />
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
      <button v-if="mt && rolled && da.every(v => !v.Owns.includes(da.find(v => v.Id === id).Pos)) && $refs.board.pay(da.find(v => v.Id === id).Pos) > 0 && $refs.board.pay(da.find(v => v.Id === id).Pos) <= da.find(v => v.Id === id).Money" class="a" @click="buy">Buy</button>
    </div>
    <div v-if="bid.bid">
      The highest bid is currently <b>${{bid.bidd}}</b>
      <form @submit="sendBid" v-if="!bid.pa.includes(id)">
        <input type="number" v-model="bid.mybidd" :min="bid.bidd+1" step="1" :max="da.find(v => v.Id === id).Money">
        <input type="submit" value="OK" class="b">
      </form>
      <button class="a" :disabled="bid.pa.includes(id)" @click="pas">Pass</button>
      <ul v-if="bid.pa.length > 0">
        Pass:
        <li v-for="v, s in bid.pa" :key="s">{{da.find(g => g.Id === v).Name}}</li>
      </ul>
    </div>
  </div>
  </div>
  <div v-if="err.e" class="er">{{err.m}}</div>
</template>

<style scoped>
  .board {
    width: 45vw;
    height: calc(100vh - 4em);
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

  form {
  text-align: center;
}

label {
  display: block;
  font-weight: bold;
  font-size: 1.1em;
}

input[type="text"],
input[type="number"] {
  padding: .4em;
  border-radius: 6px;
  border: solid 2px #000;
  font-size: 1.3em;
}

label, input {
  margin: .3em;
  vertical-align: middle;
}

button:disabled {
  background-color: rgb(97, 97, 97);
  box-shadow: 0 4px rgb(44, 44, 44);
}

.er {
  background-color: rgb(255, 182, 182);
  color: #f83c48;
  padding: 1em;
  border-radius: 6px;
  margin: 1em;
  font-weight: bold;
}
</style>

<script lang="ts">
let s: WebSocket;
export default {
  data() {
    return {
      as: '',
      str: false,
      id: '',
      da: [],
      start: false,
      mt: false,
      rolled: false,
      bid: {
        bid: false,
        bidd: 0,
        mybidd: 0,
        pa: []
      },
      err: {
        e: false,
        m: ''
      }
    }
  },
  unmounted() {
    s.close();
  },
  methods: {
    join(e: Event) {
      e.preventDefault()
      this.str = true;
      s = new WebSocket((import.meta.env.VITE_API||'').replace(/https?/, 'ws')+'/api/ws')
      s.addEventListener('open', () => {
        s.send(JSON.stringify({
          s: 'join',
          as: this.as,
          id: this.$route.params.id
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
        if(this.mt) {
          // @ts-ignore
          this.$refs.board.scr(this.da.find(v => v.Id === this.id).Pos)
        }
        break;
        case 'start':
        this.start = ws.Start
        break;
        case 'turn':
        this.mt = true;
        break;
        case 'bid':
        this.bid.pa = ws.Pa
        this.bid.bidd = ws.Bid
        this.bid.bid = ws.Biddi
        break;
        case 'err':
        this.err = {
          e: true,
          m: ws.E
        }
      }
    })
    },
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
    },
    sendBid(e: Event) {
      e.preventDefault();
      s.send(JSON.stringify({
        s: "bid",
        bid: this.bid.mybidd
      }))
      this.bid.mybidd = 0;
    },
    pas() {
      s.send(JSON.stringify({
        s: 'pass'
      }))
    }
  }
}
</script>