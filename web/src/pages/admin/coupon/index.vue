<template>
  <div class="page-container">
    <div class="card">
      <div class="table-head">
        <h3 style="margin: 0">优惠券管理</h3>
        <t-space size="small">
          <t-button theme="primary" @click="openCreate">新增优惠券</t-button>
          <t-button theme="default" variant="outline" @click="load">刷新</t-button>
        </t-space>
      </div>
      <t-table :data="list" :columns="columns" row-key="id" :loading="loading">
        <template #op="{ row }">
          <t-space size="small">
            <t-button theme="primary" variant="text" @click="openEdit(row)">编辑</t-button>
            <t-button theme="default" variant="text" @click="toggleShow(row)">{{ row.show ? '隐藏' : '展示' }}</t-button>
            <t-button theme="danger" variant="text" @click="onDrop(row)">删除</t-button>
          </t-space>
        </template>
      </t-table>
    </div>

    <t-dialog v-model:visible="dialogVisible" :header="form.id ? '编辑优惠券' : '新增优惠券'" width="520px" :confirm-btn="{ onClick: submit }" :on-close="resetForm">
      <t-form :data="form" label-width="110px">
        <t-form-item label="券码" required>
          <t-input v-model="form.code" placeholder="优惠券码" />
        </t-form-item>
        <t-form-item label="名称">
          <t-input v-model="form.name" placeholder="优惠券名称" />
        </t-form-item>
        <t-form-item label="类型">
          <t-select v-model="form.type" :options="[{ label: '折扣券', value: 0 }, { label: '金额券', value: 1 }]" />
        </t-form-item>
        <t-form-item :label="form.type === 1 ? '金额(元)' : '折扣(%)'">
          <t-input-number v-model="form.value" :min="0" :step="0.5" />
        </t-form-item>
        <t-form-item label="限用次数">
          <t-input-number v-model="form.limit_use" :min="0" :step="1" />
        </t-form-item>
        <t-form-item label="每人限用">
          <t-input-number v-model="form.limit_use_with_user" :min="0" :step="1" />
        </t-form-item>
        <t-form-item label="生效时间">
          <t-date-picker v-model="form.started_at_ms" value-type="time-stamp" mode="date" clearable style="width: 100%" />
        </t-form-item>
        <t-form-item label="结束时间">
          <t-date-picker v-model="form.ended_at_ms" value-type="time-stamp" mode="date" clearable style="width: 100%" />
        </t-form-item>
        <t-form-item label="展示">
          <t-switch v-model="form.show" />
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { MessagePlugin } from 'tdesign-vue-next';
import { listCoupons, saveCoupon, dropCoupon, toggleCouponShow, Coupon, AdminCouponSaveReq } from '@/api/admin';

const list = ref<Coupon[]>([]);
const loading = ref(false);
const dialogVisible = ref(false);
const form = reactive<AdminCouponSaveReq & { started_at_ms?: number; ended_at_ms?: number }>({
  id: undefined,
  code: '',
  name: '',
  type: 0,
  value: 0,
  show: true,
  limit_use: undefined,
  limit_use_with_user: undefined,
  started_at: 0,
  ended_at: 0,
  started_at_ms: undefined,
  ended_at_ms: undefined,
});

const columns = [
  { colKey: 'id', title: 'ID', width: 70 },
  { colKey: 'code', title: '券码' },
  { colKey: 'name', title: '名称' },
  { colKey: 'type', title: '类型', cell: (_: unknown, row: Coupon) => (row.type === 1 ? '金额券' : '折扣券') },
  { colKey: 'value', title: '面值', cell: (_: unknown, row: Coupon) => (row.type === 1 ? `${Number(row.value) / 100} 元` : `${row.value}%`) },
  { colKey: 'limit_use', title: '限用', cell: (_: unknown, row: Coupon) => row.limit_use ?? '-' },
  { colKey: 'show', title: '展示', cell: (_: unknown, row: Coupon) => (row.show ? '是' : '否') },
  { colKey: 'started_at', title: '生效', cell: (_: unknown, row: Coupon) => formatTime(row.started_at) },
  { colKey: 'ended_at', title: '结束', cell: (_: unknown, row: Coupon) => formatTime(row.ended_at) },
  { colKey: 'op', title: '操作', width: 200 },
];

async function load() {
  loading.value = true;
  try {
    list.value = await listCoupons();
  } finally {
    loading.value = false;
  }
}

function resetForm() {
  Object.assign(form, {
    id: undefined,
    code: '',
    name: '',
    type: 0,
    value: 0,
    show: true,
    limit_use: undefined,
    limit_use_with_user: undefined,
    started_at_ms: undefined,
    ended_at_ms: undefined,
  });
}

function openCreate() {
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: Coupon) {
  resetForm();
  form.id = row.id;
  form.code = row.code;
  form.name = row.name;
  form.type = row.type;
  form.value = row.type === 1 ? Number(row.value) / 100 : row.value;
  form.show = row.show;
  form.limit_use = row.limit_use ?? undefined;
  form.limit_use_with_user = row.limit_use_with_user ?? undefined;
  form.started_at_ms = row.started_at ? row.started_at * 1000 : undefined;
  form.ended_at_ms = row.ended_at ? row.ended_at * 1000 : undefined;
  dialogVisible.value = true;
}

async function submit() {
  if (!form.code) {
    MessagePlugin.error('请输入券码');
    return;
  }
  const data: AdminCouponSaveReq = {
    id: form.id,
    code: form.code,
    name: form.name,
    type: form.type ?? 0,
    value: form.type === 1 ? Math.round((form.value ?? 0) * 100) : Math.round(form.value ?? 0),
    show: form.show,
    limit_use: form.limit_use,
    limit_use_with_user: form.limit_use_with_user,
    started_at: form.started_at_ms ? Math.floor(form.started_at_ms / 1000) : 0,
    ended_at: form.ended_at_ms ? Math.floor(form.ended_at_ms / 1000) : 0,
  };
  await saveCoupon(data);
  MessagePlugin.success('已保存');
  dialogVisible.value = false;
  load();
}

async function toggleShow(row: Coupon) {
  await toggleCouponShow(row.id, !row.show);
  MessagePlugin.success('已更新');
  load();
}

async function onDrop(row: Coupon) {
  await dropCoupon(row.id);
  MessagePlugin.success('已删除');
  load();
}

function formatTime(v?: number): string {
  if (!v) return '-';
  return new Date(v * 1000).toLocaleDateString();
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
