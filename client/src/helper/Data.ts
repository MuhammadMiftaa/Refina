import {
  ArrowLeftRight,
  ChartNoAxesCombined,
  Landmark,
  Wallet,
} from "lucide-react";

export const data = {
  user: {
    name: "refina",
    email: "user@refina.com",
    avatar: "/avatars/shadcn.jpg",
  },
  // teams: [
  //   {
  //     name: "Refina",
  //     logo: <img src="/logo.png" alt="Logo Refina" />,
  //     plan: "Verified",
  //   },
  //   {
  //     name: "Acme Corp.",
  //     logo: <img src="/logo.png" alt="Logo Refina" />,
  //     plan: "Startup",
  //   },
  //   {
  //     name: "Evil Corp.",
  //     logo: <img src="/logo.png" alt="Logo Refina" />,
  //     plan: "Free",
  //   },
  // ],
  navMain: [
    {
      title: "Analytics",
      url: "/",
      icon: ChartNoAxesCombined,
    },
    {
      title: "Wallets",
      url: "/wallets",
      icon: Wallet,
    },
    {
      title: "Transactions",
      url: "/transactions",
      icon: ArrowLeftRight,
    },
    {
      title: "Investments",
      url: "/investments",
      icon: Landmark,
    },
  ],
  projects: [
    // {
    //   name: "Analytics",
    //   url: "/",
    //   icon: ChartNoAxesCombined,
    // },
    // {
    //   name: "Wallets",
    //   url: "/wallets",
    //   icon: Wallet,
    // },
    // {
    //   name: "Transactions",
    //   url: "/transactions",
    //   icon: ArrowLeftRight,
    // },
    // {
    //   name: "Investments",
    //   url: "/investments",
    //   icon: Landmark,
    // },
  ],
};

export const INVESTMENT_TYPE = {
  Gold: {
    id: "gold",
    label: "Gold",
    color: "bg-yellow-300",
  },
  Bonds: {
    id: "bonds",
    label: "Bonds",
    color: "bg-green-300",
  },
  Stocks: {
    id: "stocks",
    label: "Stocks",
    color: "bg-blue-300",
  },
  Deposits: {
    id: "deposits",
    label: "Deposits",
    color: "bg-purple-300",
  },
  "Mutual Funds": {
    id: "mutual_funds",
    label: "Mutual Funds",
    color: "bg-pink-300",
  },
  "Government Securities": {
    id: "government_securities",
    label: "Government Securities",
    color: "bg-red-300",
  },
  Others: {
    id: "others",
    label: "Others",
    color: "bg-gray-300",
  },
} as const;
