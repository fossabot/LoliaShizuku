<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { storeToRefs } from "pinia";
import { useGlobalLoadingStore } from "@/stores/globalLoading";
import { useFrpcInstallStore } from "@/stores/frpcInstall";
import {
  getFrpcStatus,
  removeFrpc,
  setGitHubMirrorURL,
  type FrpcStatus,
} from "@/services/frpc";

defineOptions({
  name: "SettingsPage",
});

const status = ref<FrpcStatus | null>(null);
const activePanel = ref<"frpc">("frpc");
const snackbar = ref(false);
const snackbarText = ref("");
const snackbarColor = ref<"success" | "error" | "info">("info");
const githubMirrorURLInput = ref("");
const builtinMirrorURL = "https://cdn.akaere.online/github.com";
type MirrorMode = "official" | "builtin" | "custom";
const mirrorMode = ref<MirrorMode>("official");
const mirrorModeItems = [
  { title: "github.com", value: "official" as const },
  { title: "cdn.akaere.online/github.com", value: "builtin" as const },
  { title: "自定义网址", value: "custom" as const },
];

const globalLoadingStore = useGlobalLoadingStore();
const withGlobalLoading = <T>(task: () => Promise<T>) =>
  globalLoadingStore.withGlobalLoading(task);
const frpcInstallStore = useFrpcInstallStore();
const { installing, canceling } = storeToRefs(frpcInstallStore);
const { startInstall, cancelInstall } = frpcInstallStore;

const showMessage = (
  text: string,
  color: "success" | "error" | "info" = "info",
) => {
  snackbarText.value = text;
  snackbarColor.value = color;
  snackbar.value = true;
};

const installedVersion = computed(
  () => status.value?.installed?.version || "未安装",
);
const latestVersion = computed(() => status.value?.latest?.tag_name || "-");
const actionText = computed(() => {
  if (!status.value?.installed?.binary_exists) {
    return "安装 frpc";
  }
  if (status.value.update_available) {
    return "更新 frpc";
  }
  return "重装 frpc";
});

const formatTime = (value?: string) => {
  if (!value) {
    return "-";
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return date.toLocaleString();
};

const loadStatus = async () => {
  await withGlobalLoading(async () => {
    try {
      status.value = await getFrpcStatus();
      const mirrorURL = (status.value.github_mirror_url || "").trim();
      if (mirrorURL === "") {
        mirrorMode.value = "official";
        githubMirrorURLInput.value = "";
      } else if (mirrorURL === builtinMirrorURL) {
        mirrorMode.value = "builtin";
        githubMirrorURLInput.value = "";
      } else {
        mirrorMode.value = "custom";
        githubMirrorURLInput.value = mirrorURL;
      }
    } catch (error) {
      showMessage(
        error instanceof Error ? error.message : "获取 frpc 状态失败",
        "error",
      );
    }
  });
};

const handleInstallOrUpdate = async () => {
  try {
    await withGlobalLoading(async () => {
      const result = await startInstall();
      status.value = result.status;
      showMessage(`frpc 已安装到 ${result.status.paths.binary_path}`, "success");
    });
  } catch (error) {
    const message = error instanceof Error ? error.message : "安装/更新 frpc 失败";
    if (message.includes("已终止")) {
      showMessage(message, "info");
      await loadStatus();
      return;
    }
    showMessage(message, "error");
  }
};

const handleCancelInstall = async () => {
  if (!installing.value || canceling.value) {
    return;
  }

  try {
    await cancelInstall();
    showMessage("已发送终止下载请求", "info");
  } catch (error) {
    showMessage(error instanceof Error ? error.message : "终止下载失败", "error");
  }
};

const handleRemove = async () => {
  await withGlobalLoading(async () => {
    try {
      await removeFrpc();
      showMessage("本地 frpc 已移除", "success");
      await loadStatus();
    } catch (error) {
      showMessage(error instanceof Error ? error.message : "移除 frpc 失败", "error");
    }
  });
};

const handleSaveMirrorURL = async () => {
  await withGlobalLoading(async () => {
    try {
      let mirrorURL = "";
      if (mirrorMode.value === "builtin") {
        mirrorURL = builtinMirrorURL;
      } else if (mirrorMode.value === "custom") {
        mirrorURL = githubMirrorURLInput.value.trim();
        if (!mirrorURL) {
          showMessage("请填写自定义加速 URL", "error");
          return;
        }
      }

      await setGitHubMirrorURL(mirrorURL);
      showMessage("下载源设置已保存", "success");
      await loadStatus();
    } catch (error) {
      showMessage(error instanceof Error ? error.message : "保存加速 URL 失败", "error");
    }
  });
};

const handleClearMirrorURL = async () => {
  mirrorMode.value = "official";
  githubMirrorURLInput.value = "";
  await handleSaveMirrorURL();
};

onMounted(() => {
  void loadStatus();
});
</script>

<template>
  <div class="settings-page d-flex flex-column ga-4">
    <v-snackbar
      v-model="snackbar"
      :color="snackbarColor"
      location="bottom"
      :timeout="2600"
    >
      {{ snackbarText }}
    </v-snackbar>

    <v-row dense class="settings-layout flex-grow-1">
      <v-col cols="12" md="3">
        <v-card elevation="2" class="h-100">
          <v-list nav density="comfortable">
            <v-list-subheader>设置菜单</v-list-subheader>
            <v-list-item
              prepend-icon="fas fa-cloud-arrow-down"
              title="frpc 管理"
              :active="activePanel === 'frpc'"
              @click="activePanel = 'frpc'"
            />
          </v-list>
        </v-card>
      </v-col>

      <v-col cols="12" md="9">
        <v-card elevation="2" class="h-100">
          <v-card-title class="d-flex align-center justify-space-between">
            <div class="text-h6 font-weight-bold">frpc 管理</div>
            <v-chip size="small" color="primary" variant="tonal">
              {{ status?.goos }}/{{ status?.goarch }}
            </v-chip>
          </v-card-title>
          <v-divider />
          <v-card-text class="d-flex flex-column ga-4">
            <div class="d-flex flex-wrap ga-3">
              <v-btn
                color="primary"
                :loading="installing"
                :disabled="installing"
                @click="handleInstallOrUpdate"
              >
                <v-icon start>fas fa-download</v-icon>
                {{ actionText }}
              </v-btn>
              <v-btn
                color="warning"
                variant="tonal"
                :loading="canceling"
                :disabled="!installing || canceling"
                @click="handleCancelInstall"
              >
                <v-icon start>fas fa-stop</v-icon>
                终止下载
              </v-btn>
              <v-btn
                color="info"
                variant="tonal"
                :disabled="installing"
                @click="loadStatus"
              >
                <v-icon start>fas fa-rotate</v-icon>
                检查更新
              </v-btn>
              <v-btn
                color="error"
                variant="tonal"
                :disabled="installing"
                @click="handleRemove"
              >
                <v-icon start>fas fa-trash</v-icon>
                删除本地 frpc
              </v-btn>
            </div>

            <v-alert
              v-if="status?.latest_error"
              type="warning"
              variant="tonal"
              density="compact"
              class="message-alert"
            >
              获取最新版本失败：{{ status.latest_error }}
            </v-alert>

            <v-row dense>
              <v-col cols="12" md="6">
                <v-sheet border rounded="lg" class="pa-3">
                  <div class="text-subtitle-2 mb-2">安装状态</div>
                  <div class="text-body-2">当前版本：{{ installedVersion }}</div>
                  <div class="text-body-2">
                    二进制：{{ status?.installed?.binary_exists ? "已安装" : "未安装" }}
                  </div>
                  <div class="text-body-2">
                    安装时间：{{ formatTime(status?.installed?.installed_at) }}
                  </div>
                </v-sheet>
              </v-col>
              <v-col cols="12" md="6">
                <v-sheet border rounded="lg" class="pa-3">
                  <div class="text-subtitle-2 mb-2">最新版本</div>
                  <div class="text-body-2">标签：{{ latestVersion }}</div>
                  <div class="text-body-2">
                    发布：{{ formatTime(status?.latest?.published_at) }}
                  </div>
                  <div class="text-body-2">
                    可更新：{{ status?.update_available ? "是" : "否" }}
                  </div>
                </v-sheet>
              </v-col>
            </v-row>

            <v-divider />

            <div class="text-subtitle-1 font-weight-bold">本地路径</div>
            <v-sheet border rounded="lg" class="pa-3">
              <div class="text-body-2 text-wrap">
                userdata: {{ status?.paths.userdata_dir }}
              </div>
              <div class="text-body-2 text-wrap">frpc: {{ status?.paths.frpc_dir }}</div>
              <div class="text-body-2 text-wrap">bin: {{ status?.paths.bin_dir }}</div>
              <div class="text-body-2 text-wrap">
                binary: {{ status?.paths.binary_path }}
              </div>
              <div class="text-body-2 text-wrap">
                downloads: {{ status?.paths.download_dir }}
              </div>
              <div class="text-body-2 text-wrap">state: {{ status?.paths.state_path }}</div>
              <div class="text-body-2 text-wrap">
                settings: {{ status?.paths.settings_path }}
              </div>
            </v-sheet>

            <v-divider />

            <v-sheet border rounded="lg" class="pa-3 d-flex flex-column ga-3">
              <div class="text-subtitle-2">GitHub 下载源</div>
              <v-select
                v-model="mirrorMode"
                :items="mirrorModeItems"
                item-title="title"
                item-value="value"
                hide-details="auto"
                :disabled="installing"
              />
              <v-text-field
                v-if="mirrorMode === 'custom'"
                v-model="githubMirrorURLInput"
                hide-details="auto"
                placeholder="https://example.com/github.com"
                :disabled="installing"
              />
              <div class="d-flex flex-wrap ga-2">
                <v-btn
                  color="primary"
                  variant="tonal"
                  :disabled="installing"
                  @click="handleSaveMirrorURL"
                >
                  保存设置
                </v-btn>
                <v-btn
                  color="secondary"
                  variant="text"
                  :disabled="installing"
                  @click="handleClearMirrorURL"
                >
                  使用 github.com
                </v-btn>
              </div>
            </v-sheet>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<style scoped>
.settings-page {
  min-height: calc(100vh - 64px - 32px);
}

.settings-layout :deep(.v-col) {
  display: flex;
}

.settings-layout :deep(.v-card) {
  flex: 1;
}

</style>
