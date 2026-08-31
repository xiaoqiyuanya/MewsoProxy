<template>
  <div class="page-container">
    <div class="card">
      <t-card :bordered="false">
        <t-form :data="form" label-width="110px" style="max-width: 560px">
          <t-form-item label="站点名称">
            <t-input v-model="form.app_name" placeholder="MewsoProxy" />
          </t-form-item>
          <t-form-item label="开放注册">
            <t-switch v-model="form.register_enabled" />
          </t-form-item>
          <t-form-item label="订阅地址">
            <t-input v-model="form.subscribe_url" placeholder="http://127.0.0.1:8080" />
          </t-form-item>
          <t-form-item label="后台路径">
            <t-input v-model="form.secure_path" placeholder="/admin" />
          </t-form-item>
          <t-form-item>
            <t-button theme="primary" :loading="saving" @click="submit">保存配置</t-button>
            <t-button theme="default" variant="outline" style="margin-left: 12px" @click="load">刷新</t-button>
          </t-form-item>
        </t-form>
      </t-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { MessagePlugin } from 'tdesign-vue-next';
import { fetchConfig, saveConfig, AdminConfig } from '@/api/admin';

const saving = ref(false);
const form = reactive<{ app_name: string; register_enabled: boolean; subscribe_url: string; secure_path: string }>({
  app_name: '',
  register_enabled: true,
  subscribe_url: '',
  secure_path: '/admin',
});

async function load() {
  const list: AdminConfig[] = await fetchConfig();
  for (const item of list) {
    if (item.key === 'app_name') form.app_name = item.value;
    else if (item.key === 'register_enabled') form.register_enabled = item.value === 'true';
    else if (item.key === 'subscribe_url') form.subscribe_url = item.value;
    else if (item.key === 'secure_path') form.secure_path = item.value;
  }
}

async function submit() {
  saving.value = true;
  try {
    const items: AdminConfig[] = [
      { key: 'app_name', value: form.app_name },
      { key: 'register_enabled', value: String(form.register_enabled) },
      { key: 'subscribe_url', value: form.subscribe_url },
      { key: 'secure_path', value: form.secure_path },
    ];
    for (const item of items) {
      await saveConfig(item);
    }
    MessagePlugin.success('配置已保存');
  } finally {
    saving.value = false;
  }
}

onMounted(load);
</script>
