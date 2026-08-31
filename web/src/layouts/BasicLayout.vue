<template>
  <t-layout style="height: 100%">
    <t-aside width="232px" style="border-right: 1px solid var(--td-border-level-2-color)">
      <div class="logo">MewsoProxy</div>
      <t-menu v-model="active" theme="light" style="border-right: none">
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
  height: 56px;
  padding: 0 20px;
  font-size: 18px;
  font-weight: 600;
  color: var(--td-brand-color);
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
