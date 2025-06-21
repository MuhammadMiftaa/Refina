import { WalletType } from "@/types/UserWallet";
import { useQuery } from "@tanstack/react-query";
import Cookies from "js-cookie";
import { useState } from "react";
import { GiPayMoney, GiReceiveMoney } from "react-icons/gi";
import { DataTable } from "./data-table";
import { ArrowUpDown, MoreHorizontal } from "lucide-react";
import { ColumnDef } from "@tanstack/react-table";
import { Button } from "../../../components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "../../../components/ui/dropdown-menu";
import { BsArrowDownLeftCircle, BsArrowUpRightCircle } from "react-icons/bs";
import { formatRawDate, generateColorByType } from "@/helper/Helper";
import { PiArrowsLeftRightLight } from "react-icons/pi";
import { getBackendURL } from "@/lib/readenv";
import { TransactionType } from "@/types/UserTransaction";

async function fetchTransactions() {
  const backendURL = getBackendURL();

  const token = Cookies.get("token");

  const res = await fetch(`${backendURL}/users/transactions`, {
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

export default function Transactions() {
  const { data: transactionsData, isLoading: transactionsLoading } = useQuery({
    queryKey: ["transactions"],
    queryFn: fetchTransactions,
  });

  const Transactions: TransactionType[] = transactionsData?.data ?? [];

  if (transactionsLoading) {
    return (
      <div className="flex h-screen w-full items-center justify-center">
        <div className="loader"></div>
      </div>
    );
  }

  return (
    <div className="font-poppins min-h-screen w-screen p-4 md:w-full md:p-6">
      <div className="flex flex-col items-start justify-between gap-4 md:flex-row md:items-center">
        <h1 className="text-3xl font-semibold md:text-4xl">Your Transaction</h1>
      </div>

      <div className="mt-6 rounded-2xl">
        <DataTable columns={columns} data={Transactions} />
      </div>
    </div>
  );
}

const columns: ColumnDef<TransactionType>[] = [
  {
    accessorKey: "description",
    header: ({ column }) => (
      <button
        className="flex items-center"
        onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
      >
        Description
        <ArrowUpDown className="ml-2 h-4 w-4 cursor-pointer" />
      </button>
    ),
  },
  {
    accessorKey: "wallet_type",
    header: ({ column }) => (
      <button
        className="flex items-center"
        onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
      >
        Wallet Type
        <ArrowUpDown className="ml-2 h-4 w-4 cursor-pointer" />
      </button>
    ),
  },
  {
    accessorKey: "transaction_date",
    header: ({ column }) => (
      <button
        onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        className="mx-auto flex items-center justify-center text-center"
      >
        Date
        <ArrowUpDown className="ml-2 h-4 w-4 cursor-pointer" />
      </button>
    ),
    cell: ({ row }: { row: any }) => {
      const date: string = row.getValue("transaction_date");
      const formattedDate = formatRawDate(date);

      return (
        <div className="flex flex-col items-center">
          <h1 className="font-light">{formattedDate[1]}</h1>
          <p className="text-sm text-nowrap text-zinc-500">
            {formattedDate[0]}, {formattedDate[2]}
          </p>
        </div>
      );
    },
  },
  {
    accessorKey: "category_type",
    header: ({ column }) => (
      <button
        className="mx-auto flex items-center"
        onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
      >
        Category
        <ArrowUpDown className="ml-2 h-4 w-4 cursor-pointer" />
      </button>
    ),
    cell: ({ row }: { row: any }) => {
      const type: string = row.getValue("category_type");
      return (
        <h1
          className={`mx-auto flex w-fit items-center gap-2 rounded-2xl border px-3 py-1 font-light text-nowrap uppercase text-${generateColorByType(type)} border-${generateColorByType(type)}`}
        >
          {type.replace(/_/g, " ")}{" "}
          <span className={`text-${generateColorByType(type)}`}>
            {type === "expense" ? (
              <BsArrowUpRightCircle />
            ) : type === "income" ? (
              <BsArrowDownLeftCircle />
            ) : (
              <PiArrowsLeftRightLight />
            )}
          </span>
        </h1>
      );
    },
  },
  {
    accessorKey: "category_name",
    header: ({ column }) => (
      <button
        className="flex items-center"
        onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
      >
        Category
        <ArrowUpDown className="ml-2 h-4 w-4 cursor-pointer" />
      </button>
    ),
  },
  {
    accessorKey: "amount",
    header: ({ column }) => (
      <button
        className="flex w-full items-center justify-end"
        onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
      >
        Amount
        <ArrowUpDown className="ml-2 h-4 w-4 cursor-pointer" />
      </button>
    ),
    cell: ({ row }: { row: any }) => {
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
    cell: ({ row }: { row: any }) => {
      const transaction = row.original;
      const backendURL = getBackendURL();
      const [deleteConfirm, setDeleteConfirm] = useState(false);

      const deleteTransaction = async (id: string) => {
        const token = Cookies.get("token") || "";

        try {
          const res = await fetch(`${backendURL}/transactions/${id}`, {
            method: "DELETE",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
          });

          if (!res.ok) {
            throw new Error("Failed to add transaction");
          }

          setDeleteConfirm(false);
        } catch (error) {
          console.error("Error deleting transaction:", error);
          setDeleteConfirm(false);
        }
      };

      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">Open menu</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="font-poppins border">
            <DropdownMenuLabel>Actions</DropdownMenuLabel>
            <DropdownMenuItem
              className=""
              onClick={() => navigator.clipboard.writeText(transaction.id)}
            >
              Copy transaction ID
            </DropdownMenuItem>
            <DropdownMenuItem className="">Update Transaction</DropdownMenuItem>
            <DropdownMenuItem
              onClick={() => setDeleteConfirm(true)}
              className=""
            >
              Delete Transaction
            </DropdownMenuItem>
          </DropdownMenuContent>

          <div
            className={`fixed inset-0 flex items-center justify-center bg-zinc-800/70 duration-100 ${deleteConfirm ? "scale-100" : "scale-0"}`}
          >
            <div className="flex flex-col rounded-xl bg-sky-100 p-7 text-black">
              <h1 className="text-lg font-bold">
                Are you sure to delete this transaction?
              </h1>
              <hr />
              <p>This action cannot be canceled.</p>
              <div className="mt-4 flex justify-end gap-3">
                <button
                  onClick={() => setDeleteConfirm(false)}
                  className="cursor-pointer rounded-lg bg-zinc-300 px-4 py-1 text-lg duration-300 hover:bg-zinc-600 hover:text-white"
                >
                  Cancel
                </button>
                <button
                  onClick={() => deleteTransaction(transaction.id)}
                  className="cursor-pointer rounded-lg bg-red-500 px-4 py-1 text-lg duration-300 hover:bg-black hover:text-white"
                >
                  Delete
                </button>
              </div>
            </div>
          </div>
        </DropdownMenu>
      );
    },
  },
];
