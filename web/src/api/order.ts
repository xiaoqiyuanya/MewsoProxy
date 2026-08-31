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

export function createOrder(data: CreateOrderReq) {
  return request.post<OrderDTO, OrderDTO>('/order/create', data);
}

export function listOrders(params?: { page?: number; page_size?: number }) {
  return request.get<{ list: OrderDTO[]; total: number }, { list: OrderDTO[]; total: number }>(
    '/order/list',
    { params },
  );
}
