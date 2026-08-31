<template>
  <div class="auth-wrap">
    <t-card class="auth-card" :bordered="false">
      <div class="auth-head">
        <h2>MewsoProxy</h2>
        <p>注册新账号</p>
      </div>
      <t-form :data="form" :rules="rules" @submit="onSubmit">
        <t-form-item label="邮箱" name="email">
          <t-input v-model="form.email" placeholder="请输入邮箱" size="large" />
        </t-form-item>
        <t-form-item label="密码" name="password">
          <t-input v-model="form.password" type="password" placeholder="至少 8 位" size="large" />
        </t-form-item>
        <t-form-item label="邀请码" name="invite_code">
          <t-input v-model="form.invite_code" placeholder="可选" size="large" />
        </t-form-item>
        <t-form-item>
          <t-button type="submit" block size="large" :loading="loading">注册</t-button>
        </t-form-item>
      </t-form>
      <div class="auth-foot">
        已有账号？<t-link theme="primary" @click="$router.push('/login')">去登录</t-link>
      </div>
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import { MessagePlugin } from 'tdesign-vue-next';
import { register, RegisterReq } from '@/api/auth';
import { useUserStore } from '@/store/user';

const router = useRouter();
const store = useUserStore();

const loading = ref(false);
const form = reactive<RegisterReq>({ email: '', password: '', invite_code: '' });
const rules = {
  email: [{ required: true, message: '请输入邮箱', type: 'error' as const }, { email: true, message: '邮箱格式不正确' }],
  password: [
    { required: true, message: '请输入密码' },
    { min: 8, message: '密码至少 8 位' },
  ],
};

async function onSubmit({ validateResult }: { validateResult?: boolean }) {
  if (validateResult !== true) return;
  loading.value = true;
  try {
    const res = await register(form);
    store.saveToken(res.token.access_token);
    await store.fetchMe();
    MessagePlugin.success('注册成功');
    router.replace('/dashboard');
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
