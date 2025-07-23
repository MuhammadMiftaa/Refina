import { ChartConfig } from "@/components/ui/chart";
import { AreaChartType, BarChartType, PieChartType } from "@/types/Chart";
import { ChartArea, ChartBar, ChartPie } from "./charts";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { IoIosTrendingDown, IoIosTrendingUp } from "react-icons/io";
import { PiBoxArrowDownLight } from "react-icons/pi";
import { TbMoneybag } from "react-icons/tb";
import { IoWalletOutline } from "react-icons/io5";
import { CiMoneyCheck1 } from "react-icons/ci";
import { getBackendURL } from "@/lib/readenv";
import Cookies from "js-cookie";
import { useQuery } from "@tanstack/react-query";
import { buildPieChartConfig } from "@/helper/Helper";
import { useEffect, useState } from "react";
import { UserSummaryType } from "@/types/UserSummary";
import { INVESTMENT_TYPE } from "@/helper/Data";
import { BsFillSafeFill } from "react-icons/bs";
import { NumberTicker } from "@/components/ui/number-ticker";

async function fetchUserSummary() {
  const backendURL = getBackendURL();

  const token = Cookies.get("token");

  const res = await fetch(`${backendURL}/transactions/user-summary/detail`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });
  if (!res.ok) {
    throw new Error("Failed to fetch transactions");
  }

  return res.json();
}

async function fetchUserMonthlySummary() {
  const backendURL = getBackendURL();

  const token = Cookies.get("token");

  const res = await fetch(
    `${backendURL}/transactions/user-monthly-summary/detail`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
    },
  );
  if (!res.ok) {
    throw new Error("Failed to fetch transactions");
  }

  return res.json();
}

async function fetchUserMostExpenses() {
  const backendURL = getBackendURL();

  const token = Cookies.get("token");

  const res = await fetch(
    `${backendURL}/transactions/user-most-expenses/detail`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
    },
  );
  if (!res.ok) {
    throw new Error("Failed to fetch transactions");
  }

  return res.json();
}

async function fetchUserWalletDailySummary() {
  const backendURL = getBackendURL();

  const token = Cookies.get("token");

  const res = await fetch(
    `${backendURL}/transactions/user-wallet-daily-summary/detail`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
    },
  );
  if (!res.ok) {
    throw new Error("Failed to fetch transactions");
  }

  return res.json();
}

const areaChartConfig = {
  balance: {
    label: "Balance",
  },
  bank: {
    label: "Bank",
    color: "#2B7FFF",
    fill: "#539CFF",
    stroke: "#1E6EEB",
    stopColor: "#155FCF",
  },
  "e-wallet": {
    label: "E-Wallet",
    color: "#8EC5FF",
    fill: "#A3D1FF",
    stroke: "#6FB5FF",
    stopColor: "#479EFF",
  },
  physical: {
    label: "Cash",
    color: "#66D9EF",
    fill: "#8BE9FD",
    stroke: "#42C9E3",
    stopColor: "#1BAFCB",
  },
  others: {
    label: "Others",
    color: "#A78BFA",
    fill: "#C4B5FD",
    stroke: "#8B5CF6",
    stopColor: "#7C3AED",
  },
} satisfies ChartConfig;

const barChartConfig = {
  total_income: {
    label: "Income",
    color: "#8EC5FF",
  },
  total_expense: {
    label: "Expense",
    color: "#2A7CFA",
  },
} satisfies ChartConfig;

const InvestmentsData = [
  {
    id: "873fd980-6663-466f-9d9d-16a5dddb7a03",
    investment_type: "Gold",
    user_id: "a79a39e5-d70f-41ae-b7e6-36246a99172d",
    name: "Emas Antam",
    amount: 9335000,
    quantity: 5,
    unit: "Gram",
    investment_date: "2024-06-09T15:04:05Z",
    description: "Emas Batang Antam",
  },
  {
    id: "d8da83b3-15ee-4b0e-8289-0924b61c1a0f",
    investment_type: "Gold",
    user_id: "a79a39e5-d70f-41ae-b7e6-36246a99172d",
    name: "Emas Pluang",
    amount: 3713738,
    quantity: 2,
    unit: "Gram",
    investment_date: "2024-04-13T15:04:05Z",
    description: "Emas Digital Pluang",
  },
  {
    id: "bf679652-b03e-4dc7-9dcf-1b9d43015f52",
    investment_type: "Others",
    user_id: "a79a39e5-d70f-41ae-b7e6-36246a99172d",
    name: "Bitcoin Binance",
    amount: 135336,
    quantity: 0.00006993,
    unit: "BTC",
    investment_date: "2025-07-23T15:04:05Z",
    description: "Bitcoin/USDT Binance",
  },
];

export default function Dashboard() {
  const { data: userSummary, isLoading: userSummaryLoading } = useQuery({
    queryKey: ["user-summary"],
    queryFn: fetchUserSummary,
  });
  const { data: userMonthlySummary, isLoading: userMonthlySummaryLoading } =
    useQuery({
      queryKey: ["user-monthly-summary"],
      queryFn: fetchUserMonthlySummary,
    });
  const { data: userMostExpenses, isLoading: userMostExpensesLoading } =
    useQuery({
      queryKey: ["user-most-expenses"],
      queryFn: fetchUserMostExpenses,
    });
  const {
    data: userWalletDailySummary,
    isLoading: userWalletDailySummaryLoading,
  } = useQuery({
    queryKey: ["user-wallet-daily-summary"],
    queryFn: fetchUserWalletDailySummary,
  });

  const UserSummary: UserSummaryType = userSummary?.data[0] || {};
  const UserMonthlySummary: BarChartType[] = userMonthlySummary?.data || [];
  const UserMostExpenses: PieChartType[] = userMostExpenses?.data || [];
  const UserWalletDailySummary: AreaChartType[] =
    userWalletDailySummary?.data || [];

  const [UserMostExpensesWithFill, setUserMostExpensesWithFill] = useState<
    PieChartType[]
  >([]);
  const [pieChartConfig, setPieChartConfig] = useState<ChartConfig>();

  useEffect(() => {
    if (UserMostExpenses.length > 0) {
      const pieChartConfig = buildPieChartConfig(UserMostExpenses);
      setPieChartConfig(pieChartConfig);
      const updatedData = UserMostExpenses.map((item) => ({
        ...item,
        fill: pieChartConfig[item.parent_category_name]?.color || "#000000", // Default color if not found
      }));
      setUserMostExpensesWithFill(updatedData);
    }
  }, [UserMostExpenses]);

  return (
    <div className="font-poppins min-h-screen w-screen p-4 md:w-full md:p-6">
      <div className="flex flex-col items-start justify-between gap-4 md:flex-row md:items-center">
        <h1 className="text-3xl font-semibold md:text-4xl">Dashboard</h1>
      </div>
      <div className="mt-4 grid grid-cols-1 gap-4 md:grid-cols-3">
        {!userSummaryLoading && UserSummary && (
          <div className="col-span-3 grid grid-cols-1 gap-4 md:grid-cols-4">
            {/* INCOME */}
            <Card>
              <CardHeader className="flex justify-between">
                <div>
                  <CardTitle className="text-sm md:text-lg">
                    Total Income
                  </CardTitle>
                  <CardDescription className="text-xs md:text-sm">
                    {new Intl.DateTimeFormat("en-US", {
                      month: "long",
                      year: "numeric",
                    }).format(new Date())}
                  </CardDescription>
                </div>
                <div className="rounded-full bg-sky-50 p-2 text-base text-sky-500 md:text-xl">
                  <PiBoxArrowDownLight />
                </div>
              </CardHeader>
              <CardContent className="flex items-center justify-between gap-4 py-0 text-lg font-bold md:text-2xl">
                <h1>
                  Rp <NumberTicker value={UserSummary?.income_now} />
                </h1>
                <div
                  className={`${UserSummary?.user_income_growth_percentage >= 0 ? "bg-emerald-50 text-emerald-600" : "bg-rose-50 text-rose-600"} flex items-center gap-2 rounded px-1 py-0.5 text-xs font-medium md:text-base`}
                >
                  {UserSummary?.user_income_growth_percentage >= 0 ? (
                    <IoIosTrendingUp />
                  ) : (
                    <IoIosTrendingDown />
                  )}
                  <h2>
                    <NumberTicker value={Math.abs(UserSummary?.user_income_growth_percentage)} decimalPlaces={2} /> %
                  </h2>
                </div>
              </CardContent>
            </Card>
            {/* EXPENSE */}
            <Card>
              <CardHeader className="flex justify-between">
                <div>
                  <CardTitle className="text-sm md:text-lg">
                    Total Expense
                  </CardTitle>
                  <CardDescription className="text-xs md:text-sm">
                    {new Intl.DateTimeFormat("en-US", {
                      month: "long",
                      year: "numeric",
                    }).format(new Date())}
                  </CardDescription>
                </div>
                <div className="rounded-full bg-sky-50 p-2 text-base text-sky-500 md:text-xl">
                  <TbMoneybag />
                </div>
              </CardHeader>
              <CardContent className="flex items-center justify-between gap-4 py-0 text-lg font-bold md:text-2xl">
                <h1>
                  Rp <NumberTicker value={UserSummary?.expense_now} />
                </h1>
                <div
                  className={`${UserSummary?.user_expense_growth_percentage <= 0 ? "bg-emerald-50 text-emerald-600" : "bg-rose-50 text-rose-600"} flex items-center gap-2 rounded px-1 py-0.5 text-xs font-medium md:text-base`}
                >
                  {UserSummary?.user_expense_growth_percentage >= 0 ? (
                    <IoIosTrendingUp />
                  ) : (
                    <IoIosTrendingDown />
                  )}
                  <h2>
                    <NumberTicker value={Math.abs(UserSummary?.user_expense_growth_percentage)} decimalPlaces={2} /> %
                  </h2>
                </div>
              </CardContent>
            </Card>
            {/* PROFIT */}
            <Card>
              <CardHeader className="flex justify-between">
                <div>
                  <CardTitle className="text-sm md:text-lg">
                    Total Profit
                  </CardTitle>
                  <CardDescription className="text-xs md:text-sm">
                    {new Intl.DateTimeFormat("en-US", {
                      month: "long",
                      year: "numeric",
                    }).format(new Date())}
                  </CardDescription>
                </div>
                <div className="rounded-full bg-sky-50 p-2 text-base text-sky-500 md:text-xl">
                  <CiMoneyCheck1 />
                </div>
              </CardHeader>
              <CardContent className="flex items-center justify-between gap-4 py-0 text-lg font-bold md:text-2xl">
                <h1>
                  Rp <NumberTicker value={UserSummary?.profit_now} />
                </h1>
                <div
                  className={`${UserSummary?.user_profit_growth_percentage >= 0 ? "bg-emerald-50 text-emerald-600" : "bg-rose-50 text-rose-600"} flex items-center gap-2 rounded px-1 py-0.5 text-xs font-medium md:text-base`}
                >
                  {UserSummary?.user_profit_growth_percentage >= 0 ? (
                    <IoIosTrendingUp />
                  ) : (
                    <IoIosTrendingDown />
                  )}
                  <h2>
                    <NumberTicker value={Math.abs(UserSummary?.user_profit_growth_percentage)} decimalPlaces={2} /> %
                  </h2>
                </div>
              </CardContent>
            </Card>
            {/* TOTAL BALANCE */}
            <Card className="bg-gradient-to-r from-sky-200 to-sky-300">
              <CardHeader className="flex justify-between">
                <div>
                  <CardTitle className="text-sm md:text-lg">
                    Total Balance
                  </CardTitle>
                  <CardDescription className="text-xs md:text-sm">
                    {new Intl.DateTimeFormat("en-US", {
                      month: "long",
                      year: "numeric",
                    }).format(new Date())}
                  </CardDescription>
                </div>
                <div className="rounded-full bg-sky-50 p-2 text-base text-sky-500 md:text-xl">
                  <IoWalletOutline />
                </div>
              </CardHeader>
              <CardContent className="flex items-center justify-between gap-4 py-0 text-lg font-bold md:text-2xl">
                <h1>
                  Rp <NumberTicker value={UserSummary?.balance_now} />
                </h1>
                <div className="flex items-center gap-2 rounded px-1 py-0.5 text-xs font-medium text-black md:text-base">
                  {UserSummary?.user_balance_growth_percentage >= 0 ? (
                    <IoIosTrendingUp />
                  ) : (
                    <IoIosTrendingDown />
                  )}
                  <h2>
                    <NumberTicker value={Math.abs(UserSummary?.user_balance_growth_percentage)} decimalPlaces={2} /> %
                  </h2>
                </div>
              </CardContent>
            </Card>
          </div>
        )}
        <div className="col-span-3 row-start-2">
          {!userWalletDailySummaryLoading &&
            UserWalletDailySummary.length > 0 && (
              <ChartArea
                chartData={UserWalletDailySummary}
                chartConfig={areaChartConfig}
              />
            )}
        </div>
        <div className="col-span-3 row-start-3 h-full w-full md:col-span-1">
          {!userMonthlySummaryLoading && UserMonthlySummary.length > 0 && (
            <ChartBar
              chartData={UserMonthlySummary}
              chartConfig={barChartConfig}
            />
          )}
        </div>
        <div className="col-span-3 row-start-4 h-full md:col-span-1 md:row-start-3">
          {!userMostExpensesLoading &&
            UserMostExpenses.length > 0 &&
            pieChartConfig &&
            UserMostExpensesWithFill && (
              <ChartPie
                chartData={UserMostExpensesWithFill}
                chartConfig={pieChartConfig}
              />
            )}
        </div>
        <Card className="col-span-3 row-start-5 flex h-full flex-col md:col-span-1 md:row-start-3">
          <CardHeader className="items-center pb-0">
            <CardTitle>Total Investments</CardTitle>
            <CardDescription>
              Overview of your investments, including types, quantities, and
              values.
            </CardDescription>
          </CardHeader>
          <CardContent className="">
            <table className="w-full">
              <thead className="text-xs md:text-base">
                <tr>
                  <th className="py-2 text-left font-medium">
                    Investment Type
                  </th>
                  <th className="py-2 text-center font-medium">Quantity</th>
                  <th className="py-2 text-right font-medium">Unit</th>
                </tr>
              </thead>
              <tbody className="text-xs md:text-base">
                {InvestmentsData.map((investment) => (
                  <tr key={investment.id}>
                    <td className="flex items-center gap-2 py-1 text-left">
                      <div
                        className={`h-2 w-2 rounded-full ${investment.investment_type in INVESTMENT_TYPE && INVESTMENT_TYPE[investment.investment_type as keyof typeof INVESTMENT_TYPE].color}`}
                      />
                      {investment.investment_type}
                    </td>
                    <td className="py-1 text-center">{investment.quantity}</td>
                    <td className="py-1 text-right">{investment.unit}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
          <CardFooter className="mt-auto flex-col gap-2">
            <div className="flex items-center gap-2 text-right font-medium">
              Total investments value is{" "}
              <span className="font-bold">
                Rp {UserSummary?.balance_now?.toLocaleString("id-ID")}
              </span>
              <BsFillSafeFill className="h-4 w-4" />
            </div>
            <div className="text-muted-foreground -mt-2 text-center">
              Your investments are diversified across various types, including
              gold, cryptocurrencies, and more.
            </div>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}
