<template>
  <p v-if="err">Cannot load board!</p>
  <div id="board">
    <div class="row" :data-row="x" :key="x" v-for="(_, x) in 4">
      <div class="c" v-for="(s, i) in board.filter((_, e) => Math.floor(e / 10) == x)" :key="i" >
        <div>{{s.Name}}</div>
        <div>{{['Chance', 'Community Chest', 'IN JAIL', 'Free Parking', 'Go to Jail :)'].includes(s.Name) ? '' : '$'+s.Price}}</div>
        <div>
          <span class="player" v-for="(s, x) in dt.filter(v => v.Pos === (x*10)+i && Math.floor(v.Pos / 10) === x)" :key="x">{{s.Name}}</span>
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
    display: grid;
    grid-auto-columns: 8fr;
    grid-auto-rows: 9fr;
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

  .c {
    border: solid 3px #000;
    width: 100%;
    height: 100%;
    word-break: break-all;
    padding: 1em;
    overflow: hidden;
  }

  .player {
    padding: .2em .4em;
    border-radius: 30px;
    background-color: blue;
    color: #fff;
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
    })
  }
}
</script>