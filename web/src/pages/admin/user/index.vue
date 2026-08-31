<template>
  <div class="page-container">
    <div class="card">
      <div class="table-head">
        <h3 style="margin: 0">用户管理</h3>
        <t-space size="small">
          <t-input v-model="keyword" placeholder="搜索邮箱" clearable style="width: 240px" @enter="load" @clear="load" />
          <t-button theme="primary" variant="outline" @click="load">搜索</t-button>
          <t-button theme="default" variant="outline" @click="refresh">刷新</t-button>
        </t-space>
      </div>
      <t-table :data="list" :columns="columns" row-key="id" :loading="loading" :pagination="pagination" @page-change="onPageChange">
        <template #op="{ row }">
          <t-space size="small">
            <t-button theme="primary" variant="text" @click="openEdit(row)">编辑</t-button>
            <t-button theme="danger" variant="text" @click="toggleBan(row)">{{ row.banned ? '解封' : '封禁' }}</t-button>
            <t-button theme="default" variant="text" @click="onResetSecret(row)">重置密钥</t-button>
          </t-space>
        </template>
      </t-table>
    </div>

    <t-dialog v-model:visible="dialogVisible" header="编辑用户" :confirm-btn="{ onClick: submitEdit }" :on-close="resetEditForm">
      <t-form :data="editForm" label-width="96px">
        <t-form-item label="余额(元)">
          <t-input-number v-model="editForm.balance" :min="0" :step="1" />
        </t-form-item>
        <t-form-item label="分组ID">
          <t-input-number v-model="editForm.group_id" :min="0" :step="1" />
        </t-form-item>
        <t-form-item label="套餐ID">
          <t-input-number v-model="editForm.plan_id" :min="0" :step="1" />
        </t-form-item>
        <t-form-item label="到期时间">
          <t-date-picker
            v-model="editForm.expired_at_ms"
            value-type="time-stamp"
            mode="date"
            enable-time-picker
            clearable
            style="width: 100%"
          />
        </t-form-item>
        <t-form-item label="封禁">
          <t-switch v-model="editForm.banned" />
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { MessagePlugin } from 'tdesign-vue-next';
import { listUsers, banUser, updateUser, resetSecret, AdminUserItem, AdminUserUpdateReq } from '@/api/admin';

const list = ref<AdminUserItem[]>([]);
const loading = ref(false);
const keyword = ref('');
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showJumper: true });

const dialogVisible = ref(false);
const editForm = reactive<{ id?: number; balance?: number; group_id?: number; plan_id?: number; expired_at_ms?: number; banned?: boolean }>({});

const columns = [
  { colKey: 'id', title: 'ID', width: 70 },
  { colKey: 'email', title: '邮箱' },
  { colKey: 'balance', title: '余额', cell: (_: unknown, row: AdminUserItem) => `${Number(row.balance) / 100} 元` },
  { colKey: 'commission_balance', title: '佣金', cell: (_: unknown, row: AdminUserItem) => `${Number(row.commission_balance) / 100} 元` },
  { colKey: 'is_admin', title: '管理员', cell: (_: unknown, row: AdminUserItem) => (row.is_admin ? '是' : '否') },
  { colKey: 'banned', title: '状态', cell: (_: unknown, row: AdminUserItem) => (row.banned ? '封禁' : '正常') },
  { colKey: 'plan_id', title: '套餐ID', cell: (_: unknown, row: AdminUserItem) => row.plan_id ?? '-' },
  { colKey: 'expired_at', title: '到期时间', cell: (_: unknown, row: AdminUserItem) => formatTime(row.expired_at) },
  { colKey: 'used_traffic', title: '已用流量', cell: (_: unknown, row: AdminUserItem) => formatBytes(row.used_traffic) },
  { colKey: 'op', title: '操作', width: 260 },
];

async function load() {
  loading.value = true;
  try {
    const res = await listUsers({
      keyword: keyword.value,
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

function openEdit(row: AdminUserItem) {
  editForm.id = row.id;
  editForm.balance = Number(row.balance) / 100;
  editForm.group_id = row.group_id ?? 0;
  editForm.plan_id = row.plan_id ?? 0;
  editForm.expired_at_ms = row.expired_at ? new Date(row.expired_at).getTime() : undefined;
  editForm.banned = row.banned;
  dialogVisible.value = true;
}

function resetEditForm() {
  editForm.id = undefined;
  editForm.balance = undefined;
  editForm.group_id = undefined;
  editForm.plan_id = undefined;
  editForm.expired_at_ms = undefined;
  editForm.banned = false;
}

async function submitEdit() {
  if (!editForm.id) return;
  const data: AdminUserUpdateReq = {
    id: editForm.id,
    balance: Math.round((editForm.balance ?? 0) * 100),
    group_id: editForm.group_id ?? 0,
    plan_id: editForm.plan_id ?? 0,
    expired_at: editForm.expired_at_ms ? Math.floor(editForm.expired_at_ms / 1000) : 0,
    banned: editForm.banned,
  };
  await updateUser(data);
  MessagePlugin.success('已保存');
  dialogVisible.value = false;
  load();
}

async function toggleBan(row: AdminUserItem) {
  await banUser(row.id, !row.banned);
  MessagePlugin.success(row.banned ? '已解封' : '已封禁');
  load();
}

async function onResetSecret(row: AdminUserItem) {
  await resetSecret(row.id);
  MessagePlugin.success('已重置密钥');
  load();
}

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
