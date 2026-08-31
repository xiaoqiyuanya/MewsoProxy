<template>
  <t-dialog :visible="visible" header="支付订单" width="480px" :footer="false" @close="onClose">
    <div v-if="!result" class="pay-form">
      <div class="order-no">订单号：{{ order?.trade_no }}</div>
      <div class="order-amount">
        应付金额：<span class="money">{{ formatMoney(order?.total_amount) }}</span> 元
      </div>
      <t-loading :loading="loadingChannels" class="channel-loading">
        <t-radio-group v-if="channels.length > 1" v-model="payment_id" class="channels">
          <t-radio-button v-for="c in channels" :key="c.id" :value="c.id">{{ c.name }}</t-radio-button>
        </t-radio-group>
        <div v-else-if="channels.length === 1" class="single-channel">
          {{ channels[0].name }}
        </div>
        <div v-else class="empty">暂无可用支付渠道，请联系管理员</div>
      </t-loading>
      <t-button
        block
        theme="primary"
        :loading="loadingPay"
        :disabled="channels.length === 0"
        @click="goPay"
      >
        确认支付
      </t-button>
    </div>

    <div v-else class="pay-result">
      <template v-if="result.pay_type === 'completed'">
        <t-result theme="success" title="支付成功" />
      </template>
      <template v-else-if="result.channel === 'mock'">
        <div class="mock-qr">{{ result.qr_code }}</div>
        <div class="mock-hint">该渠道为模拟支付，点击下方按钮模拟支付成功</div>
        <t-button block theme="primary" :loading="paying" @click="mockPaid">模拟支付成功</t-button>
      </template>
      <template v-else>
        <t-result theme="success" title="请在新窗口完成支付" />
        <t-button block theme="primary" @click="openPay">打开支付页面</t-button>
      </template>
      <div v-if="polling" class="polling">等待支付结果…</div>
    </div>
  </t-dialog>
</template>

<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from 'vue';
import { MessagePlugin } from 'tdesign-vue-next';
import { OrderDTO, PaymentChannel, PayResult, payChannels, createPayment, notifyPaid, detailOrder } from '@/api/order';

const props = defineProps<{
  visible: boolean;
  order?: OrderDTO | null;
}>();

const emit = defineEmits<{
  (e: 'update:visible', v: boolean): void;
  (e: 'paid'): void;
}>();

const channels = ref<PaymentChannel[]>([]);
const loadingChannels = ref(false);
const loadingPay = ref(false);
const paying = ref(false);
const payment_id = ref<number>();
const result = ref<PayResult | null>(null);
const polling = ref(false);
let timer: ReturnType<typeof setInterval> | null = null;

watch(
  () => props.visible,
  (v) => {
    if (v) {
      reset();
      loadChannels();
    }
  },
);

watch(
  () => props.order?.id,
  () => {
    reset();
  },
);

function reset() {
  result.value = null;
  paying.value = false;
  payment_id.value = undefined;
  stopPoll();
}

function onClose() {
  close();
}

function close() {
  stopPoll();
  emit('update:visible', false);
}

async function loadChannels() {
  loadingChannels.value = true;
  try {
    channels.value = await payChannels();
    if (channels.value.length > 0) {
      payment_id.value = channels.value[0].id;
    }
  } catch (e) {
    channels.value = [];
  } finally {
    loadingChannels.value = false;
  }
}

async function goPay() {
  if (!props.order || !payment_id.value) {
    MessagePlugin.error('请选择支付渠道');
    return;
  }
  loadingPay.value = true;
  try {
    const res = await createPayment({ order_id: props.order.id, payment_id: payment_id.value });
    result.value = res;
    if (res.pay_type === 'completed') {
      MessagePlugin.success('支付成功');
      emit('paid');
      setTimeout(close, 1500);
      return;
    }
    startPoll();
  } catch (e) {
    // 请求层已提示
  } finally {
    loadingPay.value = false;
  }
}

async function mockPaid() {
  if (!props.order) return;
  paying.value = true;
  try {
    await notifyPaid({ trade_no: props.order.trade_no, callback_no: `MOCK-${Date.now()}` });
    MessagePlugin.success('支付成功');
    emit('paid');
    stopPoll();
    setTimeout(close, 1500);
  } catch (e) {
    // 请求层已提示
  } finally {
    paying.value = false;
  }
}

function openPay() {
  if (result.value?.url) {
    window.open(result.value.url, '_blank');
  }
}

function startPoll() {
  stopPoll();
  polling.value = true;
  timer = setInterval(async () => {
    if (!props.order) return;
    try {
      const o = await detailOrder(props.order.id);
      if (o.status === 3) {
        MessagePlugin.success('支付成功');
        emit('paid');
        stopPoll();
        setTimeout(close, 1500);
      }
    } catch (e) {
      // 忽略单次轮询错误
    }
  }, 3000);
}

function stopPoll() {
  polling.value = false;
  if (timer) {
    clearInterval(timer);
    timer = null;
  }
}

function formatMoney(v?: number): string {
  return (Number(v || 0) / 100).toFixed(2);
}

onBeforeUnmount(stopPoll);
</script>

<style scoped>
.order-no {
  color: var(--td-text-color-secondary);
  font-size: 13px;
  margin-bottom: 8px;
}
.order-amount {
  font-size: 16px;
  margin-bottom: 20px;
}
.order-amount .money {
  color: var(--td-brand-color);
  font-size: 24px;
  font-weight: 600;
}
.channel-loading {
  margin-bottom: 20px;
  min-height: 40px;
}
.channels {
  display: block;
}
.single-channel,
.empty {
  padding: 8px 0;
  color: var(--td-text-color-secondary);
}
.mock-qr {
  font-family: var(--td-font-family-mono, monospace);
  padding: 16px;
  border: 1px dashed var(--td-border-level-1-color);
  border-radius: 8px;
  text-align: center;
  word-break: break-all;
  font-size: 12px;
  margin-bottom: 12px;
  background: var(--td-bg-color-container-hover);
}
.mock-hint {
  color: var(--td-text-color-secondary);
  font-size: 13px;
  margin-bottom: 16px;
}
.polling {
  margin-top: 16px;
  text-align: center;
  color: var(--td-text-color-secondary);
  font-size: 13px;
}
</style>
