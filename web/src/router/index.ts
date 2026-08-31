import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
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
  if (!to.meta.public && !authed) {
    return { path: '/login', query: { redirect: to.fullPath } };
  }
  if (to.meta.public && authed) {
    return { path: '/dashboard' };
  }
  if (authed) {
    const store = useUserStore();
    if (!store.user) {
      try {
        await store.fetchMe();
      } catch {
        store.reset();
        return { path: '/login' };
      }
    }
  }
  return true;
});

router.afterEach((to) => {
  const title = to.meta.title as string | undefined;
  document.title = title ? `${title} - MewsoProxy` : 'MewsoProxy';
});

export default router;
