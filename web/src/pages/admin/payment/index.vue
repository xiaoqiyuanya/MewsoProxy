<template>
  <div class="page-container">
    <div class="card">
      <div class="table-head">
        <h3 style="margin: 0">支付管理</h3>
        <t-space size="small">
          <t-button theme="primary" @click="openCreate">新增支付</t-button>
          <t-button theme="default" variant="outline" @click="load">刷新</t-button>
        </t-space>
      </div>
      <t-table :data="list" :columns="columns" row-key="id" :loading="loading">
        <template #op="{ row }">
          <t-space size="small">
            <t-button theme="primary" variant="text" @click="openEdit(row)">编辑</t-button>
            <t-button theme="default" variant="text" @click="toggleEnable(row)">{{ row.enable ? '停用' : '启用' }}</t-button>
            <t-button theme="danger" variant="text" @click="onDrop(row)">删除</t-button>
          </t-space>
        </template>
      </t-table>
    </div>

    <t-dialog v-model:visible="dialogVisible" :header="form.id ? '编辑支付' : '新增支付'" width="560px" :confirm-btn="{ onClick: submit }" :on-close="resetForm">
      <t-form :data="form" label-width="110px">
        <t-form-item label="UUID" required>
          <t-input v-model="form.uuid" placeholder="支付 UUID" />
        </t-form-item>
        <t-form-item label="渠道" required>
          <t-input v-model="form.payment" placeholder="如 stripe / epay / alipay" />
        </t-form-item>
        <t-form-item label="名称" required>
          <t-input v-model="form.name" placeholder="支付名称" />
        </t-form-item>
        <t-form-item label="图标">
          <t-input v-model="form.icon" placeholder="图标 URL" />
        </t-form-item>
        <t-form-item label="回调域名">
          <t-input v-model="form.notify_domain" placeholder="通知域名" />
        </t-form-item>
        <t-form-item label="固定手续费(元)">
          <t-input-number v-model="form.handling_fee_fixed" :min="0" :step="0.5" />
        </t-form-item>
        <t-form-item label="费率(%)">
          <t-input-number v-model="form.handling_fee_percent" :min="0" :step="0.1" />
        </t-form-item>
        <t-form-item label="配置">
          <t-textarea v-model="form.config" placeholder="渠道配置 JSON" :autosize="{ minRows: 3 }" />
        </t-form-item>
        <t-form-item label="排序">
          <t-input-number v-model="form.sort" :min="0" :step="1" />
        </t-form-item>
        <t-form-item label="启用">
          <t-switch v-model="form.enable" />
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { MessagePlugin } from 'tdesign-vue-next';
import { listPayments, savePayment, dropPayment, togglePaymentShow, Payment, AdminPaymentSaveReq } from '@/api/admin';

const list = ref<Payment[]>([]);
const loading = ref(false);
const dialogVisible = ref(false);
const form = reactive<AdminPaymentSaveReq & { handling_fee_fixed?: number }>({
  id: undefined,
  uuid: '',
  payment: '',
  name: '',
  icon: '',
  notify_domain: '',
  handling_fee_fixed: 0,
  handling_fee_percent: 0,
  config: '',
  sort: 0,
  enable: true,
});

const columns = [
  { colKey: 'id', title: 'ID', width: 70 },
  { colKey: 'name', title: '名称' },
  { colKey: 'payment', title: '渠道' },
  { colKey: 'uuid', title: 'UUID', width: 200 },
  {
    colKey: 'handling_fee_fixed',
    title: '固定手续费',
    cell: (_: unknown, row: Payment) => (row.handling_fee_fixed ? `${Number(row.handling_fee_fixed) / 100} 元` : '-'),
  },
  {
    colKey: 'handling_fee_percent',
    title: '费率',
    cell: (_: unknown, row: Payment) => (row.handling_fee_percent ? `${row.handling_fee_percent}%` : '-'),
  },
  { colKey: 'enable', title: '启用', cell: (_: unknown, row: Payment) => (row.enable ? '是' : '否') },
  { colKey: 'sort', title: '排序', cell: (_: unknown, row: Payment) => row.sort ?? '-' },
  { colKey: 'op', title: '操作', width: 200 },
];

async function load() {
  loading.value = true;
  try {
    list.value = await listPayments();
  } finally {
    loading.value = false;
  }
}

function resetForm() {
  Object.assign(form, {
    id: undefined,
    uuid: '',
    payment: '',
    name: '',
    icon: '',
    notify_domain: '',
    handling_fee_fixed: 0,
    handling_fee_percent: 0,
    config: '',
    sort: 0,
    enable: true,
  });
}

function openCreate() {
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: Payment) {
  resetForm();
  form.id = row.id;
  form.uuid = row.uuid;
  form.payment = row.payment;
  form.name = row.name;
  form.icon = row.icon ?? '';
  form.notify_domain = row.notify_domain ?? '';
  form.handling_fee_fixed = row.handling_fee_fixed ? Number(row.handling_fee_fixed) / 100 : 0;
  form.handling_fee_percent = row.handling_fee_percent ?? 0;
  form.config = row.config ?? '';
  form.sort = row.sort ?? 0;
  form.enable = row.enable;
  dialogVisible.value = true;
}

async function submit() {
  if (!form.uuid || !form.payment || !form.name) {
    MessagePlugin.error('请填写 UUID、渠道和名称');
    return;
  }
  const data: AdminPaymentSaveReq = {
    id: form.id,
    uuid: form.uuid,
    payment: form.payment,
    name: form.name,
    icon: form.icon || undefined,
    notify_domain: form.notify_domain || undefined,
    handling_fee_fixed: Math.round((form.handling_fee_fixed ?? 0) * 100),
    handling_fee_percent: form.handling_fee_percent,
    config: form.config || '',
    sort: form.sort,
    enable: form.enable,
  };
  await savePayment(data);
  MessagePlugin.success('已保存');
  dialogVisible.value = false;
  load();
}

async function toggleEnable(row: Payment) {
  await togglePaymentShow(row.id, !row.enable);
  MessagePlugin.success('已更新');
  load();
}

async function onDrop(row: Payment) {
  await dropPayment(row.id);
  MessagePlugin.success('已删除');
  load();
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
