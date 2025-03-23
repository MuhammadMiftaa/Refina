import {
  Banknote,
  Check,
  Copy,
  Eye,
  EyeClosed,
  Plus,
  ReceiptText,
  Wallet,
} from "lucide-react";
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
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "../ui/accordion";
import Cookies from "js-cookie";
import { useQuery } from "@tanstack/react-query";
import { WalletType } from "@/types/Wallet";

async function fetchWallets() {
  const token = Cookies.get("token");

  const res = await fetch("http://localhost:8080/v1/users/wallets", {
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

export default function Wallets() {
  const navigate = useNavigate();
  const [search, setSearch] = useState("");
  const [showNumber, setShowNumber] = useState(false);
  const [isCopied, setIsCopied] = useState(false);
  const { data, isError, isLoading } = useQuery({
    queryKey: ["wallets"],
    queryFn: fetchWallets,
  });
  const Wallets: WalletType = data?.data ?? {
    user_id: "",
    name: "",
    email: "",
    wallets: [],
  };

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

  if (isLoading) return <p>Loading...</p>;
  if (isError) return <p>Terjadi kesalahan saat mengambil data.</p>;

  return (
    <div className="font-poppins min-h-screen w-full md:px-4">
      <div className="flex items-start md:items-center justify-between gap-4">
        <h1 className="text-3xl md:text-4xl font-semibold">Your Wallet</h1>
        <button
          className="cursor-pointer rounded-full bg-zinc-100 p-2 text-black/50 shadow duration-200 hover:text-black"
          onClick={() => navigate("/wallets/create")}
        >
          <Plus />
        </button>
      </div>

      <div className="mt-4 grid grid-cols-2 items-center gap-4 md:grid-cols-3 md:gap-8">
        <div className="flex w-full flex-col items-center gap-0 rounded-lg bg-zinc-100 p-2 shadow-md md:flex-row md:gap-4 md:rounded-2xl md:p-6">
          <Wallet className="w-8 md:w-fit" size={48} />
          <div>
            <h1 className="text-center text-2xl font-semibold text-nowrap md:text-left md:text-4xl">
              4
            </h1>
            <h2 className="text-center text-base font-light text-nowrap text-zinc-400 md:text-left md:text-xl">
              Total Wallets
            </h2>
          </div>
        </div>
        <div className="col-span-2 row-start-2 flex w-full flex-col items-center gap-0 rounded-lg bg-zinc-100 p-2 shadow-md md:col-span-1 md:col-start-2 md:row-start-1 md:flex-row md:gap-4 md:rounded-2xl md:p-6">
          <Banknote className="w-8 md:w-fit" size={48} />
          <div>
            <h1 className="text-center text-2xl font-semibold text-nowrap md:text-left md:text-4xl">
              RP 15.5 M
            </h1>
            <h2 className="text-center text-base font-light text-nowrap text-zinc-400 md:text-left md:text-xl">
              Total Balance
            </h2>
          </div>
        </div>
        <div className="flex w-full flex-col items-center gap-0 rounded-lg bg-zinc-100 p-2 shadow-md md:flex-row md:gap-4 md:rounded-2xl md:p-6">
          <ReceiptText className="w-8 md:w-fit" size={48} />
          <div>
            <h1 className="text-center text-2xl font-semibold text-nowrap md:text-left md:text-4xl">
              24
            </h1>
            <h2 className="text-center text-base font-light text-nowrap text-zinc-400 md:text-left md:text-xl">
              Total Transactions
            </h2>
          </div>
        </div>
      </div>

      <div className="mt-8 flex items-center justify-between gap-4">
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

      <div className="mt-8 grid grid-cols-1 gap-8 md:grid-cols-2">
        {Wallets?.wallets?.map((wallet) => (
          <div className="min-h-72 rounded-2xl bg-slate-100 p-4 shadow-xl">
            <h1 className="pt-2 pb-4 pl-2 text-2xl font-semibold">
              {wallet.name}
            </h1>
            <AspectRatio
              ratio={1.586 / 1}
              className="relative flex min-h-64 flex-col justify-between rounded-xl bg-linear-to-br/hsl from-indigo-500 to-teal-400 p-10"
            >
              <img
                className="absolute top-16 right-10 w-20"
                src="/mandiri.svg"
                alt=""
              />
              <div className="flex flex-col text-white">
                <h1 className="text-zinc-300">Balance</h1>
                <h2 className="text-4xl font-semibold">RP {wallet.balance}</h2>
              </div>
              <div className="flex flex-col">
                <div className="flex flex-col text-white">
                  <h1 className="text-zinc-300">Name</h1>
                  <h2 className="-mt-1 text-xl">{Wallets.name}</h2>
                </div>
                <div className="mt-5 flex w-full flex-row items-center justify-between text-white">
                  <div>
                    <h1 className="text-zinc-300">Account Number</h1>
                    <h2 className="-mt-1 text-xl">{wallet.number}</h2>
                  </div>
                  <div className="flex gap-4">
                    <button
                      onClick={() => setShowNumber(!showNumber)}
                      className="relative h-6 w-6"
                    >
                      {showNumber ? (
                        <Eye className="absolute top-0" />
                      ) : (
                        <EyeClosed className="absolute top-0" />
                      )}
                    </button>
                    <button
                      onClick={() => setIsCopied(!isCopied)}
                      className="relative h-6 w-6"
                    >
                      {isCopied ? (
                        <Check className="absolute top-0" />
                      ) : (
                        <Copy className="absolute top-0" />
                      )}
                    </button>
                  </div>
                </div>
              </div>
            </AspectRatio>
            <Accordion type="single" collapsible>
              <AccordionItem value="item-1">
                <AccordionTrigger>Is it accessible?</AccordionTrigger>
                <AccordionContent>
                  Lorem ipsum dolor sit amet consectetur adipisicing elit. Dolor
                  adipisci reiciendis accusamus consectetur nobis doloribus!
                  Esse nostrum consectetur quas veniam at ullam possimus
                  reiciendis earum molestiae error voluptatum praesentium, fuga,
                  magni dicta sed iusto facilis labore repellendus dolor rerum
                  et, laboriosam illo pariatur! Mollitia provident ut modi,
                  eligendi optio ducimus, laboriosam, consectetur harum aut
                  neque a. Incidunt eius repellendus cum dicta error iusto,
                  dolore dolorum rem? Asperiores rem pariatur tempora voluptates
                  praesentium eum tenetur, sapiente aliquam libero ipsa soluta,
                  maxime officiis beatae itaque. Provident nostrum magnam
                  deserunt illum labore ullam vero hic aut, eligendi est
                  laudantium consectetur ipsa corporis placeat!
                </AccordionContent>
              </AccordionItem>
            </Accordion>
          </div>
        ))}
      </div>
    </div>
  );
}
