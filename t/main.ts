import {createApp} from 'vue';
import {createRouter, createWebHistory} from 'vue-router';

import App from './m/l.vue';
import New from './m/n.vue';
import Game from './m/g.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [{
    path: '/', component: App
  }, {
    path: '/n',
    component: New
  }, {
    path: '/g/:id',
    component: Game
  }]
})

createApp({}).use(router).mount('#app')