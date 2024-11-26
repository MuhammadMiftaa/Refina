"use client";
import { MoreHorizontal } from "lucide-react";
import { ColumnDef } from "@tanstack/react-table";
import { Button } from "../../components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "../../components/ui/dropdown-menu";
import { transactionsType } from "../../types/Transactions";
import { BsArrowDownLeftCircle, BsArrowUpRightCircle } from "react-icons/bs";

export const transactions: transactionsType[] = [
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

export const columns: ColumnDef<transactionsType>[] = [
  {
    accessorKey: "description",
    header: "Description",
  },
  {
    accessorKey: "date",
    header: "Date",
  },
  {
    accessorKey: "transaction_type",
    header: "Type",
    cell: ({ row }) => {
      const type: string = row.getValue("transaction_type");
      return (
        <h1
          className={`flex items-center gap-2 font-light uppercase px-3 py-1 w-fit rounded-2xl border ${
            type === "expense" ? "border-red-500 " : "border-green-500 "
          }`}
        >
          {type}{" "}
          <span
            className={` ${
              type === "expense" ? "text-red-500 " : "text-green-500 "
            }`}
          >
            {type === "expense" ? (
              <BsArrowUpRightCircle />
            ) : (
              <BsArrowDownLeftCircle />
            )}
          </span>
        </h1>
      );
    },
  },
  {
    accessorKey: "category",
    header: "Category",
  },
  {
    accessorKey: "amount",
    header: () => <div className="text-right">Amount</div>,
    cell: ({ row }) => {
      const amount = parseFloat(row.getValue("amount"));
      const formatted = new Intl.NumberFormat("id-ID", {
        style: "currency",
        currency: "IDR",
        minimumFractionDigits: 0,
      }).format(amount);

      return <div className="text-right font-medium">{formatted}</div>;
    },
  },
  {
    id: "actions",
    cell: ({ row }) => {
      const transaction = row.original;

      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="ghost"
              className="h-8 w-8 p-0 hover:bg-purple-500 hover:text-primary"
            >
              <span className="sr-only">Open menu</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            align="end"
            className="bg-primary border border-purple-600 text-purple-300 font-poppins"
          >
            <DropdownMenuLabel>Actions</DropdownMenuLabel>
            <DropdownMenuItem
              className="hover:bg-slate-900"
              onClick={() => navigator.clipboard.writeText(transaction.id)}
            >
              Copy transaction ID
            </DropdownMenuItem>
            <DropdownMenuItem className="hover:bg-slate-900">
              View transactions attachments
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      );
    },
  },
];
