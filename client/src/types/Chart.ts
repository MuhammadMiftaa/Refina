export type AreaChartType = {
  date: string;
  bank: number;
  "e-wallet": number;
  cash: number;
  others: number;
};

export type BarChartType = {
  month: string;
  income: number;
  expense: number;
};

export type PieChartType = {
  category: string;
  value: number;
  fill: string;
}