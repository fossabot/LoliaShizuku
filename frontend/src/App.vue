<script lang="ts" setup>
import AppHeader from "./components/AppHeader.vue";
</script>

<template>
  <v-app>
    <AppHeader style="--wails-draggable: drag" />
    <v-main>
      <router-view v-slot="{ Component, route }">
        <v-fade-transition :leave-absolute="true">
          <div :key="route.name" style="height: 100%">
            <component :is="Component" />
          </div>
        </v-fade-transition>
      </router-view>
    </v-main>
  </v-app>
</template>
