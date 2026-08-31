<template>
  <div class="page-container">
    <div class="card">
      <div class="table-head">
        <h3 style="margin: 0">公告管理</h3>
        <t-space size="small">
          <t-button theme="primary" @click="openCreate">新增公告</t-button>
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

    <t-dialog v-model:visible="dialogVisible" :header="form.id ? '编辑公告' : '新增公告'" width="560px" :confirm-btn="{ onClick: submit }" :on-close="resetForm">
      <t-form :data="form" label-width="80px">
        <t-form-item label="标题" required>
          <t-input v-model="form.title" placeholder="公告标题" />
        </t-form-item>
        <t-form-item label="内容" required>
          <t-textarea v-model="form.content" placeholder="公告内容" :autosize="{ minRows: 4 }" />
        </t-form-item>
        <t-form-item label="图片">
          <t-input v-model="form.img_url" placeholder="图片 URL" />
        </t-form-item>
        <t-form-item label="标签">
          <t-input v-model="form.tags" placeholder="标签，逗号分隔" />
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
import { listNotices, saveNotice, dropNotice, toggleNoticeShow, Notice, AdminNoticeSaveReq } from '@/api/admin';

const list = ref<Notice[]>([]);
const loading = ref(false);
const dialogVisible = ref(false);
const form = reactive<AdminNoticeSaveReq>({
  id: undefined,
  title: '',
  content: '',
  show: true,
  img_url: '',
  tags: '',
});

const columns = [
  { colKey: 'id', title: 'ID', width: 70 },
  { colKey: 'title', title: '标题' },
  { colKey: 'tags', title: '标签', cell: (_: unknown, row: Notice) => row.tags ?? '-' },
  { colKey: 'show', title: '展示', cell: (_: unknown, row: Notice) => (row.show ? '是' : '否') },
  { colKey: 'created_at', title: '创建时间', cell: (_: unknown, row: Notice) => formatTime(row.created_at) },
  { colKey: 'op', title: '操作', width: 200 },
];

async function load() {
  loading.value = true;
  try {
    list.value = await listNotices();
  } finally {
    loading.value = false;
  }
}

function resetForm() {
  Object.assign(form, { id: undefined, title: '', content: '', show: true, img_url: '', tags: '' });
}

function openCreate() {
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: Notice) {
  resetForm();
  form.id = row.id;
  form.title = row.title;
  form.content = row.content;
  form.show = row.show;
  form.img_url = row.img_url ?? '';
  form.tags = row.tags ?? '';
  dialogVisible.value = true;
}

async function submit() {
  if (!form.title || !form.content) {
    MessagePlugin.error('请填写标题和内容');
    return;
  }
  const data: AdminNoticeSaveReq = {
    id: form.id,
    title: form.title,
    content: form.content,
    show: form.show,
    img_url: form.img_url || undefined,
    tags: form.tags || undefined,
  };
  await saveNotice(data);
  MessagePlugin.success('已保存');
  dialogVisible.value = false;
  load();
}

async function toggleShow(row: Notice) {
  await toggleNoticeShow(row.id, !row.show);
  MessagePlugin.success('已更新');
  load();
}

async function onDrop(row: Notice) {
  await dropNotice(row.id);
  MessagePlugin.success('已删除');
  load();
}

function formatTime(v?: number): string {
  if (!v) return '-';
  return new Date(v * 1000).toLocaleString();
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
