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
import { Input } from "../../ui/input";
import { useEffect, useState } from "react";
import { AspectRatio } from "../../ui/aspect-ratio";
import Cookies from "js-cookie";
import { useQuery } from "@tanstack/react-query";
import { WalletType } from "@/types/UserWallet";
import { formatCurrency, handleCopy, shortenMoney } from "@/helper/Helper";
import { TransactionType } from "@/types/UserTransaction";
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

  const res = await fetch(`${backendURL}/users/transactions`, {
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

  const { data: walletData, isLoading: walletLoading } = useQuery({
    queryKey: ["wallets"],
    queryFn: fetchWallets,
  });
  const { data: transactionData, isLoading: transactionLoading } = useQuery({
    queryKey: ["transactions"],
    queryFn: fetchTransactions,
  });

  const Wallets: WalletType[] = walletData?.data ?? [];

  const Transactions: TransactionType[] = transactionData?.data ?? [];

  const [type, setType] = useState("all");
  const [search, setSearch] = useState("");
  const [showNumbers, setShowNumbers] = useState<{ [key: string]: boolean }>(
    {},
  );
  const [copiedStates, setCopiedStates] = useState<{ [key: string]: boolean }>(
    {},
  );
  const [wallets, setWallets] = useState(Wallets);
  const [transactions, setTransactions] = useState(Transactions);

  useEffect(() => {
    setWallets(Wallets);
  }, [Wallets]);

  useEffect(() => {
    setTransactions(Transactions);
  }, [Transactions]);

  useEffect(() => {
    if (wallets) {
      const filtered = Wallets.filter((wallet) =>
        type !== "all" ? wallet.wallet_type === type : true,
      ).filter((wallet) =>
        search !== ""
          ? wallet.wallet_name.toLowerCase().includes(search)
          : true,
      );

      setWallets(filtered);

      if (Transactions?.length > 0) {
        const filteredTransactions = Transactions.filter((transaction) =>
          type !== "all" ? transaction.wallet_type_name === type : true,
        );

        setTransactions(filteredTransactions);
      }
    }
  }, [type, search]);

  const toggleShowNumber = (id: string) => {
    setShowNumbers((prev) => ({
      ...prev,
      [id]: !prev[id], // Toggle showNumber untuk wallet tertentu
    }));
  };

  const copyWalletNumber = (id: string, number: string) => {
    handleCopy(number);
    setCopiedStates((prev) => ({
      ...prev,
      [id]: true, // Set copied hanya untuk wallet tertentu
    }));
    setTimeout(() => {
      setCopiedStates((prev) => ({
        ...prev,
        [id]: false, // Reset copied setelah 2 detik
      }));
    }, 2000);
  };

  if (walletLoading || transactionLoading) return <Skeleton />;

  return (
    <div className="font-poppins min-h-screen w-full p-4 md:p-6">
      <div className="flex items-start justify-between gap-4 md:items-center">
        <h1 className="text-3xl font-semibold md:text-4xl">Your Wallet</h1>
        <button
          className="cursor-pointer rounded-full bg-zinc-100 p-2 text-black/50 shadow duration-200 hover:text-black"
          onClick={() => navigate("/wallets/create")}
        >
          <Plus />
        </button>
      </div>
      {wallets && wallets?.length > 0 ? (
        <>
          <div className="mt-4 grid grid-cols-2 items-center gap-4 md:grid-cols-3 md:gap-8">
            <div className="flex w-full flex-col items-center gap-0 rounded-lg bg-zinc-100 p-2 shadow-md md:flex-row md:gap-4 md:rounded-2xl md:p-6">
              <Wallet className="w-8 md:w-fit" size={48} />
              <div>
                <h1 className="text-center text-2xl font-semibold text-nowrap md:text-left md:text-4xl">
                  {wallets.length}
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
                  RP{" "}
                  {shortenMoney(
                    wallets.reduce(
                      (acc, wallet) => acc + wallet.wallet_balance,
                      0,
                    ),
                  )}
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
                  {transactions.length}
                </h1>
                <h2 className="text-center text-base font-light text-nowrap text-zinc-400 md:text-left md:text-xl">
                  Total Transactions
                </h2>
              </div>
            </div>
          </div>

          <div className="mt-8 flex items-center justify-between gap-4">
            <Select onValueChange={(value) => setType(value)}>
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
              onChange={(event) => setSearch(event.target.value.toLowerCase())}
              className="max-w-sm shadow"
              value={search}
            />
          </div>

          <div className="mt-8 grid grid-cols-1 gap-8 md:grid-cols-2">
            {wallets.map((wallet) => (
              <div
                key={wallet.id}
                className="min-h-72 rounded-2xl bg-slate-100 p-4 shadow-xl"
              >
                <h1 className="pt-2 pb-4 pl-2 text-2xl font-semibold">
                  {wallet.wallet_name}
                </h1>
                <AspectRatio
                  ratio={1.586 / 1}
                  className="relative flex min-h-64 flex-col justify-between rounded-xl bg-linear-to-br/hsl from-indigo-500 to-teal-400 p-10"
                >
                  {/* <img
                  className="absolute top-16 right-10 w-20"
                  src="/mandiri.svg"
                  alt=""
                /> */}
                  <div className="flex flex-col text-white">
                    <h1 className="text-zinc-300">Balance</h1>
                    <h2 className="flex items-center justify-between gap-2 text-4xl font-semibold">
                      RP{" "}
                      {showNumbers[wallet.id]
                        ? formatCurrency(wallet.wallet_balance)
                        : "*****"}
                      <button
                        onClick={() => toggleShowNumber(wallet.id)}
                        className="relative h-6 w-6"
                      >
                        {showNumbers[wallet.id] ? (
                          <Eye className="absolute top-0" />
                        ) : (
                          <EyeClosed className="absolute top-0" />
                        )}
                      </button>
                    </h2>
                  </div>
                  <div className="flex flex-col">
                    <div className="flex flex-col text-white">
                      <h1 className="text-zinc-300">Name</h1>
                      <h2 className="-mt-1 text-xl">{wallet.wallet_name}</h2>
                    </div>
                    <div className="mt-5 flex w-full flex-row items-center justify-between text-white">
                      <div>
                        <h1 className="text-zinc-300">Account Number</h1>
                        <h2 className="-mt-1 text-xl">
                          {wallet.wallet_number === "—"
                            ? wallet.wallet_number
                            : wallet.wallet_number.slice(0, 4) +
                              "–" +
                              wallet.wallet_number.slice(4, 8) +
                              "–" +
                              wallet.wallet_number.slice(8)}
                        </h2>
                      </div>
                      <div className="flex gap-4">
                        <button
                          onClick={() =>
                            copyWalletNumber(wallet.id, wallet.wallet_number)
                          }
                          className="relative h-6 w-6"
                        >
                          {copiedStates[wallet.id] ? (
                            <Check className="absolute top-0" />
                          ) : (
                            <Copy className="absolute top-0" />
                          )}
                        </button>
                      </div>
                    </div>
                  </div>
                </AspectRatio>
                {/* <Accordion type="single" collapsible>
                <AccordionItem value="item-1">
                  <AccordionTrigger>Is it accessible?</AccordionTrigger>
                  <AccordionContent>
                    Lorem ipsum dolor sit amet consectetur adipisicing elit.
                    Dolor adipisci reiciendis accusamus consectetur nobis
                    doloribus! Esse nostrum consectetur quas veniam at ullam
                    possimus reiciendis earum molestiae error voluptatum
                    praesentium, fuga, magni dicta sed iusto facilis labore
                    repellendus dolor rerum et, laboriosam illo pariatur!
                    Mollitia provident ut modi, eligendi optio ducimus,
                    laboriosam, consectetur harum aut neque a. Incidunt eius
                    repellendus cum dicta error iusto, dolore dolorum rem?
                    Asperiores rem pariatur tempora voluptates praesentium eum
                    tenetur, sapiente aliquam libero ipsa soluta, maxime
                    officiis beatae itaque. Provident nostrum magnam deserunt
                    illum labore ullam vero hic aut, eligendi est laudantium
                    consectetur ipsa corporis placeat!
                  </AccordionContent>
                </AccordionItem>
              </Accordion> */}
              </div>
            ))}
          </div>
        </>
      ) : (
        <NotFound />
      )}
    </div>
  );
}

function Skeleton() {
  return (
    <div className="min-h-screen w-full">
      <div className="h-15 w-full animate-pulse rounded-md bg-gray-200 md:w-100" />
      <div className="mt-4 grid grid-cols-2 gap-4 md:grid-cols-3">
        <div className="h-30 w-full animate-pulse rounded-md bg-gray-200" />
        <div className="h-30 w-full animate-pulse rounded-md bg-gray-200" />
        <div className="col-span-2 h-30 w-full animate-pulse rounded-md bg-gray-200 md:col-span-1" />
      </div>
      <div className="mt-4 flex justify-between gap-4">
        <div className="h-10 w-full animate-pulse rounded-md bg-gray-200 md:w-50" />
        <div className="h-10 w-full animate-pulse rounded-md bg-gray-200 md:w-50" />
      </div>
      <div className="mt-4 grid grid-cols-1 gap-4 md:grid-cols-2">
        <div className="h-72 w-full animate-pulse rounded-md bg-gray-200" />
        <div className="h-72 w-full animate-pulse rounded-md bg-gray-200" />
      </div>
    </div>
  );
}

function NotFound() {
  return (
    <div className="flex h-[80vh] w-full flex-col items-center justify-center">
      <img className="h-50" src="/assets/notfound.svg" alt="wallet not found" />
      <h1 className="-mt-5 text-lg font-light md:mt-0 md:text-2xl">
        No matching wallets found.
      </h1>
    </div>
  );
}
