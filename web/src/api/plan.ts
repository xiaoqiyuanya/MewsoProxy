import request from '@/utils/request';

export interface PlanDTO {
  id: number;
  name: string;
  group_id: number;
  transfer_enable: number;
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

export function listPlans() {
  return request.get<PlanDTO[], PlanDTO[]>('/plan/list');
}
