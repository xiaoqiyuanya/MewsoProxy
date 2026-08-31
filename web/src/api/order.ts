import request from '@/utils/request';

export interface OrderDTO {
  id: number;
  plan_id: number;
  period: string;
  trade_no: string;
  type: number;
  total_amount: number;
  balance_amount: number;
  discount_amount: number;
  status: number;
  created_at: string;
  paid_at: string;
}

export interface CreateOrderReq {
  plan_id: number;
  period: string;
  coupon_code?: string;
  payment_id?: number;
  use_balance?: boolean;
}

export interface PaymentChannel {
  id: number;
  payment: string;
  name: string;
  icon?: string;
}

export interface PayResult {
  channel: string;
  pay_type: 'redirect' | 'qrcode' | 'completed';
  url?: string;
  qr_code?: string;
  completed?: boolean;
}

export function createOrder(data: CreateOrderReq) {
  return request.post<OrderDTO, OrderDTO>('/order/create', data);
}

export function listOrders(params?: { page?: number; page_size?: number }) {
  return request.get<{ list: OrderDTO[]; total: number }, { list: OrderDTO[]; total: number }>(
    '/order/list',
    { params },
  );
}

export function detailOrder(id: number) {
  return request.post<OrderDTO, OrderDTO>('/order/detail', { id });
}

export function payChannels() {
  return request.get<PaymentChannel, PaymentChannel[]>('/payment/channels');
}

export function createPayment(data: { order_id: number; payment_id?: number }) {
  return request.post<PayResult, PayResult>('/payment/create', data);
}

export function notifyPaid(data: { trade_no: string; callback_no: string }) {
  return request.post<unknown, unknown>('/payment/notify', data);
}
