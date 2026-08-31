<template>
  <div class="page-container">
    <div class="card">
      <h3 style="margin-top: 0">购买套餐</h3>
      <t-loading :loading="loading">
        <t-row :gutter="[16, 16]">
          <t-col v-for="p in plans" :key="p.id" :span="8">
            <t-card :bordered="true">
              <div class="plan-name">{{ p.name }}</div>
              <div class="plan-traffic">{{ formatTraffic(p.transfer_enable) }}</div>
              <t-divider />
              <div class="plan-price">
                <span class="money">{{ formatMoney(p.month_price) }}</span>
                <span class="unit">/月</span>
              </div>
              <t-space direction="vertical" size="small" style="margin-top: 16px">
                <t-button block theme="primary" variant="outline" @click="buy(p, 'month_price')">
                  包月
                </t-button>
                <t-button block theme="primary" variant="outline" @click="buy(p, 'quarter_price')">
                  包季 ({{ formatMoney(p.quarter_price) }})
                </t-button>
                <t-button block theme="primary" @click="buy(p, 'year_price')">
                  包年 ({{ formatMoney(p.year_price) }})
                </t-button>
              </t-space>
            </t-card>
          </t-col>
        </t-row>
      </t-loading>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { MessagePlugin } from 'tdesign-vue-next';
import { listPlans, PlanDTO } from '@/api/plan';
import { createOrder } from '@/api/order';

const plans = ref<PlanDTO[]>([]);
const loading = ref(false);

onMounted(async () => {
  loading.value = true;
  try {
    plans.value = await listPlans();
  } finally {
    loading.value = false;
  }
});

async function buy(plan: PlanDTO, period: string) {
  await createOrder({ plan_id: plan.id, period });
  MessagePlugin.success('订单已创建，请在「订单」页完成支付');
}

function formatMoney(v?: number): string {
  return `${Number(v || 0) / 100}`;
}

function formatTraffic(bytes: number): string {
  return `${(bytes / 1024 / 1024 / 1024).toFixed(2)} GB`;
}
</script>

<style scoped>
.plan-name {
  font-size: 16px;
  font-weight: 600;
}
.plan-traffic {
  color: var(--td-text-color-secondary);
  margin-top: 8px;
}
.plan-price {
  display: flex;
  align-items: baseline;
  gap: 4px;
}
.plan-price .money {
  font-size: 24px;
  color: var(--td-brand-color);
}
.unit {
  color: var(--td-text-color-secondary);
}
</style>
