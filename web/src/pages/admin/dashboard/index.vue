<template>
  <div class="page-container">
    <t-row :gutter="[16, 16]">
      <t-col :span="4">
        <t-card :bordered="false">
          <t-statistic title="用户总数" :value="status.user_count" suffix="人" />
        </t-card>
      </t-col>
      <t-col :span="4">
        <t-card :bordered="false">
          <t-statistic title="订单总数" :value="status.order_count" suffix="单" />
        </t-card>
      </t-col>
      <t-col :span="4">
        <t-card :bordered="false">
          <t-statistic title="今日实付" :value="(status.today_paid_total / 100).toFixed(2)" prefix="¥" />
        </t-card>
      </t-col>
      <t-col :span="4">
        <t-card :bordered="false">
          <t-statistic title="在线用户" :value="status.online_user_count" suffix="人" />
        </t-card>
      </t-col>
      <t-col :span="4">
        <t-card :bordered="false">
          <t-statistic title="在线节点" :value="status.online_node_count" suffix="个" />
        </t-card>
      </t-col>
      <t-col :span="4">
        <t-card :bordered="false">
          <t-statistic title="今日流量" :value="formatBytes(status.today_traffic)" />
        </t-card>
      </t-col>
      <t-col :span="4">
        <t-card :bordered="false">
          <t-statistic title="活跃用户" :value="status.active_user_count" suffix="人" />
        </t-card>
      </t-col>
      <t-col :span="4">
        <t-card :bordered="false">
          <t-statistic title="数据库" :value="status.db_status === 'ok' ? '正常' : '异常'" />
        </t-card>
      </t-col>
      <t-col :span="4">
        <t-card :bordered="false">
          <t-statistic title="Redis" :value="status.redis_status === 'ok' ? '正常' : '异常'" />
        </t-card>
      </t-col>
    </t-row>

    <t-card :bordered="false" style="margin-top: 16px">
      <t-descriptions :column="2">
        <t-descriptions-item label="服务器时间">{{ formatTime(status.server_time) }}</t-descriptions-item>
        <t-descriptions-item label="状态">运行中</t-descriptions-item>
      </t-descriptions>
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { fetchSystemStatus, AdminSystemStatus } from '@/api/admin';

const status = reactive<AdminSystemStatus>({
  server_time: 0,
  db_status: '-',
  redis_status: '-',
  online_user_count: 0,
  user_count: 0,
  order_count: 0,
  today_paid_total: 0,
  online_node_count: 0,
  active_user_count: 0,
  today_traffic: 0,
});

onMounted(async () => {
  const res = await fetchSystemStatus();
  Object.assign(status, res);
});

function formatTime(v: number): string {
  if (!v) return '-';
  return new Date(v * 1000).toLocaleString();
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
