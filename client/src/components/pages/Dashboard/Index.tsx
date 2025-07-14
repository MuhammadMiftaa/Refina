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

const areaChartData: AreaChartType[] = [
  { date: "2024-04-01", bank: 222, "e-wallet": 150, cash: 278, others: 181 },
  { date: "2024-04-02", bank: 97, "e-wallet": 180, cash: 321, others: 180 },
  { date: "2024-04-03", bank: 167, "e-wallet": 120, cash: 222, others: 140 },
  { date: "2024-04-04", bank: 242, "e-wallet": 260, cash: 332, others: 370 },
  { date: "2024-04-05", bank: 373, "e-wallet": 290, cash: 276, others: 287 },
  { date: "2024-04-06", bank: 301, "e-wallet": 340, cash: 421, others: 301 },
  { date: "2024-04-07", bank: 245, "e-wallet": 180, cash: 359, others: 189 },
  { date: "2024-04-08", bank: 409, "e-wallet": 320, cash: 284, others: 240 },
  { date: "2024-04-09", bank: 59, "e-wallet": 110, cash: 291, others: 338 },
  { date: "2024-04-10", bank: 261, "e-wallet": 190, cash: 272, others: 386 },
  { date: "2024-04-11", bank: 327, "e-wallet": 350, cash: 301, others: 130 },
  { date: "2024-04-12", bank: 292, "e-wallet": 210, cash: 488, others: 131 },
  { date: "2024-04-13", bank: 342, "e-wallet": 380, cash: 334, others: 368 },
  { date: "2024-04-14", bank: 137, "e-wallet": 220, cash: 277, others: 126 },
  { date: "2024-04-15", bank: 120, "e-wallet": 170, cash: 229, others: 142 },
  { date: "2024-04-16", bank: 138, "e-wallet": 190, cash: 414, others: 126 },
  { date: "2024-04-17", bank: 446, "e-wallet": 360, cash: 152, others: 281 },
  { date: "2024-04-18", bank: 364, "e-wallet": 410, cash: 333, others: 236 },
  { date: "2024-04-19", bank: 243, "e-wallet": 180, cash: 395, others: 254 },
  { date: "2024-04-20", bank: 89, "e-wallet": 150, cash: 479, others: 267 },
  { date: "2024-04-21", bank: 137, "e-wallet": 200, cash: 390, others: 232 },
  { date: "2024-04-22", bank: 224, "e-wallet": 170, cash: 175, others: 216 },
  { date: "2024-04-23", bank: 138, "e-wallet": 230, cash: 335, others: 170 },
  { date: "2024-04-24", bank: 387, "e-wallet": 290, cash: 229, others: 196 },
  { date: "2024-04-25", bank: 215, "e-wallet": 250, cash: 482, others: 375 },
  { date: "2024-04-26", bank: 75, "e-wallet": 130, cash: 468, others: 388 },
  { date: "2024-04-27", bank: 383, "e-wallet": 420, cash: 426, others: 172 },
  { date: "2024-04-28", bank: 122, "e-wallet": 180, cash: 240, others: 167 },
  { date: "2024-04-29", bank: 315, "e-wallet": 240, cash: 316, others: 188 },
  { date: "2024-04-30", bank: 454, "e-wallet": 380, cash: 459, others: 393 },
  { date: "2024-05-01", bank: 165, "e-wallet": 220, cash: 283, others: 223 },
  { date: "2024-05-02", bank: 293, "e-wallet": 310, cash: 390, others: 156 },
  { date: "2024-05-03", bank: 247, "e-wallet": 190, cash: 187, others: 104 },
  { date: "2024-05-04", bank: 385, "e-wallet": 420, cash: 451, others: 211 },
  { date: "2024-05-05", bank: 481, "e-wallet": 390, cash: 117, others: 300 },
  { date: "2024-05-06", bank: 498, "e-wallet": 520, cash: 455, others: 388 },
  { date: "2024-05-07", bank: 388, "e-wallet": 300, cash: 364, others: 296 },
  { date: "2024-05-08", bank: 149, "e-wallet": 210, cash: 275, others: 111 },
  { date: "2024-05-09", bank: 227, "e-wallet": 180, cash: 280, others: 176 },
  { date: "2024-05-10", bank: 293, "e-wallet": 330, cash: 356, others: 348 },
  { date: "2024-05-11", bank: 335, "e-wallet": 270, cash: 289, others: 388 },
  { date: "2024-05-12", bank: 197, "e-wallet": 240, cash: 412, others: 232 },
  { date: "2024-05-13", bank: 197, "e-wallet": 160, cash: 439, others: 221 },
  { date: "2024-05-14", bank: 448, "e-wallet": 490, cash: 345, others: 309 },
  { date: "2024-05-15", bank: 473, "e-wallet": 380, cash: 293, others: 243 },
  { date: "2024-05-16", bank: 338, "e-wallet": 400, cash: 466, others: 239 },
  { date: "2024-05-17", bank: 499, "e-wallet": 420, cash: 342, others: 247 },
  { date: "2024-05-18", bank: 315, "e-wallet": 350, cash: 193, others: 354 },
  { date: "2024-05-19", bank: 235, "e-wallet": 180, cash: 238, others: 265 },
  { date: "2024-05-20", bank: 177, "e-wallet": 230, cash: 162, others: 192 },
  { date: "2024-05-21", bank: 82, "e-wallet": 140, cash: 482, others: 260 },
  { date: "2024-05-22", bank: 81, "e-wallet": 120, cash: 347, others: 310 },
  { date: "2024-05-23", bank: 252, "e-wallet": 290, cash: 262, others: 343 },
  { date: "2024-05-24", bank: 294, "e-wallet": 220, cash: 102, others: 251 },
  { date: "2024-05-25", bank: 201, "e-wallet": 250, cash: 445, others: 281 },
  { date: "2024-05-26", bank: 213, "e-wallet": 170, cash: 306, others: 147 },
  { date: "2024-05-27", bank: 420, "e-wallet": 460, cash: 428, others: 384 },
  { date: "2024-05-28", bank: 233, "e-wallet": 190, cash: 157, others: 194 },
  { date: "2024-05-29", bank: 78, "e-wallet": 130, cash: 364, others: 120 },
  { date: "2024-05-30", bank: 340, "e-wallet": 280, cash: 365, others: 270 },
  { date: "2024-05-31", bank: 178, "e-wallet": 230, cash: 182, others: 161 },
  { date: "2024-06-01", bank: 178, "e-wallet": 200, cash: 487, others: 159 },
  { date: "2024-06-02", bank: 470, "e-wallet": 410, cash: 248, others: 319 },
  { date: "2024-06-03", bank: 103, "e-wallet": 160, cash: 175, others: 381 },
  { date: "2024-06-04", bank: 439, "e-wallet": 380, cash: 234, others: 212 },
  { date: "2024-06-05", bank: 88, "e-wallet": 140, cash: 238, others: 249 },
  { date: "2024-06-06", bank: 294, "e-wallet": 250, cash: 204, others: 375 },
  { date: "2024-06-07", bank: 323, "e-wallet": 370, cash: 341, others: 362 },
  { date: "2024-06-08", bank: 385, "e-wallet": 320, cash: 354, others: 295 },
  { date: "2024-06-09", bank: 438, "e-wallet": 480, cash: 431, others: 370 },
  { date: "2024-06-10", bank: 155, "e-wallet": 200, cash: 392, others: 249 },
  { date: "2024-06-11", bank: 92, "e-wallet": 150, cash: 184, others: 180 },
  { date: "2024-06-12", bank: 492, "e-wallet": 420, cash: 326, others: 137 },
  { date: "2024-06-13", bank: 81, "e-wallet": 130, cash: 284, others: 145 },
  { date: "2024-06-14", bank: 426, "e-wallet": 380, cash: 350, others: 106 },
  { date: "2024-06-15", bank: 307, "e-wallet": 350, cash: 310, others: 185 },
  { date: "2024-06-16", bank: 371, "e-wallet": 310, cash: 127, others: 303 },
  { date: "2024-06-17", bank: 475, "e-wallet": 520, cash: 221, others: 194 },
  { date: "2024-06-18", bank: 107, "e-wallet": 170, cash: 199, others: 229 },
  { date: "2024-06-19", bank: 341, "e-wallet": 290, cash: 305, others: 381 },
  { date: "2024-06-20", bank: 408, "e-wallet": 450, cash: 312, others: 251 },
  { date: "2024-06-21", bank: 169, "e-wallet": 210, cash: 466, others: 360 },
  { date: "2024-06-22", bank: 317, "e-wallet": 270, cash: 411, others: 284 },
  { date: "2024-06-23", bank: 480, "e-wallet": 530, cash: 194, others: 173 },
  { date: "2024-06-24", bank: 132, "e-wallet": 180, cash: 180, others: 132 },
  { date: "2024-06-25", bank: 141, "e-wallet": 190, cash: 444, others: 283 },
  { date: "2024-06-26", bank: 434, "e-wallet": 380, cash: 236, others: 165 },
  { date: "2024-06-27", bank: 448, "e-wallet": 490, cash: 404, others: 348 },
  { date: "2024-06-28", bank: 149, "e-wallet": 200, cash: 127, others: 326 },
  { date: "2024-06-29", bank: 103, "e-wallet": 160, cash: 450, others: 397 },
  { date: "2024-06-30", bank: 446, "e-wallet": 400, cash: 358, others: 282 },
];

const areaChartConfig = {
  visitors: {
    label: "Visitors",
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
  cash: {
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

const barChartData: BarChartType[] = [
  { month: "January", income: 40000, expense: 24000 },
  { month: "February", income: 30000, expense: 13980 },
  { month: "March", income: 20000, expense: 98000 },
  { month: "April", income: 27800, expense: 39080 },
  { month: "May", income: 18900, expense: 48000 },
  { month: "June", income: 23900, expense: 38000 },
];

const barChartConfig = {
  income: {
    label: "Income",
    color: "#8EC5FF",
  },
  expense: {
    label: "Expense",
    color: "#2A7CFA",
  },
} satisfies ChartConfig;

const pieChartData: PieChartType[] = [
  { category: "Food", value: 3000, fill: "#FF6384" },
  { category: "Transport", value: 2000, fill: "#36A2EB" },
  { category: "Entertainment", value: 1500, fill: "#FFCE56" },
  { category: "Utilities", value: 2500, fill: "#4BC0C0" },
];

const pieChartConfig = {
  value: {
    label: "Value",
  },
  food: {
    label: "Food",
    color: "#FF6384",
  },
  transport: {
    label: "Transport",
    color: "#36A2EB",
  },
  entertainment: {
    label: "Entertainment",
    color: "#FFCE56",
  },
  utilities: {
    label: "Utilities",
    color: "#4BC0C0",
  },
} satisfies ChartConfig;

export default function Dashboard() {
  return (
    <div className="font-poppins min-h-screen w-screen p-4 md:w-full md:p-6">
      <div className="flex flex-col items-start justify-between gap-4 md:flex-row md:items-center">
        <h1 className="text-3xl font-semibold md:text-4xl">Dashboard</h1>
      </div>
      <div className="mt-4 grid grid-cols-3 gap-4">
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
              <h1>RP 4.500.000</h1>
              <div className="flex items-center gap-2 rounded bg-emerald-50 px-1 py-0.5 text-base font-medium text-emerald-600">
                <IoIosTrendingUp />
                <h2>4.5 %</h2>
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
              <h1>RP 3.500.000</h1>
              <div className="flex items-center gap-2 rounded bg-rose-50 px-1 py-0.5 text-base font-medium text-rose-600">
                <IoIosTrendingDown />
                <h2>2.5 %</h2>
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
              <h1>RP 1.000.000</h1>
              <div className="flex items-center gap-2 rounded bg-rose-50 px-1 py-0.5 text-base font-medium text-rose-600">
                <IoIosTrendingDown />
                <h2>0.5 %</h2>
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
              <h1>RP 12.000.000</h1>
              <div className="flex items-center gap-2 rounded px-1 py-0.5 text-base font-medium text-black">
                <IoIosTrendingDown />
                <h2>5.5 %</h2>
              </div>
            </CardContent>
          </Card>
        </div>
        <div className="col-span-3">
          <ChartArea chartData={areaChartData} chartConfig={areaChartConfig} />
        </div>
        <div className="h-full">
          <ChartBar chartData={barChartData} chartConfig={barChartConfig} />
        </div>
        <div className="h-full">
          <ChartPie chartData={pieChartData} chartConfig={pieChartConfig} />
        </div>
        <div className="h-full">
          <Card className="h-full">
            <CardHeader>
              <CardTitle className="text-lg">Total Invesment</CardTitle>
              {/* <CardDescription>June 2024</CardDescription> */}
            </CardHeader>
            <CardContent className="">
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
