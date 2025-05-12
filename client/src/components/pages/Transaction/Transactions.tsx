import { WalletType } from "@/types/UserWallet";
import { useQuery } from "@tanstack/react-query";
import Cookies from "js-cookie";
import { useEffect, useState } from "react";
import { FaMoneyBillTransfer } from "react-icons/fa6";
import { GiPayMoney, GiReceiveMoney } from "react-icons/gi";
import { useNavigate } from "react-router";
import styled from "styled-components";
import { DataTable } from "./data-table";
import { UserTransactionsType } from "@/types/Transactions";
import { MoreHorizontal } from "lucide-react";
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

async function fetchWallets() {
  const backendURL = getBackendURL();

  const token = Cookies.get("token");

  const res = await fetch(`${backendURL}/users/wallets`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });

  if (!res.ok) {
    throw new Error("Failed to fetch wallets");
  }

  return res.json();
}

async function fetchTransactions() {
  const backendURL = getBackendURL();

  const token = Cookies.get("token");

  const res = await fetch(`${backendURL}/transactions/user`, {
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
  const navigate = useNavigate();

  const { data: walletData, isLoading: walletLoading } = useQuery({
    queryKey: ["wallets"],
    queryFn: fetchWallets,
  });
  const { data: transactionsData, isLoading: transactionsLoading } = useQuery({
    queryKey: ["transactions"],
    queryFn: fetchTransactions,
  });

  const Wallets: WalletType = walletData?.data ?? {
    user_id: "",
    name: "",
    email: "",
    wallets: [],
  };

  const Transactions: UserTransactionsType[] = transactionsData?.data ?? [];

  const [wallets, setWallets] = useState(Wallets.wallets);

  useEffect(() => {
    setWallets(Wallets.wallets);
  }, [Wallets]);

  if (walletLoading || transactionsLoading) {
    return (
      <div className="flex h-screen w-full items-center justify-center">
        <div className="loader"></div>
      </div>
    );
  }

  return (
    <div className="font-poppins min-h-screen w-full md:px-4">
      <div className="flex flex-col items-start justify-between gap-4 md:flex-row md:items-center">
        <h1 className="text-3xl font-semibold md:text-4xl">Your Transaction</h1>
        {wallets.length > 0 && (
          <div className="flex flex-wrap items-center gap-5">
            <FundTransfer
              onclick={() => navigate("/transactions/add/fund_transfer")}
            />
            <ExpenseButton
              onclick={() => navigate("/transactions/add/expense")}
            />
            <IncomeButton
              onclick={() => navigate("/transactions/add/income")}
            />
          </div>
        )}
      </div>

      <div className="mt-6 rounded-2xl p-4">
        <DataTable columns={columns} data={Transactions} />
      </div>
    </div>
  );
}

const ExpenseButton = ({ onclick }: { onclick: () => void }) => {
  return (
    <StyledWrapper onClick={onclick}>
      <button className="Btn w-32 bg-[radial-gradient(100%_100%_at_100%_0%,_#FF7F7F_0%,_#D50000_100%)]">
        Expense
        <GiPayMoney className="icon" />
      </button>
    </StyledWrapper>
  );
};

const IncomeButton = ({ onclick }: { onclick: () => void }) => {
  return (
    <StyledWrapper onClick={onclick}>
      <button className="Btn w-32 bg-[radial-gradient(100%_100%_at_100%_0%,_#A8FF78_0%,_#00A86B_100%)]">
        Income
        <GiReceiveMoney className="icon" />
      </button>
    </StyledWrapper>
  );
};

const FundTransfer = ({ onclick }: { onclick: () => void }) => {
  return (
    <StyledWrapper onClick={onclick}>
      <button className="Btn w-44 bg-[radial-gradient(100%_100%_at_100%_0%,_#FFE177_0%,_#FFA500_100%)]">
        Fund Transfer
        <FaMoneyBillTransfer className="icon" />
      </button>
    </StyledWrapper>
  );
};

const StyledWrapper = styled.div`
  .Btn {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: flex-start;
    height: 40px;
    border: none;
    padding: 0px 20px;
    color: black;
    font-weight: 500;
    cursor: pointer;
    border-radius: 10px;
    box-shadow: 5px 5px 0px #000;
    transition-duration: 0.3s;
  }

  .icon {
    width: 13px;
    position: absolute;
    right: 0;
    margin-right: 20px;
    fill: black;
    transition-duration: 0.3s;
  }

  .Btn:hover {
    color: transparent;
  }

  .Btn:hover .icon {
    right: 43%;
    margin: 0;
    padding: 0;
    border: none;
    transition-duration: 0.3s;
  }

  .Btn:active {
    transform: translate(3px, 3px);
    transition-duration: 0.3s;
    box-shadow: 2px 2px 0px rgb(0, 0, 0);
  }
`;

const columns: ColumnDef<UserTransactionsType>[] = [
  {
    accessorKey: "description",
    header: "Description",
  },
  {
    accessorKey: "wallet_name",
    header: "Wallet Name",
  },
  {
    accessorKey: "date",
    header: () => <div className="text-center">Date</div>,
    cell: ({ row }: { row: any }) => {
      const date: string = row.getValue("date");
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
    header: () => <div className="text-center">Type</div>,
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
