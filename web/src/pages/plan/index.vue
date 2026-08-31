<template>
  <div class="page-container">
    <div class="page-title">
      <div>
        <h2>购买套餐</h2>
        <p class="desc">选择适合你的套餐，订阅即刻生效</p>
      </div>
    </div>

    <t-loading :loading="loading">
      <div class="plan-grid">
        <div v-for="p in plans" :key="p.id" class="plan-card card">
          <div class="plan-name">{{ p.name }}</div>
          <div class="plan-traffic">{{ formatTraffic(p.transfer_enable) }}</div>
          <div class="plan-price">
            <span class="plan-symbol">¥</span>
            <span class="plan-amount">{{ formatMoney(p.month_price) }}</span>
            <span class="plan-unit">/月</span>
          </div>
          <t-divider />
          <div class="plan-buttons">
            <t-button
              v-for="opt in periods(p)"
              :key="opt.key"
              block
              :theme="opt.featured ? 'primary' : 'default'"
              :variant="opt.featured ? 'solid' : 'outline'"
              @click="buy(p, opt.key)"
            >
              {{ opt.label }} · {{ formatMoney(opt.price) }} 元
            </t-button>
          </div>
        </div>
      </div>
      <t-empty v-if="!loading && plans.length === 0" description="暂无可用套餐" />
    </t-loading>

    <PayDialog v-model:visible="payVisible" :order="payOrder" @paid="onPaid" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { listPlans, PlanDTO } from '@/api/plan';
import { createOrder, OrderDTO } from '@/api/order';
import PayDialog from '@/components/PayDialog.vue';

interface PeriodOpt {
  key: string;
  label: string;
  price?: number;
  featured: boolean;
}

const PERIOD_LABELS: Record<string, string> = {
  month_price: '包月',
  quarter_price: '包季',
  half_year_price: '半年',
  year_price: '包年',
  two_year_price: '两年',
  three_year_price: '三年',
  onetime_price: '一次性',
};

const PERIOD_ORDER = [
  'month_price',
  'quarter_price',
  'half_year_price',
  'year_price',
  'two_year_price',
  'three_year_price',
  'onetime_price',
];

const plans = ref<PlanDTO[]>([]);
const loading = ref(false);
const payVisible = ref(false);
const payOrder = ref<OrderDTO | null>(null);

onMounted(async () => {
  loading.value = true;
  try {
    plans.value = await listPlans();
  } finally {
    loading.value = false;
  }
});

function periods(p: PlanDTO): PeriodOpt[] {
  const list: PeriodOpt[] = [];
  for (const key of PERIOD_ORDER) {
    const price = (p as unknown as Record<string, number | undefined>)[key];
    if (typeof price === 'number' && price > 0) {
      list.push({ key, label: PERIOD_LABELS[key] || key, price, featured: false });
    }
  }
  if (list.length > 0) {
    // 突出性价比最高的包年选项
    const year = list.find((o) => o.key === 'year_price');
    const featured = year || list[list.length - 1];
    featured.featured = true;
  }
  return list;
}

async function buy(plan: PlanDTO, period: string) {
  const order = await createOrder({ plan_id: plan.id, period });
  payOrder.value = order;
  payVisible.value = true;
}

function onPaid() {
  // 支付完成后无需额外操作
}

function formatMoney(v?: number): string {
  return (Number(v || 0) / 100).toFixed(2);
}

function formatTraffic(bytes: number): string {
  return `${(bytes / 1024 / 1024 / 1024).toFixed(2)} GB`;
}
</script>

<style scoped>
.plan-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 16px;
}
.plan-card {
  display: flex;
  flex-direction: column;
  padding: 24px;
  transition: box-shadow 0.2s ease, transform 0.2s ease;
}
.plan-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}
.plan-name {
  font-size: 16px;
  font-weight: 600;
  color: #1f2329;
}
.plan-traffic {
  font-size: 13px;
  color: #8a919f;
  margin-top: 6px;
}
.plan-price {
  display: flex;
  align-items: baseline;
  gap: 4px;
  margin-top: 12px;
}
.plan-symbol {
  font-size: 16px;
  color: var(--td-brand-color);
}
.plan-amount {
  font-size: 30px;
  font-weight: 700;
  color: var(--td-brand-color);
  line-height: 1;
}
.plan-unit {
  font-size: 13px;
  color: #8a919f;
}
.plan-buttons {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 4px;
}
</style>
