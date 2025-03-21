import { LucideHistory } from "lucide-react";
import { Link } from "react-router";
import { FiPlus } from "react-icons/fi";
import { DatePickerWithRange } from "../ui/date-picker";
import { FaFilePdf } from "react-icons/fa";
import { FaArrowTrendUp } from "react-icons/fa6";
import { DataTable } from "./data-table";
import { useState } from "react";
import { DateRange } from "react-day-picker";
import { addDays } from "date-fns";
import { columns } from "./columns";
import { transactionsType } from "@/types/Transactions";

export default function Analytics() {
  const transactions: transactionsType[] = [
    {
      id: "1",
      amount: 100000,
      transaction_type: "expense",
      date: "2024-11-20",
      description: "Belanja di supermarket",
      category: "Groceries",
      user_id: "101",
    },
    {
      id: "2",
      amount: 50000,
      transaction_type: "expense",
      date: "2024-11-21",
      description: "Transportasi online",
      category: "Transport",
      user_id: "101",
    },
    {
      id: "3",
      amount: 200000,
      transaction_type: "income",
      date: "2024-11-21",
      description: "Gaji mingguan",
      category: "Salary",
      user_id: "101",
    },
    {
      id: "4",
      amount: 80000,
      transaction_type: "expense",
      date: "2024-11-22",
      description: "Makan malam",
      category: "Food & Drinks",
      user_id: "102",
    },
    {
      id: "5",
      amount: 150000,
      transaction_type: "income",
      date: "2024-11-22",
      description: "Bonus proyek",
      category: "Bonus",
      user_id: "101",
    },
    {
      id: "6",
      amount: 75000,
      transaction_type: "expense",
      date: "2024-11-23",
      description: "Beli buku",
      category: "Education",
      user_id: "102",
    },
    {
      id: "7",
      amount: 120000,
      transaction_type: "expense",
      date: "2024-11-24",
      description: "Tagihan listrik",
      category: "Utilities",
      user_id: "101",
    },
    {
      id: "8",
      amount: 200000,
      transaction_type: "income",
      date: "2024-11-24",
      description: "Penjualan barang bekas",
      category: "Other Income",
      user_id: "102",
    },
    {
      id: "9",
      amount: 50000,
      transaction_type: "expense",
      date: "2024-11-25",
      description: "Ngopi di cafe",
      category: "Food & Drinks",
      user_id: "101",
    },
    {
      id: "10",
      amount: 500000,
      transaction_type: "income",
      date: "2024-11-26",
      description: "Transfer keluarga",
      category: "Family",
      user_id: "103",
    },
    {
      id: "11",
      amount: 300000,
      transaction_type: "expense",
      date: "2024-11-26",
      description: "Servis motor",
      category: "Transport",
      user_id: "102",
    },
    {
      id: "12",
      amount: 100000,
      transaction_type: "income",
      date: "2024-11-27",
      description: "Pemasukan tambahan",
      category: "Other Income",
      user_id: "101",
    },
    {
      id: "13",
      amount: 120000,
      transaction_type: "expense",
      date: "2024-11-28",
      description: "Pulsa internet",
      category: "Utilities",
      user_id: "102",
    },
    {
      id: "14",
      amount: 45000,
      transaction_type: "expense",
      date: "2024-11-28",
      description: "Cemilan malam",
      category: "Food & Drinks",
      user_id: "103",
    },
    {
      id: "15",
      amount: 600000,
      transaction_type: "income",
      date: "2024-11-29",
      description: "Pendapatan investasi",
      category: "Investment",
      user_id: "101",
    },
    {
      id: "16",
      amount: 100000,
      transaction_type: "expense",
      date: "2024-11-29",
      description: "Langganan streaming",
      category: "Entertainment",
      user_id: "102",
    },
    {
      id: "17",
      amount: 85000,
      transaction_type: "expense",
      date: "2024-11-30",
      description: "Makan siang",
      category: "Food & Drinks",
      user_id: "103",
    },
    {
      id: "18",
      amount: 250000,
      transaction_type: "income",
      date: "2024-11-30",
      description: "Hadiah teman",
      category: "Gift",
      user_id: "101",
    },
    {
      id: "19",
      amount: 90000,
      transaction_type: "expense",
      date: "2024-12-01",
      description: "Beli pakaian",
      category: "Shopping",
      user_id: "102",
    },
    {
      id: "20",
      amount: 150000,
      transaction_type: "income",
      date: "2024-12-01",
      description: "Jual buku bekas",
      category: "Other Income",
      user_id: "103",
    },
  ];

  const [date, setDate] = useState<DateRange | undefined>({
    from: new Date(2022, 0, 20),
    to: addDays(new Date(2022, 0, 20), 20),
  });

  return (
    <div className="bg-secondary min-h-screen w-full p-5 text-white">
      <div className="rounded-2xl bg-gradient-to-b from-purple-600 to-[#00D47E] p-0.5 text-purple-500">
        <div className="bg-primary flex w-full items-center justify-between rounded-2xl px-4 py-3">
          <div className="flex items-center gap-2">
            <LucideHistory />
            <h1 className="font-poppins text-2xl font-bold">
              Transactions History
            </h1>
          </div>
          <Link
            to={"add"}
            className="text-primary font-poppins flex items-center rounded-lg bg-[#00D47E] px-4 py-2 font-light"
          >
            <FiPlus />
            <span className="ml-2">Add</span>
          </Link>
        </div>
      </div>
      <div className="mt-6 flex gap-5">
        <div className="font-poppins bg-primary w-fit basis-2/5 rounded-2xl p-0.5">
          <h1 className="pt-3 pl-8 text-xs text-zinc-400">Select Date Range</h1>
          <DatePickerWithRange
            date={date}
            setDate={setDate}
          ></DatePickerWithRange>
        </div>
        <div className="font-poppins w-72 basis-2/5 rounded-2xl bg-[#00D47E] px-8 py-3">
          <h1 className="text-xs text-black">Current Balance</h1>
          <div className="flex items-center justify-between gap-3">
            <h2 className="text-2xl font-bold text-black">Rp. 8.000.000</h2>
            <div className="flex flex-col-reverse items-center text-xl text-black">
              <p className="-mt-1 text-xs">10%</p>
              <FaArrowTrendUp />
            </div>
          </div>
        </div>
        <div className="bg-primary font-poppins flex basis-1/5 cursor-pointer flex-col items-center justify-center rounded-2xl px-8 py-3 duration-100 hover:text-purple-500">
          <div className="w-fit text-3xl">
            <FaFilePdf />
          </div>
          <h1>Print to PDF</h1>
        </div>
      </div>
      <div className="bg-primary mt-6 rounded-2xl p-4">
        <DataTable columns={columns} data={transactions} />
      </div>
    </div>
  );
}
