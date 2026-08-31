import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import { MessagePlugin } from 'tdesign-vue-next';
import { getToken } from '@/utils/token';
import { useUserStore } from '@/store/user';

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/login/index.vue'),
    meta: { title: '登录', public: true },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/pages/register/index.vue'),
    meta: { title: '注册', public: true },
  },
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: () => import('@/pages/admin/login/index.vue'),
    meta: { title: '管理员登录', public: true },
  },
  {
    path: '/',
    component: () => import('@/layouts/BasicLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/pages/dashboard/index.vue'),
        meta: { title: '概览' },
      },
      {
        path: 'plan',
        name: 'Plan',
        component: () => import('@/pages/plan/index.vue'),
        meta: { title: '套餐' },
      },
      {
        path: 'order',
        name: 'Order',
        component: () => import('@/pages/order/index.vue'),
        meta: { title: '订单' },
      },
    ],
  },
  {
    path: '/admin',
    component: () => import('@/layouts/AdminLayout.vue'),
    redirect: '/admin/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'AdminDashboard',
        component: () => import('@/pages/admin/dashboard/index.vue'),
        meta: { title: '概览', admin: true },
      },
      {
        path: 'user',
        name: 'AdminUser',
        component: () => import('@/pages/admin/user/index.vue'),
        meta: { title: '用户', admin: true },
      },
      {
        path: 'plan',
        name: 'AdminPlan',
        component: () => import('@/pages/admin/plan/index.vue'),
        meta: { title: '套餐', admin: true },
      },
      {
        path: 'order',
        name: 'AdminOrder',
        component: () => import('@/pages/admin/order/index.vue'),
        meta: { title: '订单', admin: true },
      },
      {
        path: 'server',
        name: 'AdminServer',
        component: () => import('@/pages/admin/server/index.vue'),
        meta: { title: '节点', admin: true },
      },
      {
        path: 'coupon',
        name: 'AdminCoupon',
        component: () => import('@/pages/admin/coupon/index.vue'),
        meta: { title: '优惠券', admin: true },
      },
      {
        path: 'notice',
        name: 'AdminNotice',
        component: () => import('@/pages/admin/notice/index.vue'),
        meta: { title: '公告', admin: true },
      },
      {
        path: 'payment',
        name: 'AdminPayment',
        component: () => import('@/pages/admin/payment/index.vue'),
        meta: { title: '支付', admin: true },
      },
      {
        path: 'config',
        name: 'AdminConfig',
        component: () => import('@/pages/admin/config/index.vue'),
        meta: { title: '配置', admin: true },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/dashboard',
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach(async (to) => {
  const authed = !!getToken();
  const isAdminArea = to.path.startsWith('/admin');
  const isAdminLogin = to.path === '/admin/login';

  const store = useUserStore();

  if (authed && !store.user) {
    try {
      await store.fetchMe();
    } catch {
      store.reset();
      return { path: isAdminArea ? '/admin/login' : '/login' };
    }
  }

  if (isAdminArea) {
    if (isAdminLogin) {
      if (authed && store.isAdmin) return { path: '/admin/dashboard' };
      if (authed && !store.isAdmin) return { path: '/dashboard' };
      return true;
    }
    if (!authed) return { path: '/admin/login', query: { redirect: to.fullPath } };
    if (!store.isAdmin) {
      MessagePlugin.error('无权访问管理后台');
      return { path: '/dashboard' };
    }
    return true;
  }

  if (!to.meta.public && !authed) {
    return { path: '/login', query: { redirect: to.fullPath } };
  }
  if (to.meta.public && authed) {
    return { path: '/dashboard' };
  }
  return true;
});

router.afterEach((to) => {
  const title = to.meta.title as string | undefined;
  document.title = title ? `${title} - MewsoProxy` : 'MewsoProxy';
});

export default router;
