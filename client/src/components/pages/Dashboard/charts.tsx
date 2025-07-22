"use client";

import { useState } from "react";
import {
  Area,
  AreaChart,
  Bar,
  BarChart,
  CartesianGrid,
  Cell,
  Pie,
  PieChart,
  XAxis,
} from "recharts";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { AreaChartType, BarChartType, PieChartType } from "@/types/Chart";
import { TrendingUp } from "lucide-react";
import { FaRegChartBar } from "react-icons/fa6";
import { getLast6MonthsRange } from "@/helper/Helper";

export function ChartArea({
  chartData,
  chartConfig,
}: {
  chartData: AreaChartType[];
  chartConfig: ChartConfig;
}) {
  const [timeRange, setTimeRange] = useState("90d");

  const filteredData = chartData.filter((item) => {
    const date = new Date(item.date);
    const referenceDate = new Date("2024-06-30");
    let daysToSubtract = 90;
    if (timeRange === "30d") {
      daysToSubtract = 30;
    } else if (timeRange === "7d") {
      daysToSubtract = 7;
    }
    const startDate = new Date(referenceDate);
    startDate.setDate(startDate.getDate() - daysToSubtract);
    return date >= startDate;
  });

  return (
    <Card className="pt-0">
      <CardHeader className="flex items-center gap-2 space-y-0 border-b py-5 sm:flex-row">
        <div className="grid flex-1 gap-1">
          <CardTitle>Daily Wallet Balance Overview</CardTitle>
          <CardDescription>
            A 3 months overview of your wallet balances across different sources
            including physical cash, bank accounts, e-wallets, and others.
          </CardDescription>
        </div>
        <Select value={timeRange} onValueChange={setTimeRange} disabled>
          <SelectTrigger
            disabled
            className="hidden w-[160px] rounded-lg sm:ml-auto sm:flex"
            aria-label="Select a value"
          >
            <SelectValue placeholder="Last 3 months" />
          </SelectTrigger>
          <SelectContent className="rounded-xl">
            <SelectItem value="90d" className="rounded-lg">
              Last 3 months
            </SelectItem>
            <SelectItem value="30d" className="rounded-lg">
              Last 30 days
            </SelectItem>
            <SelectItem value="7d" className="rounded-lg">
              Last 7 days
            </SelectItem>
          </SelectContent>
        </Select>
      </CardHeader>
      <CardContent className="px-2 pt-4 sm:px-6 sm:pt-6">
        <ChartContainer
          config={chartConfig}
          className="aspect-auto h-[250px] w-full"
        >
          <AreaChart data={filteredData}>
            <defs>
              {Object.entries(chartConfig).map(
                ([key, config]) =>
                  key !== "balance" && (
                    <linearGradient
                      key={key}
                      id={`fill${key.charAt(0).toUpperCase() + key.slice(1)}`}
                      x1="0"
                      y1="0"
                      x2="0"
                      y2="1"
                    >
                      <stop
                        offset="5%"
                        stopColor={config.stopColor}
                        stopOpacity={0.8}
                      />
                      <stop
                        offset="95%"
                        stopColor={config.stopColor}
                        stopOpacity={0.1}
                      />
                    </linearGradient>
                  ),
              )}
            </defs>
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="date"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              minTickGap={32}
              tickFormatter={(value) => {
                const date = new Date(value);
                return date.toLocaleDateString("en-US", {
                  month: "short",
                  day: "numeric",
                });
              }}
            />
            <ChartTooltip
              cursor={false}
              content={
                <ChartTooltipContent
                  labelFormatter={(value) => {
                    return new Date(value).toLocaleDateString("en-US", {
                      month: "short",
                      day: "numeric",
                    });
                  }}
                  indicator="dot"
                />
              }
            />
            {Object.entries(chartConfig).map(
              ([key, config]) =>
                key !== "balance" && (
                  <Area
                    key={key}
                    dataKey={key}
                    type="natural"
                    fill={config.fill}
                    stroke={config.stroke}
                    stackId="a"
                  />
                ),
            )}
            <ChartLegend content={<ChartLegendContent />} />
          </AreaChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}

export function ChartBar({
  chartData,
  chartConfig,
}: {
  chartData: BarChartType[];
  chartConfig: ChartConfig;
}) {
  const chartDataConverted = chartData.map((item) => ({
    ...item,
    month_name: item.month_name
      .split("")[0]
      .toUpperCase()
      .concat(item.month_name.slice(1)),
  }));

  const maxMonthIncome =
    chartData
      .reduce((max, curr) =>
        curr.total_income > max.total_income ? curr : max,
      )
      .month_name.split("")[0]
      .toUpperCase() +
    chartData
      .reduce((max, curr) =>
        curr.total_income > max.total_income ? curr : max,
      )
      .month_name.slice(1);

  const maxMonthExpense =
    chartData
      .reduce((max, curr) =>
        curr.total_expense > max.total_expense ? curr : max,
      )
      .month_name.split("")[0]
      .toUpperCase() +
    chartData
      .reduce((max, curr) =>
        curr.total_expense > max.total_expense ? curr : max,
      )
      .month_name.slice(1);

  return (
    <Card className="h-full">
      <CardHeader>
        <CardTitle>Monthly Income vs Expenses</CardTitle>
        <CardDescription>
          Total income and spending tracked per month for the past 6 months to
          help monitor your financial flow and trends.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <BarChart accessibilityLayer data={chartDataConverted}>
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="month_name"
              tickLine={false}
              tickMargin={10}
              axisLine={false}
              tickFormatter={(value) => value.slice(0, 3)}
            />
            <ChartTooltip
              cursor={false}
              content={<ChartTooltipContent indicator="dashed" />}
            />
            {Object.entries(chartConfig).map(([key, config]) => (
              <Bar key={key} dataKey={key} fill={config.color} radius={4} />
            ))}
          </BarChart>
        </ChartContainer>
      </CardContent>
      <CardFooter className="flex-col items-start gap-2 text-sm">
        <div className="flex items-center gap-2 leading-none font-medium">
          Spending peaked in {maxMonthExpense}, while income was highest in{" "}
          {maxMonthIncome}. <FaRegChartBar className="h-4 w-4" />
        </div>
        <div className="text-muted-foreground leading-none">
          Data covers the last 6 months ({getLast6MonthsRange()})
        </div>
      </CardFooter>
    </Card>
  );
}

const formatRupiah = (angka: number) => `Rp ${angka.toLocaleString("id-ID")}`;

const renderLabel = ({
  value,
  cx,
  cy,
  midAngle,
  innerRadius,
  outerRadius,
}: any) => {
  const RADIAN = Math.PI / 180;
  const radius = innerRadius + (outerRadius - innerRadius) * 0.5;
  const x = cx + radius * Math.cos(-midAngle * RADIAN) * 2.5;
  const y = cy + radius * Math.sin(-midAngle * RADIAN) * 2.5;

  return (
    <text
      x={x}
      y={y}
      fill="#000"
      textAnchor={x > cx ? "start" : "end"}
      dominantBaseline="central"
      fontSize={12}
    >
      {formatRupiah(value)}
    </text>
  );
};

export function ChartPie({
  chartData,
  chartConfig,
}: {
  chartData: PieChartType[];
  chartConfig: ChartConfig;
}) {
  return (
    <Card className="flex h-full flex-col">
      <CardHeader className="items-center pb-0">
        <CardTitle>Top 7 Spending Categories</CardTitle>
        <CardDescription>
          Shows your highest spending categories over the last 3 months, helping
          you identify where most of your money goes.
        </CardDescription>
      </CardHeader>
      <CardContent className="flex-1 pb-0">
        <ChartContainer
          config={chartConfig}
          className="[&_.recharts-pie-label-text]:fill-foreground mx-auto max-h-[250px] pb-0"
        >
          <PieChart>
            <ChartTooltip content={<ChartTooltipContent hideLabel />} />
            <Pie
              data={chartData}
              dataKey="total"
              label={renderLabel}
              nameKey="parent_category_name"
            />
          </PieChart>
        </ChartContainer>
      </CardContent>
      <CardFooter className="flex-col gap-2 text-sm">
        <div className="flex items-center gap-2 leading-none font-medium">
          Most of your spending is concentrated in a few key categories.
        </div>
        <div className="text-muted-foreground leading-none">
          Your top expense categories is {chartData[0].parent_category_name} with a total of {formatRupiah(chartData[0].total)}
        </div>
      </CardFooter>
    </Card>
  );
}
