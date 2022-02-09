<template>
<h1>Hmm</h1>
<div class="s">
  <button>New Game</button>
  <button>Join Game</button>
</div>
</template>

<style scoped>
button:first-child {
  background-color: #D67021;
  box-shadow: 0 4px #ac5a1c;
}

button:last-child {
  background-color: #2C6E49;
  box-shadow: 0 4px #235739;
}

button {
  border-radius: 11px;
  padding: 1em 3em;
  width: 25vw;
  font-size: 1.2em;
  font-weight: bold;
  border: none;
  color: #fff;
  margin: 2em;
  min-width: 200px;
}

button:active {
  box-shadow: none;
  transform: translateY(4px);
}

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