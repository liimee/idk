<template>
  <title>New Game − Monopoly</title>
  <p v-if="err">an error occured</p>
</template>

<script lang="ts">
export default {
  data() {
    return {
      err: false
    }
  },
  mounted() {
      fetch((import.meta.env.VITE_API||'')+'/api/new', {
        method: 'POST'
      }).then(v => {
        if(v.ok) return v.json();
        this.err = true;
      }).then(v => {
        this.$router.replace("/g/"+v.id)
      })
  }
}
</script>
