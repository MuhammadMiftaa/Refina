import { ChartConfig } from "@/components/ui/chart";
import { AreaChartType, BarChartType, PieChartType } from "@/types/Chart";
import { ChartArea, ChartBar, ChartPie } from "./charts";
import {
  Card,
  CardContent,
  CardDescription,
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

// const pieChartData: PieChartType[] = [
//   { category: "Food", value: 3000, fill: "#FF6384" },
//   { category: "Transport", value: 2000, fill: "#36A2EB" },
//   { category: "Entertainment", value: 1500, fill: "#FFCE56" },
//   { category: "Utilities", value: 2500, fill: "#4BC0C0" },
// ];

// const pieChartConfig = {
//   value: {
//     label: "Value",
//   },
//   food: {
//     label: "Food",
//     color: "#FF6384",
//   },
//   transport: {
//     label: "Transport",
//     color: "#36A2EB",
//   },
//   entertainment: {
//     label: "Entertainment",
//     color: "#FFCE56",
//   },
//   utilities: {
//     label: "Utilities",
//     color: "#4BC0C0",
//   },
// } satisfies ChartConfig;

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
      <div className="mt-4 grid grid-cols-3 gap-4">
        {!userSummaryLoading && UserSummary && (
          <div className="col-span-3 grid grid-cols-4 gap-4">
            {/* INCOME */}
            <Card>
              <CardHeader className="flex justify-between">
                <div>
                  <CardTitle className="text-lg">Total Income</CardTitle>
                  <CardDescription>June 2024</CardDescription>
                </div>
                <div className="rounded-full bg-sky-50 p-2 text-xl text-sky-500">
                  <PiBoxArrowDownLight />
                </div>
              </CardHeader>
              <CardContent className="flex items-center gap-4 text-2xl font-bold">
                <h1>Rp {UserSummary?.income_now?.toLocaleString("id-ID")}</h1>
                <div className="flex items-center gap-2 rounded bg-emerald-50 px-1 py-0.5 text-base font-medium text-emerald-600">
                  <IoIosTrendingUp />
                  <h2>{UserSummary?.user_income_growth_percentage} %</h2>
                </div>
              </CardContent>
            </Card>
            {/* EXPENSE */}
            <Card>
              <CardHeader className="flex justify-between">
                <div>
                  <CardTitle className="text-lg">Total Expense</CardTitle>
                  <CardDescription>June 2024</CardDescription>
                </div>
                <div className="rounded-full bg-sky-50 p-2 text-xl text-sky-500">
                  <TbMoneybag />
                </div>
              </CardHeader>
              <CardContent className="flex items-center gap-4 text-2xl font-bold">
                <h1>
                  Rp {UserSummary?.expense_now?.toLocaleString("id-ID")}
                </h1>
                <div className="flex items-center gap-2 rounded bg-rose-50 px-1 py-0.5 text-base font-medium text-rose-600">
                  <IoIosTrendingDown />
                  <h2>{UserSummary?.user_expense_growth_percentage} %</h2>
                </div>
              </CardContent>
            </Card>
            {/* PROFIT */}
            <Card>
              <CardHeader className="flex justify-between">
                <div>
                  <CardTitle className="text-lg">Total Profit</CardTitle>
                  <CardDescription>June 2024</CardDescription>
                </div>
                <div className="rounded-full bg-sky-50 p-2 text-xl text-sky-500">
                  <CiMoneyCheck1 />
                </div>
              </CardHeader>
              <CardContent className="flex items-center gap-4 text-2xl font-bold">
                <h1>Rp {UserSummary?.profit_now?.toLocaleString("id-ID")}</h1>
                <div className="flex items-center gap-2 rounded bg-rose-50 px-1 py-0.5 text-base font-medium text-rose-600">
                  <IoIosTrendingDown />
                  <h2>{UserSummary?.user_profit_growth_percentage} %</h2>
                </div>
              </CardContent>
            </Card>
            {/* TOTAL BALANCE */}
            <Card className="bg-gradient-to-r from-sky-200 to-sky-300">
              <CardHeader className="flex justify-between">
                <div>
                  <CardTitle className="text-lg">Total Balance</CardTitle>
                  <CardDescription>June 2024</CardDescription>
                </div>
                <div className="rounded-full bg-sky-50 p-2 text-xl text-sky-500">
                  <IoWalletOutline />
                </div>
              </CardHeader>
              <CardContent className="flex items-center gap-4 text-2xl font-bold">
                <h1>
                  Rp {UserSummary?.balance_now?.toLocaleString("id-ID")}
                </h1>
                <div className="flex items-center gap-2 rounded px-1 py-0.5 text-base font-medium text-black">
                  <IoIosTrendingDown />
                  <h2>{UserSummary?.user_balance_growth_percentage} %</h2>
                </div>
              </CardContent>
            </Card>
          </div>
        )}
        <div className="col-span-3">
          {!userWalletDailySummaryLoading &&
            UserWalletDailySummary.length > 0 && (
              <ChartArea
                chartData={UserWalletDailySummary}
                chartConfig={areaChartConfig}
              />
            )}
        </div>
        <div className="h-full">
          {!userMonthlySummaryLoading && UserMonthlySummary.length > 0 && (
            <ChartBar
              chartData={UserMonthlySummary}
              chartConfig={barChartConfig}
            />
          )}
        </div>
        <div className="h-full">
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
        <div className="h-full">
          <Card className="h-full">
            <CardHeader>
              <CardTitle className="text-lg">Total Invesment</CardTitle>
              {/* <CardDescription>June 2024</CardDescription> */}
            </CardHeader>
            <CardContent className=""></CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
