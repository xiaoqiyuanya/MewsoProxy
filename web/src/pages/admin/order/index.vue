<template>
  <div class="page-container">
    <div class="card">
      <div class="table-head">
        <h3 style="margin: 0">订单管理</h3>
        <t-space size="small">
          <t-input v-model="keyword" placeholder="搜索订单号" clearable style="width: 240px" @enter="load" @clear="load" />
          <t-select v-model="status" :options="statusOptions" style="width: 140px" @change="load" />
          <t-button theme="primary" variant="outline" @click="load">搜索</t-button>
          <t-button theme="default" variant="outline" @click="refresh">刷新</t-button>
        </t-space>
      </div>
      <t-table :data="list" :columns="columns" row-key="id" :loading="loading" :pagination="pagination" @page-change="onPageChange">
        <template #op="{ row }">
          <t-space size="small">
            <t-button theme="default" variant="text" @click="viewDetail(row)">详情</t-button>
            <t-button v-if="row.status === 0" theme="primary" variant="text" @click="onPaid(row)">标记已支付</t-button>
            <t-button v-if="row.status === 0" theme="danger" variant="text" @click="onCancel(row)">取消</t-button>
          </t-space>
        </template>
      </t-table>
    </div>

    <t-dialog v-model:visible="detailVisible" header="订单详情" width="480px">
      <t-descriptions :column="1" bordered v-if="detail">
        <t-descriptions-item label="订单号">{{ detail.trade_no }}</t-descriptions-item>
        <t-descriptions-item label="套餐ID">{{ detail.plan_id }}</t-descriptions-item>
        <t-descriptions-item label="周期">{{ detail.period }}</t-descriptions-item>
        <t-descriptions-item label="类型">{{ typeText(detail.type) }}</t-descriptions-item>
        <t-descriptions-item label="金额">{{ Number(detail.total_amount) / 100 }} 元</t-descriptions-item>
        <t-descriptions-item label="余额抵扣">{{ Number(detail.balance_amount) / 100 }} 元</t-descriptions-item>
        <t-descriptions-item label="优惠">{{ Number(detail.discount_amount) / 100 }} 元</t-descriptions-item>
        <t-descriptions-item label="状态">{{ statusText(detail.status) }}</t-descriptions-item>
        <t-descriptions-item label="创建时间">{{ formatTime(detail.created_at) }}</t-descriptions-item>
        <t-descriptions-item label="支付时间">{{ formatTime(detail.paid_at) }}</t-descriptions-item>
      </t-descriptions>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { MessagePlugin } from 'tdesign-vue-next';
import { listOrdersAdmin, getOrder, cancelOrder, markOrderPaid } from '@/api/admin';
import type { OrderDTO } from '@/api/order';

const list = ref<OrderDTO[]>([]);
const loading = ref(false);
const keyword = ref('');
const status = ref<number | undefined>(undefined);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showJumper: true });
const detailVisible = ref(false);
const detail = ref<OrderDTO | null>(null);

const statusOptions = [
  { label: '全部', value: -1 },
  { label: '待支付', value: 0 },
  { label: '已取消', value: 2 },
  { label: '已完成', value: 3 },
];

const columns = [
  { colKey: 'id', title: 'ID', width: 70 },
  { colKey: 'trade_no', title: '订单号', width: 220 },
  { colKey: 'plan_id', title: '套餐ID', width: 90 },
  { colKey: 'period', title: '周期' },
  { colKey: 'type', title: '类型', cell: (_: unknown, row: OrderDTO) => typeText(row.type) },
  { colKey: 'total_amount', title: '金额', cell: (_: unknown, row: OrderDTO) => `${Number(row.total_amount) / 100} 元` },
  { colKey: 'status', title: '状态', cell: (_: unknown, row: OrderDTO) => statusText(row.status) },
  { colKey: 'created_at', title: '创建时间', cell: (_: unknown, row: OrderDTO) => formatTime(row.created_at) },
  { colKey: 'op', title: '操作', width: 200 },
];

async function load() {
  loading.value = true;
  try {
    const res = await listOrdersAdmin({
      keyword: keyword.value,
      status: status.value,
      page: pagination.current,
      page_size: pagination.pageSize,
    });
    list.value = res.list;
    pagination.total = res.total;
  } finally {
    loading.value = false;
  }
}

function refresh() {
  pagination.current = 1;
  load();
}

function onPageChange({ current, pageSize }: { current: number; pageSize: number }) {
  pagination.current = current;
  pagination.pageSize = pageSize;
  load();
}

async function viewDetail(row: OrderDTO) {
  detail.value = await getOrder(row.id);
  detailVisible.value = true;
}

async function onCancel(row: OrderDTO) {
  await cancelOrder(row.id);
  MessagePlugin.success('已取消');
  load();
}

async function onPaid(row: OrderDTO) {
  await markOrderPaid(row.id);
  MessagePlugin.success('已标记为已支付');
  load();
}

function typeText(t: number): string {
  return ({ 1: '新购', 2: '续费', 3: '升级' } as Record<number, string>)[t] || '-';
}

function statusText(s: number): string {
  return ({ 0: '待支付', 3: '已完成', 2: '已取消' } as Record<number, string>)[s] || '处理中';
}

function formatTime(v?: string): string {
  if (!v) return '-';
  return new Date(v).toLocaleString();
}

onMounted(load);
</script>

<style scoped>
.table-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
</style>
