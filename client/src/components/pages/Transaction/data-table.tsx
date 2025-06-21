"use client";

import {
  ColumnDef,
  ColumnFiltersState,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  SortingState,
  useReactTable,
  VisibilityState,
} from "@tanstack/react-table";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../../ui/table";
import { useState } from "react";
import { Input } from "@/components/ui/input";
import { DropdownMenu } from "@radix-ui/react-dropdown-menu";
import {
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { IoSearchOutline } from "react-icons/io5";
import { HiMiniPlus } from "react-icons/hi2";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Link, useNavigate } from "react-router";

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
}

export function DataTable<TData, TValue>({
  columns,
  data,
}: DataTableProps<TData, TValue>) {
  const navigate = useNavigate();

  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({});
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    getFilteredRowModel: getFilteredRowModel(),
    getSortedRowModel: getSortedRowModel(),
    onColumnVisibilityChange: setColumnVisibility,
    state: {
      sorting,
      columnFilters,
      columnVisibility,
    },
  });

  return (
    <div>
      <div className="flex items-center py-4">
        <div className="flex w-full flex-col items-start justify-between gap-2 md:flex-row md:items-center md:gap-0">
          <div className="flex h-full items-center justify-center gap-12">
            <div className="rounded-[16px] bg-gradient-to-b from-stone-300/40 to-transparent p-[4px]">
              <button className="group rounded-[12px] bg-gradient-to-b from-white to-stone-200/40 p-[2px] shadow-[0_1px_3px_rgba(0,0,0,0.5)] active:scale-[0.995] active:shadow-[0_0px_1px_rgba(0,0,0,0.5)]">
                <div className="rounded-[8px] bg-gradient-to-b from-stone-200/40 to-white/80 px-2">
                  <div className="flex items-center justify-between gap-2 overflow-hidden">
                    <Input
                      placeholder="Search by description..."
                      value={
                        (table
                          .getColumn("description")
                          ?.getFilterValue() as string) ?? ""
                      }
                      onChange={(event) =>
                        table
                          .getColumn("description")
                          ?.setFilterValue(event.target.value)
                      }
                      className="max-w-sm border-none bg-transparent p-0 text-sm shadow-none focus:bg-transparent focus:shadow-none focus:ring-0 focus:outline-none"
                    />
                    <IoSearchOutline />
                  </div>
                </div>
              </button>
            </div>
          </div>
          <div className="flex flex-row-reverse items-center gap-2">
            <Dialog>
              <DialogTrigger>
                <div className="flex h-full items-center justify-center gap-12">
                  <div className="rounded-[16px] bg-gradient-to-b from-stone-300/40 to-transparent p-[4px]">
                    <button className="group rounded-[12px] bg-gradient-to-b from-white to-stone-200/40 p-[2px] shadow-[0_1px_3px_rgba(0,0,0,0.5)] active:scale-[0.995] active:shadow-[0_0px_1px_rgba(0,0,0,0.5)]">
                      <div className="rounded-[8px] bg-gradient-to-b from-stone-200/40 to-white/80 px-2 py-2">
                        <div className="flex cursor-pointer items-center gap-2">
                          <span className="text-sm font-semibold">New</span>
                          <HiMiniPlus />
                        </div>
                      </div>
                    </button>
                  </div>
                </div>
              </DialogTrigger>
              <DialogContent className="w-120 border-none bg-[rgba(255,255,255,0.5)] backdrop-blur-sm">
                <div className="flex gap-1 p-0 md:p-2">
                  <div>
                    <span className="box center inline-block h-3 w-3 rounded-full bg-rose-500"></span>
                  </div>
                  <div>
                    <span className="center inline-block h-3 w-3 rounded-full bg-blue-500"></span>
                  </div>
                  <div>
                    <span className="center inline-block h-3 w-3 rounded-full bg-indigo-500"></span>
                  </div>
                </div>
                <DialogHeader>
                  <DialogTitle>Select Transaction Type</DialogTitle>
                </DialogHeader>
                <div className="group/cards flex flex-col gap-[15px] md:flex-row">
                  <Link
                    to="/transactions/add/expense"
                    className="flex h-[100px] w-full transform cursor-pointer flex-col items-center justify-center rounded-[10px] bg-rose-500 text-center text-white transition duration-400 group-hover/cards:scale-90 group-hover/cards:blur-sm hover:scale-110 hover:blur-none"
                  >
                    <p className="text-base font-bold">Income</p>
                    <p className="text-[0.7em] italic">Pemasukan</p>
                  </Link>
                  <Link
                    to="/transactions/add/income"
                    className="flex h-[100px] w-full transform cursor-pointer flex-col items-center justify-center rounded-[10px] bg-blue-500 text-center text-white transition duration-400 group-hover/cards:scale-90 group-hover/cards:blur-sm hover:scale-110 hover:blur-none"
                  >
                    <p className="text-base font-bold">Expense</p>
                    <p className="text-[0.7em] italic">Pengeluaran</p>
                  </Link>
                  <Link
                    to="/transactions/add/fund_transfer"
                    className="flex h-[100px] w-full transform cursor-pointer flex-col items-center justify-center rounded-[10px] bg-indigo-500 text-center text-white transition duration-400 group-hover/cards:scale-90 group-hover/cards:blur-sm hover:scale-110 hover:blur-none"
                  >
                    <p className="text-base font-bold">Fund Transfer</p>
                    <p className="text-[0.7em] italic">Pindah Dana</p>
                  </Link>
                </div>
              </DialogContent>
            </Dialog>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <div className="flex h-full items-center justify-center gap-12">
                  <div className="rounded-[16px] bg-gradient-to-b from-stone-300/40 to-transparent p-[4px]">
                    <button className="group rounded-[12px] bg-gradient-to-b from-white to-stone-200/40 p-[2px] shadow-[0_1px_3px_rgba(0,0,0,0.5)] active:scale-[0.995] active:shadow-[0_0px_1px_rgba(0,0,0,0.5)]">
                      <div className="rounded-[8px] bg-gradient-to-b from-stone-200/40 to-white/80 px-2 py-2">
                        <div className="flex cursor-pointer items-center gap-2">
                          <span className="text-sm font-semibold">
                            Show Columns
                          </span>
                        </div>
                      </div>
                    </button>
                  </div>
                </div>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                {table
                  .getAllColumns()
                  .filter((column) => column.getCanHide())
                  .map((column) => {
                    return (
                      <DropdownMenuCheckboxItem
                        key={column.id}
                        className="capitalize"
                        checked={column.getIsVisible()}
                        onCheckedChange={(value) =>
                          column.toggleVisibility(!!value)
                        }
                      >
                        {column.id}
                      </DropdownMenuCheckboxItem>
                    );
                  })}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </div>
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead key={header.id}>
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext(),
                          )}
                    </TableHead>
                  );
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && "selected"}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext(),
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  No results.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      <div className="flex items-center justify-between md:justify-end space-x-2 py-4">
        <button
          onClick={() => table.previousPage()}
          disabled={!table.getCanPreviousPage()}
          className="group cursor-pointer rounded-xl border-[1px] border-slate-500 bg-gradient-to-b from-indigo-500 to-indigo-600 px-6 py-3 font-medium text-white shadow-[0px_4px_22px_0_rgba(99,102,241,.70)]"
        >
          <div className="relative overflow-hidden">
            <p className="duration-[1.125s] ease-[cubic-bezier(0.19,1,0.22,1)] group-hover:-translate-y-7">
              Previous
            </p>
            <p className="absolute top-7 left-0 duration-[1.125s] ease-[cubic-bezier(0.19,1,0.22,1)] group-hover:top-0">
              Previous
            </p>
          </div>
        </button>
        <button
          onClick={() => table.nextPage()}
          disabled={!table.getCanNextPage()}
          className="group cursor-pointer rounded-xl border-[1px] border-slate-500 bg-gradient-to-b from-indigo-500 to-indigo-600 px-6 py-3 font-medium text-white shadow-[0px_4px_22px_0_rgba(99,102,241,.70)]"
        >
          <div className="relative overflow-hidden">
            <p className="duration-[1.125s] ease-[cubic-bezier(0.19,1,0.22,1)] group-hover:-translate-y-7">
              Next
            </p>
            <p className="absolute top-7 left-0 duration-[1.125s] ease-[cubic-bezier(0.19,1,0.22,1)] group-hover:top-0">
              Next
            </p>
          </div>
        </button>
      </div>
    </div>
  );
}
