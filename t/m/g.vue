<script setup lang="ts">
import Board from './board.vue';
</script>

<template>
  <title>Game {{$route.params.id}} − Monopoly</title>
  <div v-if="!str">
    <form @submit="join">
      <label for="joinas">Join as...</label>
      <input type="text" placeholder="Gopher" id="joinas" @input="e => as = e.target.value" required>
      <input type="submit" class="a" value="Join">
    </form>
  </div>
  <div class="parent" v-if="str && !res && !win">
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
        <button v-if="(rolled || da.find(v => v.Id === id).InJail) && da.find(v => v.Id === id).Money >= 0" class="b" @click="endTurn">End Turn</button>
        <button v-if="mt && rolled && da.every(v => !v.Owns.includes(da.find(v => v.Id === id).Pos)) && $refs.board.pay(da.find(v => v.Id === id).Pos) > 0 && $refs.board.pay(da.find(v => v.Id === id).Pos) <= da.find(v => v.Id === id).Money && !da.find(v => v.Id === id).InJail" class="a" @click="buy">Buy</button>
        <button v-if="!rolled && !da.find(v => v.Id === id).InJail" class="a" @click="roll">Roll</button>
        <button v-if="da.find(v => v.Id === id).InJail" class="b" @click="payJail">Pay $50 to get out of Jail</button>
        <button v-if="da.find(v => v.Id === id).Money < 0" class="a" @click="resign">Resign</button>
      </div>
      <div v-if="bid.bid">
        <hr />
        <div style="text-align: center">
          <p style="text-align: center">Auction <b>{{$refs.board.name(bid.bp)}}</b></p>
          <p>The highest bid is currently <b>${{bid.bidd}}</b></p>
        </div>
        <form @submit="sendBid" v-if="!bid.pa.includes(id)">
          <input type="number" v-model="bid.mybidd" :min="bid.bidd+1" step="1" :max="da.find(v => v.Id === id).Money">
          <input type="submit" value="OK" class="b">
        </form>
        <div style="text-align: center"><button class="a" :disabled="bid.pa.includes(id)" @click="pas">Pass</button></div>
        <ul v-if="bid.pa.length > 0">
          <b>Pass:</b>
          <li v-for="v, s in bid.pa" :key="s">{{da.find(g => g.Id === v).Name}}</li>
        </ul>
      </div>
      <div>
        <hr />
        <h2>Players ({{da.length}})</h2>
        <ul>
          <li v-for="g, i in da" :key="i">{{g.Name}}<span style="color: grey"> − has ${{g.Money}} − owns {{g.Owns.length}} properties {{g.InJail ? '− in jail':''}}</span></li>
        </ul>
      </div>
      <div>
        <hr />
        <h2>Your Properties ({{da.find(g => g.Id === id).Owns.length}})</h2>
        <p v-if="da.find(g => g.Id === id).Owns.length < 1">You don't have anything yet; go buy a property!</p>
        <ul>
          <li v-for="g, i in da.find(g => g.Id === id).Owns" :key="i">
            <b>{{$refs.board.name(g)}}</b>
            <a v-if="(self.Ho[g]||0) == 0" @click="mortgage(g)">[{{!self.Mo.includes(g) ? `mortgage ($${$refs.board.pay(g) / 2 })` : `unmortgage ($${($refs.board.pay(g) / 2) + ((10 / 100) * ($refs.board.pay(g) / 2))})` }}]</a>
            <a v-if="$refs.board.board[g].Set !== 0 && (self.Ho[g]||0) < 5 && !self.Mo.includes(g) && sameSet(g)" @click="ho(g)">[Buy House]</a>
            <a v-if="(self.Ho[g]||0) > 0" @click="sell(g)">[Sell House]</a>
          </li>
        </ul>
      </div>
      <div v-if="card.card">
        <hr />
        <div class="card">
          <div><b>{{card.t == 'a' ? 'Chance' : 'Community Chest'}}</b></div>
          {{card.str}}
        </div>
      </div>
      <div v-if="!start">
        <hr />
        To invite someone, you can share either this page's URL, or the Game ID:
        <div class="card l">{{$route.params.id}}</div>
      </div>
    </div>
  </div>
  <div v-if="err.e" class="er">{{err.m}}</div>
  <div v-if="res">
    <h1>You resigned</h1>
  </div>
  <div v-if="win" style="text-align: center">
    <h1>Congratulations!</h1>
    <p>You have won this game</p>
  </div>
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

li > a {
  text-decoration: underline;
  user-select: none;
  -webkit-user-select: none;
  cursor: pointer;
}

.card {
  border: solid 2px grey;
  border-radius: 6px;
  padding: 1em;
  width: 60%;
}

.card.l {
  width: max-content;
  padding: .5em;
  margin: .5em 0;
  font-weight: bold;
  letter-spacing: 2px;
}

@media (max-width: 750px) {
  .parent {
    grid-template-columns: 1fr;
    grid-template-rows: 60vh max-content;
    gap: 2em;
  }

  .board {
    width: calc(100% - 2.2em);
    height: 100%;
  }
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
        pa: [],
        bp: 0,
      },
      err: {
        e: false,
        m: ''
      },
      res: false,
      win: false,
      card: {
        card: false,
        str: '',
        t: 'a'
      }
    }
  },
  unmounted() {
    try {
      s.close();
    } catch(_) {}
  },
  methods: {
    join(e: Event) {
      e.preventDefault()

      fetch(import.meta.env.VITE_API+'/api/exists/'+this.$route.params.id)
      .then(v => v.json())
      .then(v => {
        if(!v.Exists) {
          this.$router.replace('/')
          alert('Game does not exist :(')
          return;
        }
      })

      this.str = true;

      s = new WebSocket((import.meta.env.VITE_API||'').replace(/https?/, 'wss')+'/api/ws')
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
          this.bid.bp = ws.Biddd
          break;
          case 'err':
          this.err = {
            e: true,
            m: ws.E
          }
          break;
          case 'win':
          this.win = true;
          break;
          case 'card':
          this.card.card = true;
          this.card.str = ws.Str;
          this.card.t = ws.T;
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
      this.card.card = false;
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
    },
    payJail() {
      s.send(JSON.stringify({
        s: 'payjail'
      }))
    },
    resign() {
      this.res = true;
      s.send(JSON.stringify({
        s: 'resign'
      }))
      s.close();
    },
    mortgage(m: number) {
      s.send(JSON.stringify({
        s: 'mortgage',
        pos: m
      }))
    },
    ho(m: number) {
      s.send(JSON.stringify({
        s: 'ho',
        pos: m
      }))
    },
    sell(m: number) {
      s.send(JSON.stringify({
        s: 'sell',
        pos: m
      }))
    },
    sameSet(m: number) {
      const l = this.$refs.board.board.filter(v => v.Set === this.$refs.board.board[m].Set).map(v => this.$refs.board.board.indexOf(v))
      const e = this.self.Owns.filter(v => l.includes(v))
      console.log(l, e)
      return l.length === e.length
    }
  },
  computed: {
    self() {
      return this.da.find(g => g.Id === this.id)
    }
  }
}
</script>
