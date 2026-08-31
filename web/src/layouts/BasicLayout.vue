<template>
  <t-layout style="height: 100%">
    <t-aside width="232px" style="border-right: 1px solid var(--td-border-level-2-color)">
      <div class="logo">
        <span class="logo-mark">M</span>
        <span>MewsoProxy</span>
      </div>
      <t-menu v-model="active" theme="light" style="border-right: none" @change="onMenuChange">
        <t-menu-item value="/dashboard">
          <template #icon><dashboard-icon /></template>
          概览
        </t-menu-item>
        <t-menu-item value="/plan">
          <template #icon><chart-icon /></template>
          套餐
        </t-menu-item>
        <t-menu-item value="/order">
          <template #icon><cart-icon /></template>
          订单
        </t-menu-item>
      </t-menu>
    </t-aside>

    <t-layout>
      <t-header class="header">
        <t-space align="center">
          <h3 style="margin: 0">{{ title }}</h3>
        </t-space>
        <t-space align="center" size="small">
          <t-avatar size="small">{{ avatarText }}</t-avatar>
          <span>{{ user?.email || '未登录' }}</span>
          <t-button theme="default" variant="text" @click="handleLogout">退出</t-button>
        </t-space>
      </t-header>
      <t-content style="overflow: auto">
        <router-view />
      </t-content>
    </t-layout>
  </t-layout>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import {
  DashboardIcon,
  ChartIcon,
  CartIcon,
} from 'tdesign-icons-vue-next';
import { useUserStore } from '@/store/user';
import { logout } from '@/api/auth';

const route = useRoute();
const router = useRouter();
const store = useUserStore();

const active = ref<string>('/dashboard');
const user = computed(() => store.user);
const title = computed(() => (route.meta.title as string) || 'MewsoProxy');
const avatarText = computed(() => (user.value?.email || 'M').slice(0, 1).toUpperCase());

watch(
  () => route.path,
  (p) => {
    active.value = p;
  },
  { immediate: true },
);

function onMenuChange(value: string | number) {
  router.push(String(value));
}

async function handleLogout() {
  try {
    await logout();
  } finally {
    store.reset();
    router.replace('/login');
  }
}
</script>

<style scoped>
.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 56px;
  padding: 0 20px;
  font-size: 17px;
  font-weight: 600;
  color: #1f2329;
}
.logo-mark {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 8px;
  background: var(--td-brand-color);
  color: #fff;
  font-size: 15px;
  font-weight: 700;
}
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background: var(--td-bg-color-container);
  border-bottom: 1px solid var(--td-border-level-2-color);
}
</style>
