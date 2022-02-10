<template>
  <form @submit="start">
    <label for="joinas">Join as...</label>
    <input type="text" placeholder="Hm?" id="joinas" v-model="as" required>
    <input type="submit" class="a" value="Join">
  </form>
</template>

<style scoped>
form {
  text-align: center;
}

label {
  display: block;
  font-weight: bold;
  font-size: 1.1em;
}

input[type="text"] {
  padding: .4em;
  border-radius: 6px;
  border: solid 2px #000;
  font-size: 1.3em;
}

label, input {
  margin: .3em;
  vertical-align: middle;
}
</style>

<script lang="ts">
export default {
  data() {
    return {
      err: false,
      as: ''
    }
  },
  methods: {
    start(e: Event) {
      e.preventDefault();
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
}
</script>