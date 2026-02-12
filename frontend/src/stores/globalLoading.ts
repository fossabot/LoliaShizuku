import { defineStore } from "pinia";

export const useGlobalLoadingStore = defineStore("globalLoading", {
  state: () => ({
    pendingCount: 0,
  }),
  getters: {
    isLoading: (state) => state.pendingCount > 0,
  },
  actions: {
    startLoading() {
      this.pendingCount += 1;
    },
    stopLoading() {
      this.pendingCount = Math.max(0, this.pendingCount - 1);
    },
    async withGlobalLoading<T>(task: () => Promise<T>): Promise<T> {
      this.startLoading();
      try {
        return await task();
      } finally {
        this.stopLoading();
      }
    },
  },
});
