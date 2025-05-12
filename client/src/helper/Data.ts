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
