import { X } from "lucide-react";
import { FaMoneyBillWave } from "react-icons/fa";
import { useNavigate } from "react-router";
import { SlWallet } from "react-icons/sl";
import { RxCrumpledPaper } from "react-icons/rx";
import { CiBank } from "react-icons/ci";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../../ui/select";
import { useQuery } from "@tanstack/react-query";
import { WalletTypeType } from "@/types/UserWallet";
import Cookies from "js-cookie";
import { useState } from "react";
import { SubmitButton } from "../../ui/submit-button";
import { CancelButton } from "../../ui/cancel-button";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

async function fetchWalletTypes() {
  const token = Cookies.get("token");
  const res = await fetch(`${import.meta.env.VITE_API_URL}/wallet-types`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });

  if (!res.ok) {
    throw new Error("Failed to fetch wallet types");
  }

  return res.json();
}

const WalletForm = z.object({
  wallet_type_id: z.string(),
  name: z.string(),
  number: z.string(),
  balance: z.number().default(0),
});

type WalletFormType = z.infer<typeof WalletForm>;

export default function CreateWallet() {
  const navigate = useNavigate();
  const form = useForm<WalletFormType>({
    resolver: zodResolver(WalletForm),
    defaultValues: {
      number: "—",
      balance: 0,
    }
  });
  const { data } = useQuery({
    queryKey: ["wallet-types"],
    queryFn: fetchWalletTypes,
  });
  const WalletTypes: WalletTypeType[] = data?.data ?? [];
  const Type = [
    { value: "physical", label: "Physical", icon: <FaMoneyBillWave /> },
    { value: "bank", label: "Bank", icon: <CiBank /> },
    { value: "e-wallet", label: "E-Wallet", icon: <SlWallet /> },
    { value: "others", label: "Others", icon: <RxCrumpledPaper /> },
  ];

  const [type, setType] = useState("");
  const [zeroBalance, setZeroBalance] = useState(false);

  const onSubmit = form.handleSubmit(async (data) => {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/wallets`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${Cookies.get("token")}`,
      },
      body: JSON.stringify({
        ...data,
        balance: zeroBalance ? 0 : data.balance,
      }),
    });

    const res = await response.json();

    if (response.ok) navigate(-1);
    else console.error(res.message);
  });

  return (
    <div className="font-poppins min-h-screen w-full md:px-4">
      <div className="flex items-start justify-between gap-4 md:items-center">
        <h1 className="text-3xl font-semibold md:text-4xl">Add a New Wallet</h1>
        <button
          className="cursor-pointer rounded-full bg-zinc-100 p-2 text-black/50 shadow duration-200 hover:text-black"
          onClick={() => navigate(-1)}
        >
          <X />
        </button>
      </div>

      <form
        className="mt-8 flex w-full flex-col gap-4 md:gap-10"
        onSubmit={onSubmit}
      >
        <div className="flex w-full flex-col">
          <label className="mb-2" htmlFor="type">
            Type
          </label>
          <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
            {Type.map((type) => (
              <label
                htmlFor={type.value}
                className="flex cursor-pointer items-center gap-2 rounded border border-gray-200 px-6 py-3 shadow-[4px_4px_0px_rgba(0,0,0,0.05)] duration-200 hover:bg-zinc-100 active:translate-x-1 active:translate-y-1 active:bg-gray-200 active:shadow-none has-checked:translate-x-1 has-checked:translate-y-1 has-checked:bg-zinc-100 has-checked:shadow-none"
                key={type.value}
              >
                <input
                  onChange={(e) => setType(e.target.value)}
                  className="hidden"
                  type="radio"
                  id={type.value}
                  name="type"
                  value={type.value}
                />
                {type.icon}
                <span>{type.label}</span>
              </label>
            ))}
          </div>
        </div>

        <div className="flex w-full flex-col gap-4 md:flex-row">
          <div className="flex w-full flex-col">
            <label className="mb-2" htmlFor="walletType">
              Wallet Type
            </label>
            <Select
              disabled={!type}
              onValueChange={(value) => form.setValue("wallet_type_id", value)}
            >
              <SelectTrigger className="h-fit min-h-11.5 w-full px-4 py-2 text-lg shadow-md">
                <SelectValue
                  className="text-lg"
                  placeholder="Choose a Wallet Type"
                />
              </SelectTrigger>
              <SelectContent>
                {WalletTypes.filter((wt) => wt.type === type).map((wt) => (
                  <SelectItem key={wt.id} className="text-lg" value={wt.id}>
                    {wt.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>

            <input type="hidden" {...form.register("wallet_type_id")} />
          </div>

          <div className="flex w-full flex-col">
            <label className="mb-2" htmlFor="number">
              Number
            </label>
            <input
              className="w-full rounded-lg border border-gray-200 px-4 py-2 text-lg shadow-md"
              type="text"
              id="number"
              placeholder="Wallet Number"
              {...form.register("number")}
            />
          </div>
        </div>

        <div className="flex w-full flex-col">
          <label className="mb-2" htmlFor="name">
            Name
          </label>
          <input
            className="w-full rounded-lg border border-gray-200 px-4 py-2 text-lg shadow-md"
            type="text"
            id="name"
            placeholder="Wallet Name"
            {...form.register("name")}
          />
        </div>

        <div className="flex w-full flex-col">
          <label className="mb-2" htmlFor="balance">
            Balance
          </label>
          <input
            className="w-full rounded-lg border border-gray-200 px-4 py-2 text-lg shadow-md disabled:cursor-not-allowed disabled:bg-gray-200"
            type="number"
            id="balance"
            disabled={zeroBalance}
            value={zeroBalance ? 0 : form.watch("balance")}
            placeholder="Wallet Balance"
            {...form.register("balance", { valueAsNumber: true })}
          />
          <label className="group relative mt-2 flex cursor-pointer items-center justify-end">
            <input
              onChange={() => setZeroBalance(!zeroBalance)}
              className="peer sr-only"
              type="checkbox"
            />
            <div className="h-5 w-5 rounded border-2 border-purple-500 bg-white from-purple-500 to-pink-500 transition-all duration-300 ease-in-out peer-checked:rotate-0 peer-checked:border-0 peer-checked:bg-gradient-to-br after:absolute after:top-1/2 after:left-1/2 after:h-4 after:w-4 after:-translate-x-1/2 after:-translate-y-1/2 after:bg-[url('data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9IiNmZmZmZmYiIHN0cm9rZS13aWR0aD0iMyIgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIiBzdHJva2UtbGluZWpvaW49InJvdW5kIj48cG9seWxpbmUgcG9pbnRzPSIyMCA2IDkgMTcgNCAxMiI+PC9wb2x5bGluZT48L3N2Zz4=')] after:bg-contain after:bg-no-repeat after:opacity-0 after:transition-opacity after:duration-300 after:content-[''] peer-checked:after:opacity-100 hover:shadow-[0_0_15px_rgba(168,85,247,0.5)]"></div>
            <span className="ml-3 text-sm font-medium text-gray-900">
              Set balance to zero
            </span>
          </label>
          {form.formState.errors.balance && (
            <span className="mt-2 text-sm text-red-500">
              {form.formState.errors.balance.message}
            </span>
          )}
        </div>

        <div className="flex w-full items-center justify-center gap-4">
          <CancelButton
            text="Clear Form"
            onclick={() => {
              form.reset();
              setZeroBalance(false);
              setType("");
              form.reset({ balance: 0 });
            }}
          />
          <SubmitButton text="Add Wallet" />
        </div>
      </form>
    </div>
  );
}
