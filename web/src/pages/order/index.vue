<template>
  <div class="page-container">
    <div class="card">
      <div class="table-head">
        <h3 style="margin: 0">我的订单</h3>
        <t-button theme="default" variant="outline" @click="load">刷新</t-button>
      </div>
      <t-table
        :data="list"
        :columns="columns"
        row-key="id"
        :loading="loading"
        :pagination="pagination"
        @page-change="onPageChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { MessagePlugin } from 'tdesign-vue-next';
import { listOrders, OrderDTO } from '@/api/order';

const list = ref<OrderDTO[]>([]);
const loading = ref(false);
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showJumper: true,
});

const columns = [
  { colKey: 'trade_no', title: '订单号', width: 220 },
  { colKey: 'period', title: '周期' },
  { colKey: 'total_amount', title: '金额', cell: (_: unknown, row: OrderDTO) => `${Number(row.total_amount) / 100} 元` },
  {
    colKey: 'type',
    title: '类型',
    cell: (_: unknown, row: OrderDTO) =>
      ({ 1: '新购', 2: '续费', 3: '升级' }[row.type] || '-'),
  },
  {
    colKey: 'status',
    title: '状态',
    cell: (_: unknown, row: OrderDTO) =>
      ({ 0: '待支付', 3: '已完成', 2: '已取消' }[row.status] || '处理中'),
  },
  { colKey: 'created_at', title: '创建时间', cell: (_: unknown, row: OrderDTO) => formatTime(row.created_at) },
];

async function load() {
  loading.value = true;
  try {
    const res = await listOrders({ page: pagination.current, page_size: pagination.pageSize });
    list.value = res.list;
    pagination.total = res.total;
  } finally {
    loading.value = false;
  }
}

function onPageChange({ current, pageSize }: { current: number; pageSize: number }) {
  pagination.current = current;
  pagination.pageSize = pageSize;
  load();
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
