import { defineStore } from 'pinia';
import { getMe, UserDTO } from '@/api/auth';
import { getToken, setToken, clearToken } from '@/utils/token';

export const useUserStore = defineStore('user', {
  state: () => ({
    token: getToken(),
    user: null as UserDTO | null,
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
    isAdmin: (state) => !!state.user?.is_admin,
  },
  actions: {
    async fetchMe() {
      this.user = await getMe();
      return this.user;
    },
    saveToken(token: string) {
      this.token = token;
      setToken(token);
    },
    reset() {
      this.token = '';
      this.user = null;
      clearToken();
    },
  },
});
