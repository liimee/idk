<template>
  <p v-if="err">Cannot load board!</p>
  <div id="board" ref="s">
    <div id="center"><h1>Monopoly</h1></div>
    <div class="row" :data-row="x" :key="x" v-for="(_, x) in 4">
      <div class="c" :style="{background: `linear-gradient(to bottom, ${['transparent', '#CB8556', '#2AC3D1', '#FF0088', '#E88720', '#CB0000', '#F6BE16', '#00C750', '#0047CB'][board[(x*10)+i].Set]}, transparent 40%)`}" v-for="(s, i) in board.filter((_, e) => Math.floor(e / 10) == x)" :key="i" >
        <div><b>{{s.Name}}</b></div>
        <div>{{['Chance', 'Community Chest', 'IN JAIL', 'Free Parking', 'Go to Jail :)'].includes(s.Name) ? '' : '$'+s.Price}}</div>
        <div style="color: grey" v-if="dt.some(v => v.Owns.includes((x*10)+i))">({{dt.find(v => v.Owns.includes((x*10)+i)).Name}})</div>
        <div v-if="dt.some(v => v.Ho[(x*10)+i])">
          <span v-if="dt.find(v => v.Ho[(x*10)+i]).Ho[(x*10)+i] < 5" v-for="_ in dt.find(v => v.Ho[(x*10)+i]).Ho[(x*10)+i]">🏠</span>
          <span v-else>🏨</span>
        </div>
        <div>
          <span :style="{backgroundColor: s.Color}" class="player" v-for="(s, x) in dt.filter(v => v.Pos === (x*10)+i && Math.floor(v.Pos / 10) === x)" :key="x">{{s.Name}}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  * {
    box-sizing: border-box;
  }

  #board {
    background-color: #FFF9DC;
    display: grid;
    grid-auto-columns: 390px;
    grid-auto-rows: 180px;
    width: max-content;
    grid-template-areas:
      "a a a a a a a a a a a d"
      "b o o o o o o o o o o d"
      "b o o o o o o o o o o d"
      "b o o o o o o o o o o d"
      "b o o o o o o o o o o d"
      "b o o o o o o o o o o d"
      "b o o o o o o o o o o d"
      "b o o o o o o o o o o d"
      "b o o o o o o o o o o d"
      "b o o o o o o o o o o d"
      "b c c c c c c c c c c c";
  }

  .row {
    display: flex;
  }

  .row[data-row="0"] {
   grid-area: c;
   flex-direction: row-reverse;
  }

  .row[data-row="1"] {
   grid-area: b;
   flex-direction: column-reverse;
  }

  .row[data-row="2"] {
   grid-area: a;
   flex-direction: row;
  }

  .row[data-row="3"] {
   grid-area: d;
   flex-direction: column;
  }

  #center {
    grid-area: o;
    font-size: 5em;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .c {
    border: solid 3px #000;
    width: 100%;
    height: 100%;
    word-wrap: break-word;
    padding: 1em;
    text-align: center;
  }

  .player {
    padding: .2em 1em;
    border-radius: 30px;
    background-color: blue;
    color: #fff;
    margin: 0 .2em;
  }

  .c > div {
    margin: .5em;
  }
</style>

<script lang="ts">
export default {
  props: ['dt'],
  data() {
    return {
      board: [],
      err: false
    }
  },
  mounted() {
    // @ts-ignore
    fetch((import.meta.env.VITE_API||'')+'/api/board')
    .then(v => {
      if(v.ok) return v.json();
      this.err = true;
    }).then(v => {
      this.board = v
    }).then(() => {
      this.scr(0)
    })
  },
  methods: {
    scr(g: number) {
      (this.$refs.s as HTMLElement).children[Math.floor(g / 10) + 1].children[g-(Math.floor(g / 10)*10)].scrollIntoView()
    },
    pay(g: number) {
      return this.board[g].Price
    },
    name(g: number) {
      return this.board[g].Name
    }
  }
}
</script>
