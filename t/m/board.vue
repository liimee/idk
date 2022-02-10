<template>
  <p v-if="err">Cannot load board!</p>
  <table><tr v-for="i in 11" :key="i">
    <td v-for="(g, s) in board" :key="s">
      <div v-if="10-i == Math.floor(s / 11)">{{g.Name}}</div>
    </td>
  </tr></table>
</template>

<script lang="ts">
export default {
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