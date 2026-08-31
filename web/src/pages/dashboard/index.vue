<template>
  <div class="page-container">
    <div class="page-title">
      <div>
        <h2>概览</h2>
        <p class="desc">账号信息与流量使用情况</p>
      </div>
    </div>

    <div class="card account-card">
      <div class="account-main">
        <div class="stat-block">
          <div class="stat-label">账户余额</div>
          <div class="stat-value">
            <span class="money">{{ formatMoney(user?.balance) }}</span> 元
          </div>
        </div>
        <div class="stat-divider"></div>
        <div class="stat-block">
          <div class="stat-label">套餐到期</div>
          <div class="stat-value">{{ formatTime(user?.expired_at) }}</div>
        </div>
      </div>

      <div class="account-fields">
        <div class="account-field">
          <span class="field-label">邮箱</span>
          <span class="field-val">{{ user?.email || '-' }}</span>
        </div>
        <div class="account-field">
          <span class="field-label">订阅 Token</span>
          <span class="field-val">
            <code>{{ user?.token || '-' }}</code>
          </span>
        </div>
        <div class="account-field">
          <span class="field-label">订阅地址</span>
          <div class="field-val sub-field">
            <code class="sub-url">{{ subUrl || '暂无订阅' }}</code>
            <t-button v-if="subUrl" theme="primary" variant="outline" size="small" @click="copyUrl">
              复制
            </t-button>
          </div>
        </div>
      </div>
    </div>

    <div class="card traffic-card">
      <div class="traffic-head">
        <h3>流量使用</h3>
        <span class="traffic-summary">
          已用 {{ formatBytes(user?.used_traffic || 0) }} / 共 {{ formatBytes(user?.transfer_enable || 0) }}
        </span>
      </div>
      <t-progress :percentage="trafficPercent" :label="true" />
      <div class="traffic-stats">
        <div class="stat-block">
          <div class="stat-label">总流量</div>
          <div class="stat-value">{{ formatBytes(user?.transfer_enable || 0) }}</div>
        </div>
        <div class="stat-block">
          <div class="stat-label">已用流量</div>
          <div class="stat-value">{{ formatBytes(user?.used_traffic || 0) }}</div>
        </div>
        <div class="stat-block">
          <div class="stat-label">剩余流量</div>
          <div class="stat-value">{{ formatBytes(remainTraffic) }}</div>
        </div>
      </div>
    </div>
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

const remainTraffic = computed(() => {
  return Math.max(0, (user.value?.transfer_enable || 0) - (user.value?.used_traffic || 0));
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
  return new Date(v).toLocaleDateString();
}

function formatMoney(v?: number): string {
  return (Number(v || 0) / 100).toFixed(2);
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

<style scoped>
.account-card,
.traffic-card {
  margin-bottom: 16px;
}
.account-main {
  display: flex;
  align-items: center;
  gap: 32px;
  padding-bottom: 20px;
}
.stat-block {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.stat-label {
  font-size: 13px;
  color: #8a919f;
}
.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #1f2329;
}
.stat-value .money {
  color: var(--td-brand-color);
}
.stat-divider {
  width: 1px;
  height: 40px;
  background: var(--td-border-level-2-color);
}
.account-fields {
  border-top: 1px solid var(--td-border-level-2-color);
  padding-top: 4px;
}
.account-field {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 10px 0;
}
.field-label {
  flex-shrink: 0;
  width: 88px;
  font-size: 13px;
  color: #8a919f;
}
.field-val {
  flex: 1;
  min-width: 0;
  font-size: 14px;
  color: #1f2329;
  word-break: break-all;
}
.field-val code {
  font-family: var(--td-font-family-mono, Consolas, Menlo, monospace);
  font-size: 13px;
  color: #444;
}
.sub-field {
  display: flex;
  align-items: center;
  gap: 12px;
}
.sub-url {
  flex: 1;
  min-width: 0;
}
.traffic-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-bottom: 20px;
}
.traffic-head h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}
.traffic-summary {
  font-size: 13px;
  color: #8a919f;
}
.traffic-stats {
  display: flex;
  gap: 48px;
  margin-top: 20px;
}
.traffic-stats .stat-value {
  font-size: 18px;
}
</style>
