import { computed, ref } from "vue";

const pendingCount = ref(0);

const isLoading = computed(() => pendingCount.value > 0);

const startLoading = () => {
  pendingCount.value += 1;
};

const stopLoading = () => {
  pendingCount.value = Math.max(0, pendingCount.value - 1);
};

const withGlobalLoading = async <T>(task: () => Promise<T>): Promise<T> => {
  startLoading();
  try {
    return await task();
  } finally {
    stopLoading();
  }
};

export const useGlobalLoading = () => ({
  isLoading,
  startLoading,
  stopLoading,
  withGlobalLoading,
});
