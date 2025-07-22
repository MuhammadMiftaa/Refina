export type AreaChartType = {
  date: string;
  physical: number;
  "e-wallet": number;
  bank: number;
  others: number;
};

export type BarChartType = {
  month: string;
  month_name: string;
  total_income: number;
  total_expense: number;
};

export type PieChartType = {
  parent_category_name: string;
  total: number;
  fill: string;
}