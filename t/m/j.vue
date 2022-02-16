<template>
  <form @submit="$router.replace('/g/'+a)">
    <label for="gameid">Game ID</label>
    <input type="text" id="gameid" v-model="a" required><input class="b" type="submit" value="Join">
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
</style>

<script>
export default {
  data() {
    return {
      a: ''
    }
  }
}
</script>
