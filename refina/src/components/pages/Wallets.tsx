import { Banknote, Check, Copy, Eye, EyeClosed, Plus, ReceiptText, Wallet } from "lucide-react";
import { useNavigate } from "react-router";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "../ui/input";
import { useState } from "react";
import { AspectRatio } from "../ui/aspect-ratio";

export default function Wallets() {
  const navigate = useNavigate();
  const [search, setSearch] = useState("");
  const [showNumber, setShowNumber] = useState(false);
  const [isCopied, setIsCopied] = useState(false);

  const WalletType = [
    {
      value: "all",
      label: "All",
    },
    {
      value: "physical",
      label: "Physical",
    },
    {
      value: "bank",
      label: "Bank",
    },
    {
      value: "e-wallet",
      label: "E-Wallet",
    },
    {
      value: "others",
      label: "Others",
    },
  ];

  return (
    <div className="w-full min-h-screen font-poppins px-4">
      <div className="flex justify-between items-center">
        <h1 className="text-4xl font-semibold">Your Wallet</h1>
        <button
          className="bg-zinc-100 p-2 rounded-full text-black/50 hover:text-black duration-200 cursor-pointer shadow"
          onClick={() => navigate("#")}
        >
          <Plus />
        </button>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-3 items-center gap-4 md:gap-8 mt-4">
        <div className="p-2 md:p-6 w-full bg-zinc-100 rounded-lg md:rounded-2xl flex flex-col md:flex-row items-center gap-0 md:gap-4 shadow-md">
          <Wallet className="w-8 md:w-fit" size={48} />
          <div>
            <h1 className="text-center md:text-left text-nowrap text-2xl md:text-4xl font-semibold">
              4
            </h1>
            <h2 className="text-center md:text-left text-nowrap text-base md:text-xl text-zinc-400 font-light">
              Total Wallets
            </h2>
          </div>
        </div>
        <div className="row-start-2 md:row-start-1 col-span-2 md:col-span-1 md:col-start-2 p-2 md:p-6 w-full bg-zinc-100 rounded-lg md:rounded-2xl flex flex-col md:flex-row items-center gap-0 md:gap-4 shadow-md">
          <Banknote className="w-8 md:w-fit" size={48} />
          <div>
            <h1 className="text-center md:text-left text-nowrap text-2xl md:text-4xl font-semibold">
              RP 15.5 M
            </h1>
            <h2 className="text-center md:text-left text-nowrap text-base md:text-xl text-zinc-400 font-light">
              Total Balance
            </h2>
          </div>
        </div>
        <div className="p-2 md:p-6 w-full bg-zinc-100 rounded-lg md:rounded-2xl flex flex-col md:flex-row items-center gap-0 md:gap-4 shadow-md">
          <ReceiptText className="w-8 md:w-fit" size={48} />
          <div>
            <h1 className="text-center md:text-left text-nowrap text-2xl md:text-4xl font-semibold">
              24
            </h1>
            <h2 className="text-center md:text-left text-nowrap text-base md:text-xl text-zinc-400 font-light">
              Total Transactions
            </h2>
          </div>
        </div>
      </div>

      <div className="flex justify-between items-center mt-8 gap-4">
        <Select>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Wallet Type" />
          </SelectTrigger>
          <SelectContent>
            {WalletType.map((item) => (
              <SelectItem key={item.value} value={item.value}>
                {item.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        <Input
          placeholder="Find Wallet..."
          onChange={(event) => setSearch(event.target.value)}
          className="max-w-sm shadow"
        />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mt-8">
        <div className="bg-slate-100 rounded-2xl p-4 min-h-72 shadow-xl">
          <h1 className="text-2xl font-semibold pt-2 pb-4 pl-2">
            Tabungan Mandiri
          </h1>
          <AspectRatio ratio={1.586 / 1} className="relative min-h-64 rounded-xl bg-linear-to-br/hsl from-indigo-500 to-teal-400 p-10 flex flex-col justify-between">
            <img className="w-20 absolute right-10 top-16" src="/mandiri.svg" alt="" />
            <div className="text-white flex flex-col">
              <h1 className="text-zinc-300">Balance</h1>
              <h2 className="text-4xl font-semibold ">RP 700.000,00</h2>
            </div>
            <div className="flex flex-col">
              <div className="text-white flex flex-col">
                <h1 className="text-zinc-300">Name</h1>
                <h2 className="text-xl -mt-1">Muhammad Miftakul Salam</h2>
              </div>
              <div className="text-white flex flex-row justify-between items-center w-full mt-5">
                <div>
                  <h1 className="text-zinc-300">Account Number</h1>
                  <h2 className="text-xl -mt-1">1410-****-11281</h2>
                </div>
                <div className="flex gap-4">
                  <button onClick={() => setShowNumber(!showNumber)} className="relative w-6 h-6">
                    { showNumber ? <Eye className="absolute top-0" /> : <EyeClosed className="absolute top-0" />}
                  </button>
                  <button onClick={() => setIsCopied(!isCopied)} className="relative w-6 h-6">
                    { isCopied ? <Check className="absolute top-0" /> : <Copy className="absolute top-0" />}
                  </button>
                </div>
              </div>
            </div>
          </AspectRatio>
        </div>
      </div>
    </div>
  );
}
