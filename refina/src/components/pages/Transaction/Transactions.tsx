import { WalletType } from "@/types/UserWallet";
import { useQuery } from "@tanstack/react-query";
import Cookies from "js-cookie";
import { useEffect, useState } from "react";
import { FaMoneyBillTransfer } from "react-icons/fa6";
import { GiPayMoney, GiReceiveMoney } from "react-icons/gi";
import { useNavigate } from "react-router";
import styled from "styled-components";

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

export default function Transactions() {
  const navigate = useNavigate();

  const { data: walletData, isLoading: walletLoading } = useQuery({
    queryKey: ["wallets"],
    queryFn: fetchWallets,
  });

  const Wallets: WalletType = walletData?.data ?? {
    user_id: "",
    name: "",
    email: "",
    wallets: [],
  };

  const [wallets, setWallets] = useState(Wallets.wallets);

  useEffect(() => {
    setWallets(Wallets.wallets);
  }, [Wallets]);

  if (walletLoading) {
    return (
      <div className="flex h-screen w-full items-center justify-center">
        <div className="loader"></div>
      </div>
    );
  }

  return (
    <div className="font-poppins min-h-screen w-full md:px-4">
      <div className="flex items-start justify-between gap-4 md:items-center">
        <h1 className="text-3xl font-semibold md:text-4xl">Your Transaction</h1>
        {wallets.length > 0 && (
          <div className="flex items-center gap-5">
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
