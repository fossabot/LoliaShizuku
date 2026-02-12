<script lang="ts" setup>
import { onBeforeUnmount, onMounted, computed, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { getRunnerRuntimeStatus } from "@/services/center";

const route = useRoute();
const runnerRunning = ref(false);
let timer: ReturnType<typeof setInterval> | null = null;

const refreshRunnerStatus = async () => {
  try {
    const status = await getRunnerRuntimeStatus();
    runnerRunning.value = !!status.running;
  } catch {
    runnerRunning.value = false;
  }
};

const showFab = computed(
  () => runnerRunning.value && route.path !== "/oauth" && route.path !== "/runner",
);

onMounted(() => {
  void refreshRunnerStatus();
  timer = setInterval(() => {
    void refreshRunnerStatus();
  }, 3000);
});

onBeforeUnmount(() => {
  if (timer) {
    clearInterval(timer);
    timer = null;
  }
});

watch(
  () => route.path,
  () => {
    void refreshRunnerStatus();
  },
);
</script>

<template>
  <v-fab
    v-if="showFab"
    class="fab-global"
    color="primary"
    size="large"
    location="bottom end"
    prepend-icon="fas fa-terminal"
    to="/runner"
  >
    回到运行
  </v-fab>
</template>

<style scoped>
.fab-global {
  position: fixed !important;
  bottom: 24px !important;
  right: 24px !important;
  z-index: 1000;
}
</style>
