<script lang="ts" setup>
import { ref, onMounted, onUnmounted } from "vue";
import {
  WindowMinimise,
  WindowToggleMaximise,
  WindowIsMaximised,
  EventsOn,
  EventsOff,
  Quit,
} from "../../wailsjs/runtime/runtime";
import AppLogo from "./AppLogo.vue";
import { useTheme } from "vuetify";

const maximised = ref(false);
const theme = useTheme();

const onToggleMaximize = (isMaximised: boolean) => {
  maximised.value = isMaximised;
};

onMounted(async () => {
  const isMax = await WindowIsMaximised();
  onToggleMaximize(isMax);

  EventsOn(
    "window_changed",
    (info: { fullscreen?: boolean; maximised?: boolean }) => {
      const { maximised: isMaximised } = info;
      if (isMaximised !== undefined) {
        onToggleMaximize(isMaximised);
      }
    },
  );
});

onUnmounted(() => {
  EventsOff("window_changed");
});

async function handleMinimize() {
  WindowMinimise();
}

async function handleMaximize() {
  WindowToggleMaximise();
}

function handleClose() {
  Quit();
}

function toggleTheme() {
  theme.global.name.value = theme.global.current.value.dark
    ? "lightTheme"
    : "darkTheme";
}
</script>

<template>
  <v-app-bar color="appbar" app>
    <template #prepend>
      <div
        style="display: flex; align-items: center; gap: 8px; margin-left: 17px"
      >
        <AppLogo :size="20" />
        <v-app-bar-title class="font-comfortaa text-lg font-semibold">
          LoliaShizuku
        </v-app-bar-title>
      </div>
    </template>

    <template #append>
      <!-- 主题切换 -->
      <v-btn
        :icon="theme.global.current.value.dark ? 'fas fa-sun' : 'fas fa-moon'"
        variant="text"
        size="small"
        @click="toggleTheme"
      />

      <!-- 窗口控制按钮 -->
      <v-btn
        icon="fas fa-minus"
        variant="text"
        size="small"
        @click="handleMinimize"
      />

      <v-btn
        :icon="maximised ? 'fas fa-window-restore' : 'fas fa-window-maximize'"
        variant="text"
        size="small"
        @click="handleMaximize"
      />

      <v-btn
        icon="fas fa-times"
        variant="text"
        size="small"
        @click="handleClose"
        class="mr-2"
      />
    </template>
  </v-app-bar>
</template>
