<template>
  <t-layout style="height: 100%">
    <t-aside width="232px" style="border-right: 1px solid var(--td-border-level-2-color)">
      <div class="logo">
        <span>MewsoProxy</span>
        <span class="badge">Admin</span>
      </div>
      <t-menu v-model="active" theme="light" style="border-right: none">
        <t-menu-item value="/admin/dashboard">
          <template #icon><dashboard-icon /></template>
          概览
        </t-menu-item>
        <t-menu-item value="/admin/user">
          <template #icon><user-icon /></template>
          用户
        </t-menu-item>
        <t-menu-item value="/admin/plan">
          <template #icon><shop-icon /></template>
          套餐
        </t-menu-item>
        <t-menu-item value="/admin/order">
          <template #icon><cart-icon /></template>
          订单
        </t-menu-item>
        <t-menu-item value="/admin/server">
          <template #icon><server-icon /></template>
          节点
        </t-menu-item>
        <t-menu-item value="/admin/coupon">
          <template #icon><coupon-icon /></template>
          优惠券
        </t-menu-item>
        <t-menu-item value="/admin/notice">
          <template #icon><notification-icon /></template>
          公告
        </t-menu-item>
        <t-menu-item value="/admin/payment">
          <template #icon><wallet-icon /></template>
          支付
        </t-menu-item>
        <t-menu-item value="/admin/config">
          <template #icon><setting-icon /></template>
          配置
        </t-menu-item>
      </t-menu>
    </t-aside>

    <t-layout>
      <t-header class="header">
        <t-space align="center">
          <h3 style="margin: 0">{{ title }}</h3>
        </t-space>
        <t-space align="center" size="small">
          <t-dropdown :options="dropdownOptions" @click="onDropdownClick">
            <t-avatar size="small" style="cursor: pointer">{{ avatarText }}</t-avatar>
          </t-dropdown>
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
import { MessagePlugin } from 'tdesign-vue-next';
import {
  DashboardIcon,
  UserIcon,
  ShopIcon,
  CartIcon,
  ServerIcon,
  CouponIcon,
  NotificationIcon,
  WalletIcon,
  SettingIcon,
} from 'tdesign-icons-vue-next';
import { useUserStore } from '@/store/user';
import { logout } from '@/api/auth';

const route = useRoute();
const router = useRouter();
const store = useUserStore();

const active = ref<string>('/admin/dashboard');
const user = computed(() => store.user);
const title = computed(() => (route.meta.title as string) || '管理后台');
const avatarText = computed(() => (user.value?.email || 'M').slice(0, 1).toUpperCase());

const dropdownOptions = [
  { content: '返回前台', value: 'frontend' },
];

watch(
  () => route.path,
  (p) => {
    active.value = p;
  },
  { immediate: true },
);

function onDropdownClick(data: { value: string | number }) {
  if (data.value === 'frontend') {
    router.push('/dashboard');
  }
}

async function handleLogout() {
  try {
    await logout();
  } finally {
    store.reset();
    MessagePlugin.success('已退出登录');
    router.replace('/login');
  }
}
</script>

<style scoped>
.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 56px;
  padding: 0 20px;
  font-size: 18px;
  font-weight: 600;
  color: var(--td-brand-color);
}
.badge {
  font-size: 12px;
  font-weight: 500;
  color: #fff;
  background: var(--td-brand-color);
  border-radius: 10px;
  padding: 1px 8px;
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
