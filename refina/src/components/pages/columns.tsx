"use client";
import { MoreHorizontal } from "lucide-react";
import { ColumnDef } from "@tanstack/react-table";
import { Button } from "../ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import { UserTransactionsType } from "../../types/Transactions";
import { BsArrowDownLeftCircle, BsArrowUpRightCircle } from "react-icons/bs";

export const columns: ColumnDef<UserTransactionsType>[] = [
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
    cell: ({ row }: { row: any }) => {
      const type: string = row.getValue("transaction_type");
      return (
        <h1
          className={`flex w-fit items-center gap-2 rounded-2xl border px-3 py-1 font-light uppercase ${
            type === "expense" ? "border-red-500" : "border-green-500"
          }`}
        >
          {type}{" "}
          <span
            className={` ${
              type === "expense" ? "text-red-500" : "text-green-500"
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

      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="ghost"
              className="hover:text-primary h-8 w-8 p-0 hover:bg-purple-500"
            >
              <span className="sr-only">Open menu</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            align="end"
            className="bg-primary font-poppins border border-purple-600 text-purple-300"
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
