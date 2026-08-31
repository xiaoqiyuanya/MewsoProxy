<template>
  <div class="page-container">
    <t-row :gutter="[16, 16]">
      <t-col :span="8">
        <t-card title="账户概览" :bordered="false">
          <t-descriptions :column="1">
            <t-descriptions-item label="邮箱">{{ user?.email }}</t-descriptions-item>
            <t-descriptions-item label="余额">
              {{
                Number(user?.balance || 0) / 100
              }} 元
            </t-descriptions-item>
            <t-descriptions-item label="到期时间">
              {{ formatTime(user?.expired_at) }}
            </t-descriptions-item>
            <t-descriptions-item label="订阅 Token">
              <code>{{ user?.token }}</code>
            </t-descriptions-item>
            <t-descriptions-item label="订阅地址">
              <code style="word-break: break-all">{{ subUrl }}</code>
              <t-button theme="primary" variant="outline" size="small" style="margin-left: 8px" @click="copyUrl">复制</t-button>
            </t-descriptions-item>
          </t-descriptions>
        </t-card>
      </t-col>
      <t-col :span="8">
        <t-card title="流量使用" :bordered="false">
          <t-descriptions :column="1">
            <t-descriptions-item label="总流量">
              {{ formatBytes(user?.transfer_enable || 0) }}
            </t-descriptions-item>
            <t-descriptions-item label="已用流量">
              {{ formatBytes(user?.used_traffic || 0) }}
            </t-descriptions-item>
          </t-descriptions>
          <t-progress :percentage="trafficPercent" :label="true" />
        </t-card>
      </t-col>
    </t-row>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { MessagePlugin } from 'tdesign-vue-next';
import { getSubscribe } from '@/api/auth';
import { useUserStore } from '@/store/user';

const store = useUserStore();
const user = computed(() => store.user);
const subUrl = ref('');

const trafficPercent = computed(() => {
  const total = user.value?.transfer_enable || 0;
  if (total <= 0) return 0;
  return Math.min(100, Math.round(((user.value?.used_traffic || 0) / total) * 100));
});

async function loadSubscribe() {
  try {
    const r = await getSubscribe();
    subUrl.value = r.url;
  } catch {
    // 忽略，未登录时无订阅
  }
}

async function copyUrl() {
  if (!subUrl.value) return;
  await navigator.clipboard.writeText(subUrl.value);
  MessagePlugin.success('已复制订阅地址');
}

onMounted(loadSubscribe);

function formatTime(v?: string): string {
  if (!v) return '-';
  return new Date(v).toLocaleString();
}

function formatBytes(bytes: number): string {
  if (bytes <= 0) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  let i = 0;
  let n = bytes;
  while (n >= 1024 && i < units.length - 1) {
    n /= 1024;
    i++;
  }
  return `${n.toFixed(2)} ${units[i]}`;
}
</script>
