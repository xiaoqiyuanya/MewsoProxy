<template>
  <div class="page-container">
    <div class="card">
      <div class="table-head">
        <h3 style="margin: 0">套餐管理</h3>
        <t-space size="small">
          <t-button theme="primary" @click="openCreate">新增套餐</t-button>
          <t-button theme="default" variant="outline" @click="load">刷新</t-button>
        </t-space>
      </div>
      <t-table :data="list" :columns="columns" row-key="id" :loading="loading">
        <template #op="{ row }">
          <t-space size="small">
            <t-button theme="primary" variant="text" @click="openEdit(row)">编辑</t-button>
            <t-button theme="danger" variant="text" @click="onDrop(row)">删除</t-button>
          </t-space>
        </template>
      </t-table>
    </div>

    <t-dialog v-model:visible="dialogVisible" :header="form.id ? '编辑套餐' : '新增套餐'" :confirm-btn="{ onClick: submit }" :on-close="resetForm">
      <t-form :data="form" label-width="96px">
        <t-form-item label="名称" required>
          <t-input v-model="form.name" placeholder="套餐名称" />
        </t-form-item>
        <t-form-item label="分组ID" required>
          <t-input-number v-model="form.group_id" :min="1" :step="1" />
        </t-form-item>
        <t-form-item label="流量(GB)">
          <t-input-number v-model="form.transfer_gb" :min="0" :step="1" />
        </t-form-item>
        <t-form-item label="限速(Mbps)">
          <t-input-number v-model="form.speed_limit" :min="0" :step="1" />
        </t-form-item>
        <t-form-item label="月付(元)">
          <t-input-number v-model="form.month_price" :min="0" :step="0.5" />
        </t-form-item>
        <t-form-item label="季付(元)">
          <t-input-number v-model="form.quarter_price" :min="0" :step="0.5" />
        </t-form-item>
        <t-form-item label="半年付(元)">
          <t-input-number v-model="form.half_year_price" :min="0" :step="0.5" />
        </t-form-item>
        <t-form-item label="年付(元)">
          <t-input-number v-model="form.year_price" :min="0" :step="0.5" />
        </t-form-item>
        <t-form-item label="两年付(元)">
          <t-input-number v-model="form.two_year_price" :min="0" :step="0.5" />
        </t-form-item>
        <t-form-item label="三年付(元)">
          <t-input-number v-model="form.three_year_price" :min="0" :step="0.5" />
        </t-form-item>
        <t-form-item label="一次性(元)">
          <t-input-number v-model="form.onetime_price" :min="0" :step="0.5" />
        </t-form-item>
        <t-form-item label="排序">
          <t-input-number v-model="form.sort" :min="0" :step="1" />
        </t-form-item>
        <t-form-item label="展示">
          <t-switch v-model="form.show" />
        </t-form-item>
        <t-form-item label="可续费">
          <t-switch v-model="form.renew" />
        </t-form-item>
        <t-form-item label="说明">
          <t-textarea v-model="form.content" placeholder="套餐说明" />
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { MessagePlugin } from 'tdesign-vue-next';
import { listPlansAdmin, savePlan, dropPlan, AdminPlanSaveReq } from '@/api/admin';
import type { PlanDTO } from '@/api/plan';

const list = ref<PlanDTO[]>([]);
const loading = ref(false);
const dialogVisible = ref(false);
const form = reactive<AdminPlanSaveReq & { transfer_gb?: number }>({
  id: undefined,
  name: '',
  group_id: 1,
  transfer_enable: 0,
  transfer_gb: 0,
  speed_limit: undefined,
  show: true,
  renew: true,
  sort: 0,
  content: '',
  month_price: undefined,
  quarter_price: undefined,
  half_year_price: undefined,
  year_price: undefined,
  two_year_price: undefined,
  three_year_price: undefined,
  onetime_price: undefined,
});

const columns = [
  { colKey: 'id', title: 'ID', width: 70 },
  { colKey: 'name', title: '名称' },
  { colKey: 'group_id', title: '分组ID', width: 90 },
  { colKey: 'transfer_enable', title: '流量', cell: (_: unknown, row: PlanDTO) => formatGB(row.transfer_enable) },
  { colKey: 'month_price', title: '月付', cell: (_: unknown, row: PlanDTO) => formatYuan(row.month_price) },
  { colKey: 'year_price', title: '年付', cell: (_: unknown, row: PlanDTO) => formatYuan(row.year_price) },
  { colKey: 'show', title: '展示', cell: (_: unknown, row: PlanDTO) => (row.show ? '是' : '否') },
  { colKey: 'renew', title: '续费', cell: (_: unknown, row: PlanDTO) => (row.renew ? '是' : '否') },
  { colKey: 'sort', title: '排序', cell: (_: unknown, row: PlanDTO) => row.sort ?? '-' },
  { colKey: 'op', title: '操作', width: 140 },
];

async function load() {
  loading.value = true;
  try {
    list.value = await listPlansAdmin();
  } finally {
    loading.value = false;
  }
}

function resetForm() {
  Object.assign(form, {
    id: undefined,
    name: '',
    group_id: 1,
    transfer_gb: 0,
    speed_limit: undefined,
    show: true,
    renew: true,
    sort: 0,
    content: '',
    month_price: undefined,
    quarter_price: undefined,
    half_year_price: undefined,
    year_price: undefined,
    two_year_price: undefined,
    three_year_price: undefined,
    onetime_price: undefined,
  });
}

function openCreate() {
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: PlanDTO) {
  resetForm();
  form.id = row.id;
  form.name = row.name;
  form.group_id = row.group_id;
  form.transfer_gb = Number((row.transfer_enable / 1024 / 1024 / 1024).toFixed(2));
  form.speed_limit = row.speed_limit ?? undefined;
  form.show = row.show;
  form.renew = row.renew;
  form.sort = row.sort ?? 0;
  form.content = row.content ?? '';
  form.month_price = toYuan(row.month_price);
  form.quarter_price = toYuan(row.quarter_price);
  form.half_year_price = toYuan(row.half_year_price);
  form.year_price = toYuan(row.year_price);
  form.two_year_price = toYuan(row.two_year_price);
  form.three_year_price = toYuan(row.three_year_price);
  form.onetime_price = toYuan(row.onetime_price);
  dialogVisible.value = true;
}

async function submit() {
  if (!form.name) {
    MessagePlugin.error('请输入套餐名称');
    return;
  }
  const data: AdminPlanSaveReq = {
    id: form.id,
    name: form.name,
    group_id: form.group_id ?? 1,
    transfer_enable: Math.round((form.transfer_gb ?? 0) * 1024 * 1024 * 1024),
    speed_limit: form.speed_limit,
    show: form.show,
    renew: form.renew,
    sort: form.sort,
    content: form.content || undefined,
    month_price: toFen(form.month_price),
    quarter_price: toFen(form.quarter_price),
    half_year_price: toFen(form.half_year_price),
    year_price: toFen(form.year_price),
    two_year_price: toFen(form.two_year_price),
    three_year_price: toFen(form.three_year_price),
    onetime_price: toFen(form.onetime_price),
  };
  await savePlan(data);
  MessagePlugin.success('已保存');
  dialogVisible.value = false;
  load();
}

async function onDrop(row: PlanDTO) {
  await dropPlan(row.id);
  MessagePlugin.success('已删除');
  load();
}

function formatGB(bytes: number): string {
  return `${(bytes / 1024 / 1024 / 1024).toFixed(2)} GB`;
}

function formatYuan(v?: number): string {
  if (v === undefined) return '-';
  return `${Number(v) / 100}`;
}

function toYuan(v?: number): number | undefined {
  if (v === undefined) return undefined;
  return Number(v) / 100;
}

function toFen(v?: number): number | undefined {
  if (v === undefined) return undefined;
  return Math.round(v * 100);
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
