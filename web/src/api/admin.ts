import request from '@/utils/request';
import type { AuthResp } from '@/api/auth';
import type { PlanDTO } from '@/api/plan';
import type { OrderDTO } from '@/api/order';

export interface AdminLoginReq {
  email: string;
  password: string;
}

export interface AdminUserItem {
  id: number;
  email: string;
  balance: number;
  commission_balance: number;
  is_admin: boolean;
  banned: boolean;
  plan_id?: number;
  group_id?: number;
  expired_at: string;
  transfer_enable: number;
  used_traffic: number;
  created_at: string;
}

export interface AdminUserUpdateReq {
  id: number;
  balance?: number;
  group_id?: number;
  plan_id?: number;
  expired_at?: number;
  banned?: boolean;
}

export interface AdminPlanSaveReq {
  id?: number;
  group_id: number;
  transfer_enable: number;
  name: string;
  speed_limit?: number;
  show: boolean;
  sort?: number;
  renew: boolean;
  content?: string;
  month_price?: number;
  quarter_price?: number;
  half_year_price?: number;
  year_price?: number;
  two_year_price?: number;
  three_year_price?: number;
  onetime_price?: number;
}

export interface AdminConfig {
  key: string;
  value: string;
}

export interface AdminSystemStatus {
  server_time: number;
  db_status: string;
  redis_status: string;
  online_user_count: number;
  user_count: number;
  order_count: number;
  today_paid_total: number;
}

export interface ServerGroup {
  id: number;
  name: string;
  created_at: number;
  updated_at: number;
}

export interface AdminServerGroupSaveReq {
  id?: number;
  name: string;
}

export interface ServerNode {
  id: number;
  group_id: string;
  route_id?: string;
  parent_id?: number;
  tags?: string;
  name: string;
  rate: string;
  host: string;
  port: string;
  server_port: number;
  show: boolean;
  sort?: number;
  created_at: number;
  updated_at: number;
  allow_insecure?: boolean;
  server_name?: string;
  tls?: number;
  network?: string;
  cipher?: string;
  up_mbps?: number;
  down_mbps?: number;
  insecure?: boolean;
}

export interface AdminServerNodeSaveReq {
  id?: number;
  type: string;
  group_id: string;
  name: string;
  host: string;
  port: string;
  server_port: number;
  rate?: string;
  show: boolean;
  sort?: number;
  tags?: string;
  route_id?: string;
  parent_id?: number;
  cipher?: string;
  network?: string;
  tls?: number;
  allow_insecure?: boolean;
  server_name?: string;
  up_mbps?: number;
  down_mbps?: number;
  insecure?: boolean;
}

export interface Coupon {
  id: number;
  code: string;
  name: string;
  type: number;
  value: number;
  show: boolean;
  limit_use?: number;
  limit_use_with_user?: number;
  limit_plan_ids?: string;
  limit_period?: string;
  started_at: number;
  ended_at: number;
  created_at: number;
  updated_at: number;
}

export interface AdminCouponSaveReq {
  id?: number;
  code: string;
  name: string;
  type: number;
  value: number;
  show: boolean;
  limit_use?: number;
  limit_use_with_user?: number;
  limit_plan_ids?: string;
  limit_period?: string;
  started_at: number;
  ended_at: number;
}

export interface Notice {
  id: number;
  title: string;
  content: string;
  show: boolean;
  img_url?: string;
  tags?: string;
  created_at: number;
  updated_at: number;
}

export interface AdminNoticeSaveReq {
  id?: number;
  title: string;
  content: string;
  show: boolean;
  img_url?: string;
  tags?: string;
}

export interface Payment {
  id: number;
  uuid: string;
  payment: string;
  name: string;
  icon?: string;
  config: string;
  notify_domain?: string;
  handling_fee_fixed?: number;
  handling_fee_percent?: number;
  enable: boolean;
  sort?: number;
  created_at: number;
  updated_at: number;
}

export interface AdminPaymentSaveReq {
  id?: number;
  uuid: string;
  payment: string;
  name: string;
  icon?: string;
  config: string;
  notify_domain?: string;
  handling_fee_fixed?: number;
  handling_fee_percent?: number;
  enable: boolean;
  sort?: number;
}

export interface AdminListResp<T> {
  list: T[];
  total: number;
}

export function adminLogin(data: AdminLoginReq) {
  return request.post<AuthResp, AuthResp>('/admin/login', data);
}

export function fetchConfig() {
  return request.get<AdminConfig[], AdminConfig[]>('/admin/config/fetch');
}

export function saveConfig(data: AdminConfig) {
  return request.post<null, null>('/admin/config/save', data);
}

export function fetchSystemStatus() {
  return request.get<AdminSystemStatus, AdminSystemStatus>('/admin/system/status');
}

export function listPlansAdmin() {
  return request.get<PlanDTO[], PlanDTO[]>('/admin/plan/list');
}

export function savePlan(data: AdminPlanSaveReq) {
  return request.post<{ id: number }, { id: number }>('/admin/plan/save', data);
}

export function dropPlan(id: number) {
  return request.post<null, null>('/admin/plan/drop', { id });
}

export function listUsers(params?: { keyword?: string; page?: number; page_size?: number }) {
  return request.get<AdminListResp<AdminUserItem>, AdminListResp<AdminUserItem>>('/admin/user/list', { params });
}

export function getUser(id: number) {
  return request.post<AdminUserItem, AdminUserItem>('/admin/user/info', { id });
}

export function updateUser(data: AdminUserUpdateReq) {
  return request.post<null, null>('/admin/user/update', data);
}

export function banUser(id: number, ban: boolean) {
  return request.post<null, null>('/admin/user/ban', { id, ban });
}

export function resetSecret(id: number) {
  return request.post<null, null>('/admin/user/reset_secret', { id });
}

export function listOrdersAdmin(params?: { keyword?: string; status?: number; page?: number; page_size?: number }) {
  return request.get<AdminListResp<OrderDTO>, AdminListResp<OrderDTO>>('/admin/order/list', { params });
}

export function getOrder(id: number) {
  return request.get<OrderDTO, OrderDTO>('/admin/order/info', { params: { id } });
}

export function cancelOrder(id: number) {
  return request.post<null, null>('/admin/order/cancel', { id });
}

export function markOrderPaid(id: number, callback_no?: string) {
  return request.post<null, null>('/admin/order/paid', { id, callback_no });
}

export function listGroups() {
  return request.get<ServerGroup[], ServerGroup[]>('/admin/server/group/list');
}

export function saveGroup(data: AdminServerGroupSaveReq) {
  return request.post<{ id: number }, { id: number }>('/admin/server/group/save', data);
}

export function dropGroup(id: number) {
  return request.post<null, null>('/admin/server/group/drop', { id });
}

export function listNodes(type: string) {
  return request.get<ServerNode[], ServerNode[]>('/admin/server/node/list', { params: { type } });
}

export function saveNode(data: AdminServerNodeSaveReq) {
  return request.post<{ id: number }, { id: number }>('/admin/server/node/save', data);
}

export function dropNode(type: string, id: number) {
  return request.post<null, null>('/admin/server/node/drop', { type, id });
}

export function listCoupons() {
  return request.get<Coupon[], Coupon[]>('/admin/coupon/list');
}

export function saveCoupon(data: AdminCouponSaveReq) {
  return request.post<{ id: number }, { id: number }>('/admin/coupon/save', data);
}

export function dropCoupon(id: number) {
  return request.post<null, null>('/admin/coupon/drop', { id });
}

export function toggleCouponShow(id: number, show: boolean) {
  return request.post<null, null>('/admin/coupon/show', { id, show });
}

export function listNotices() {
  return request.get<Notice[], Notice[]>('/admin/notice/list');
}

export function saveNotice(data: AdminNoticeSaveReq) {
  return request.post<{ id: number }, { id: number }>('/admin/notice/save', data);
}

export function dropNotice(id: number) {
  return request.post<null, null>('/admin/notice/drop', { id });
}

export function toggleNoticeShow(id: number, show: boolean) {
  return request.post<null, null>('/admin/notice/show', { id, show });
}

export function listPayments() {
  return request.get<Payment[], Payment[]>('/admin/payment/list');
}

export function savePayment(data: AdminPaymentSaveReq) {
  return request.post<{ id: number }, { id: number }>('/admin/payment/save', data);
}

export function dropPayment(id: number) {
  return request.post<null, null>('/admin/payment/drop', { id });
}

export function togglePaymentShow(id: number, enable: boolean) {
  return request.post<null, null>('/admin/payment/show', { id, enable });
}
