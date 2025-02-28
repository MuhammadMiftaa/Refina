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

export default function Home() {
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
    <div className="bg-secondary w-full text-white p-5 min-h-screen">
      <div className="bg-gradient-to-b from-purple-600 to-[#00D47E] text-purple-500 rounded-2xl p-0.5">
        <div className="bg-primary w-full rounded-2xl flex justify-between items-center px-4 py-3">
          <div className="flex items-center gap-2">
            <LucideHistory />
            <h1 className="text-2xl font-bold font-poppins">
              Transactions History
            </h1>
          </div>
          <Link
            to={"add"}
            className="bg-[#00D47E] text-primary py-2 px-4 rounded-lg flex items-center font-poppins font-light"
          >
            <FiPlus />
            <span className="ml-2">Add</span>
          </Link>
        </div>
      </div>
      <div className="flex gap-5 mt-6">
        <div className="p-0.5 w-fit rounded-2xl font-poppins bg-primary basis-2/5">
          <h1 className="text-xs text-zinc-400 pt-3 pl-8">Select Date Range</h1>
          <DatePickerWithRange
            date={date}
            setDate={setDate}
          ></DatePickerWithRange>
        </div>
        <div className="bg-[#00D47E] rounded-2xl font-poppins py-3 px-8 w-72 basis-2/5">
          <h1 className="text-xs text-black">Current Balance</h1>
          <div className="flex items-center gap-3 justify-between">
            <h2 className=" text-2xl text-black font-bold">Rp. 8.000.000</h2>
            <div className="text-xl text-black flex flex-col-reverse items-center">
              <p className="text-xs -mt-1">10%</p>
              <FaArrowTrendUp />
            </div>
          </div>
        </div>
        <div className="rounded-2xl bg-primary font-poppins py-3 px-8 flex flex-col items-center justify-center basis-1/5 cursor-pointer duration-100 hover:text-purple-500">
          <div className="text-3xl w-fit">
            <FaFilePdf />
          </div>
          <h1>Print to PDF</h1>
        </div>
      </div>
      <div className="bg-primary rounded-2xl p-4 mt-6">
        <DataTable columns={columns} data={transactions} />
      </div>
    </div>
  );
}
