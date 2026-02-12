import { defineStore } from "pinia";
import {
  cancelInstallOrUpdateFrpc,
  installOrUpdateFrpc,
  type FrpcInstallResult,
} from "@/services/frpc";

export const useFrpcInstallStore = defineStore("frpcInstall", {
  state: () => ({
    installing: false,
    canceling: false,
    runningPromise: null as Promise<FrpcInstallResult> | null,
  }),
  actions: {
    async startInstall(): Promise<FrpcInstallResult> {
      if (this.runningPromise) {
        return this.runningPromise;
      }

      this.installing = true;
      this.canceling = false;

      this.runningPromise = (async () => {
        try {
          return await installOrUpdateFrpc();
        } finally {
          this.installing = false;
          this.canceling = false;
          this.runningPromise = null;
        }
      })();

      return this.runningPromise;
    },

    async cancelInstall(): Promise<void> {
      if (!this.runningPromise || this.canceling) {
        return;
      }

      this.canceling = true;
      try {
        await cancelInstallOrUpdateFrpc();
      } catch (error) {
        this.canceling = false;
        throw error;
      }
    },
  },
});
