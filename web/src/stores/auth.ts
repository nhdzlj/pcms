import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { login as apiLogin, register as apiRegister, getMe } from "@/api/auth";
import type { UserInfo } from "@/api/auth";
import router from "@/router";

export const useAuthStore = defineStore("auth", () => {
  const token = ref(localStorage.getItem("token") || "");
  const user = ref<UserInfo | null>(
    JSON.parse(localStorage.getItem("user") || "null")
  );

  const isLoggedIn = computed(() => !!token.value);
  const username = computed(() => user.value?.username || "");

  async function login(username: string, password: string) {
    const result = await apiLogin({ username, password });
    token.value = result.token;
    user.value = result.user;
    localStorage.setItem("token", result.token);
    localStorage.setItem("user", JSON.stringify(result.user));
  }

  async function register(username: string, password: string, email?: string) {
    const result = await apiRegister({ username, password, email });
    token.value = result.token;
    user.value = result.user;
    localStorage.setItem("token", result.token);
    localStorage.setItem("user", JSON.stringify(result.user));
  }

  async function fetchUser() {
    if (!token.value) return;
    try {
      const u = await getMe();
      user.value = u;
      localStorage.setItem("user", JSON.stringify(u));
    } catch {
      logout();
    }
  }

  function logout() {
    token.value = "";
    user.value = null;
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    router.push("/login");
  }

  return {
    token,
    user,
    isLoggedIn,
    username,
    login,
    register,
    fetchUser,
    logout,
  };
});
