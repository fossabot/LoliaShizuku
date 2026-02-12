<script lang="ts" setup>
import AppHeader from "./components/AppHeader.vue";
import FloatingActionButton from "./components/FloatingActionButton.vue";
import { useGlobalLoading } from "@/composables/globalLoading";

const { isLoading } = useGlobalLoading();
</script>

<template>
  <v-app>
    <v-main>
      <AppHeader style="--wails-draggable: drag" />
      <div class="app-content-scroll">
        <v-progress-linear
          :active="isLoading"
          indeterminate
          color="primary"
          class="app-global-loading-bar"
        />
        <router-view v-slot="{ Component, route }">
          <transition name="fade" mode="out-in">
            <div :key="route.name" class="app-page-wrap">
              <v-container>
                <component :is="Component" />
              </v-container>
            </div>
          </transition>
        </router-view>
      </div>

      <!-- 全局悬浮按钮 -->
      <FloatingActionButton />
    </v-main>
  </v-app>
</template>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.1s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.app-content-scroll {
  height: calc(100vh - 64px);
  overflow-y: auto;
}

.app-page-wrap {
  width: 100%;
  height: 100%;
}

.app-global-loading-bar {
  position: fixed;
  left: 0;
  right: 0;
  top: 0;
  transform: translateY(64px);
  z-index: 2000;
  pointer-events: none;
}
</style>
