<template>
  <div class="auth-wrap">
    <t-card class="auth-card" :bordered="false">
      <div class="auth-head">
        <h2>MewsoProxy</h2>
        <p>管理员登录</p>
      </div>
      <t-form :data="form" :rules="rules" @submit="onSubmit">
        <t-form-item label="邮箱" name="email">
          <t-input v-model="form.email" placeholder="请输入管理员邮箱" size="large" />
        </t-form-item>
        <t-form-item label="密码" name="password">
          <t-input v-model="form.password" type="password" placeholder="请输入密码" size="large" />
        </t-form-item>
        <t-form-item>
          <t-button type="submit" block size="large" :loading="loading">登录</t-button>
        </t-form-item>
      </t-form>
      <div class="auth-foot">
        <t-link theme="default" @click="$router.push('/login')">返回用户端登录</t-link>
      </div>
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { MessagePlugin } from 'tdesign-vue-next';
import { adminLogin, AdminLoginReq } from '@/api/admin';
import { useUserStore } from '@/store/user';

const router = useRouter();
const route = useRoute();
const store = useUserStore();

const loading = ref(false);
const form = reactive<AdminLoginReq>({ email: '', password: '' });
const rules = {
  email: [{ required: true, message: '请输入邮箱', type: 'error' as const }, { email: true, message: '邮箱格式不正确' }],
  password: [{ required: true, message: '请输入密码' }],
};

async function onSubmit({ validateResult }: { validateResult?: boolean }) {
  if (validateResult !== true) return;
  loading.value = true;
  try {
    const res = await adminLogin(form);
    store.saveToken(res.token.access_token);
    await store.fetchMe();
    MessagePlugin.success('登录成功');
    router.replace((route.query.redirect as string) || '/admin/dashboard');
  } catch {
    /* 由拦截器提示 */
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.auth-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  background: linear-gradient(135deg, #f3f4f6 0%, #eef2ff 100%);
}
.auth-card {
  width: 400px;
  padding: 12px;
}
.auth-head {
  text-align: center;
  margin-bottom: 24px;
}
.auth-head h2 {
  margin: 0;
  color: var(--td-brand-color);
}
.auth-head p {
  color: var(--td-text-color-secondary);
  margin-top: 8px;
}
.auth-foot {
  text-align: center;
  margin-top: 8px;
}
</style>
