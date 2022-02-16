<template>
  <title>Monopoly</title>
  <h1>Monopoly</h1>
  <div class="s">
    <router-link to="/n"><button class="a">New Game</button></router-link>
    <router-link to="/j"><button class="b">Join Game</button></router-link>
  </div>
</template>

<style scoped>
.s {
  text-align: center;
}

@media (max-width: 290px) {
  button {
    margin: 2em .5em;
    width: calc(100vw - 2em);
    min-width: 0;
  }
}
</style>
