import { X } from "lucide-react";
import { useNavigate, useParams } from "react-router";
import Cookies from "js-cookie";
import { useQuery } from "@tanstack/react-query";
import TextField from "@mui/material/TextField";
import Autocomplete from "@mui/material/Autocomplete";
import { CategoryType } from "@/types/Category";
import { WalletType } from "@/types/UserWallet";
import { useEffect, useState } from "react";
import { formatCurrency } from "@/helper/Helper";
import styled from "styled-components";
import { NumericFormat } from "react-number-format";
import { CancelButton } from "@/components/ui/cancel-button";
import { SubmitButton } from "@/components/ui/submit-button";

async function fetchCategories(type: string) {
  const token = Cookies.get("token");

  const res = await fetch("http://localhost:8080/v1/categories/type/" + type, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });

  if (!res.ok) {
    throw new Error("Failed to fetch categories");
  }

  return res.json();
}

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

export default function AddTransaction() {
  const navigate = useNavigate();
  const { type } = useParams();
  const { data: categoriesData, isLoading: categoriesLoading } = useQuery({
    queryKey: ["categories"],
    queryFn: () => fetchCategories(type as string),
  });
  const { data: walletsData, isLoading: walletsLoading } = useQuery({
    queryKey: ["wallets"],
    queryFn: fetchWallets,
  });

  const Categories: CategoryType[] = categoriesData?.data ?? [];
  const Wallets: WalletType = walletsData?.data ?? {
    user_id: "",
    name: "",
    email: "",
    wallets: [],
  };

  const [categories, setCategories] = useState([
    {
      id: "",
      name: "",
      group_name: "",
    },
  ]);
  const [wallets, setWallets] = useState(Wallets.wallets);
  const [userInput, setUserInput] = useState({
    amount: 0,
    wallet_id: "",
    category_id: "",
    date: new Date(),
    description: "",
  });

  useEffect(() => {
    const flatMap = Categories.flatMap((group) =>
      group.category.map((item) => ({
        ...item,
        group_name: group.group_name,
      })),
    );

    setCategories(flatMap);
  }, [Categories]);

  useEffect(() => {
    setWallets(Wallets.wallets);
  }, [Wallets]);

  if (categoriesLoading || walletsLoading) {
    return <div>Loading...</div>;
  }

  return (
    <div className="font-poppins min-h-screen w-full md:px-4">
      <div className="flex items-start justify-between gap-4 md:items-center">
        <h1 className="text-3xl font-semibold md:text-4xl">Add Transaction</h1>
        <button
          className="cursor-pointer rounded-full bg-zinc-100 p-2 text-black/50 shadow duration-200 hover:text-black"
          onClick={() => navigate(-1)}
        >
          <X />
        </button>
      </div>

      <form className="font-poppins mt-8 flex w-full flex-col gap-4 md:gap-10">
        <div className="flex w-full flex-col">
          <label className="mb-2" htmlFor="type">
            Type
          </label>
          <Autocomplete
            className="rounded-lg border-gray-200 shadow-md"
            options={categories.sort(
              (a, b) => -b.group_name.localeCompare(a.group_name),
            )}
            groupBy={(option) => option.group_name}
            getOptionLabel={(option) => option.name}
            sx={{
              "& .MuiOutlinedInput-root": {
                borderRadius: "8px", // Sesuai dengan rounded-lg di Tailwind
                fontFamily: "Poppins, sans-serif",
                "&:hover .MuiOutlinedInput-notchedOutline": {
                  borderColor: "#4f46e5", // Warna hover indigo-600
                },
                "&.Mui-focused .MuiOutlinedInput-notchedOutline": {
                  borderColor: "#4f46e5", // Warna focus indigo-600
                  borderWidth: "2px",
                },
              },
              "& .MuiInputLabel-root": {
                fontFamily: "Poppins, sans-serif",
                color: "#6b7280", // Warna label gray-500
                "&.Mui-focused": {
                  color: "#4f46e5", // Warna label saat focus
                },
              },
            }}
            onChange={(_, newValue) => {
              setUserInput((prev) => ({
                ...prev,
                category_id: newValue?.id || "",
              }));
            }}
            renderInput={(params) => (
              <TextField
                className="font-poppins"
                {...params}
                label="Transaction type"
              />
            )}
            renderGroup={(params) => (
              <li key={params.key}>
                <h1 className="font-poppins pt-2 pl-2 text-sm font-semibold text-indigo-600">
                  {params.group}
                </h1>
                <h2 className="font-poppins">{params.children}</h2>
              </li>
            )}
          />
        </div>

        <div className="flex w-full flex-col">
          <label className="mb-2">Amount (IDR)</label>
          <NumericFormat
            value={userInput.amount}
            onValueChange={(values) => {
              console.info("DANCOK", userInput);
              const { floatValue } = values;
              setUserInput((prev) => ({
                ...prev,
                amount: floatValue || 0,
              }));
            }}
            customInput={TextField}
            valueIsNumericString
            thousandSeparator=","
            prefix="Rp. "
            sx={{
              "& .MuiOutlinedInput-root": {
                borderRadius: "8px", // Sesuai dengan rounded-lg di Tailwind
                fontFamily: "Poppins, sans-serif",
                fontSize: "2rem",
                textAlign: "center",
                "&:hover .MuiOutlinedInput-notchedOutline": {
                  borderColor: "#4f46e5", // Warna hover indigo-600
                },
                "&.Mui-focused .MuiOutlinedInput-notchedOutline": {
                  borderColor: "#4f46e5", // Warna focus indigo-600
                  borderWidth: "2px",
                },
              },
              "& .MuiInputLabel-root": {
                fontFamily: "Poppins, sans-serif",
                color: "#6b7280", // Warna label gray-500
                "&.Mui-focused": {
                  color: "#4f46e5", // Warna label saat focus
                },
              },
            }}
          />
        </div>

        <div className="flex w-full flex-col">
          <label className="mb-2">Wallets</label>
          <div className="grid grid-cols-2 gap-4">
            {wallets.map((wallet) => (
              <label
                htmlFor={wallet.id}
                className="relative flex cursor-pointer flex-col gap-2 rounded border border-indigo-200 px-6 py-3 shadow-[4px_4px_0px_oklch(0.87_0.065_274.039)] duration-200 hover:bg-indigo-100 active:translate-x-1 active:translate-y-1 active:bg-indigo-200 active:shadow-none has-checked:translate-x-1 has-checked:translate-y-1 has-checked:bg-indigo-100 has-checked:shadow-none"
                key={wallet.id}
              >
                {/* <input
                  className="hidden"
                  onChange={(e) =>
                    setUserInput((prev) => ({
                      ...prev,
                      wallet_id: e.target.value,
                    }))
                  }
                  type="radio"
                  id={wallet.id}
                  name="type"
                  value={wallet.id}
                /> */}
                <CheckboxStyle>
                  <label className="neon-checkbox">
                    <input
                      onChange={(e) =>
                        setUserInput((prev) => ({
                          ...prev,
                          wallet_id: e.target.value,
                        }))
                      }
                      type="radio"
                      id={wallet.id}
                      name="type"
                      value={wallet.id}
                    />
                    <div className="neon-checkbox__frame">
                      <div className="neon-checkbox__box">
                        <div className="neon-checkbox__check-container">
                          <svg
                            viewBox="0 0 24 24"
                            className="neon-checkbox__check"
                          >
                            <path d="M3,12.5l7,7L21,5" />
                          </svg>
                        </div>
                        <div className="neon-checkbox__glow" />
                        <div className="neon-checkbox__borders">
                          <span />
                          <span />
                          <span />
                          <span />
                        </div>
                      </div>
                      <div className="neon-checkbox__effects">
                        <div className="neon-checkbox__particles">
                          <span />
                          <span />
                          <span />
                          <span /> <span />
                          <span />
                          <span />
                          <span /> <span />
                          <span />
                          <span />
                          <span />
                        </div>
                        <div className="neon-checkbox__rings">
                          <div className="ring" />
                          <div className="ring" />
                          <div className="ring" />
                        </div>
                        <div className="neon-checkbox__sparks">
                          <span />
                          <span />
                          <span />
                          <span />
                        </div>
                      </div>
                    </div>
                  </label>
                </CheckboxStyle>
                <img
                  className="absolute top-5 right-5 w-10"
                  src="/mandiri.svg"
                  alt=""
                />
                <div className="">
                  <span className="line-clamp-1">{wallet.name}</span>
                  <span className="line-clamp-1">{wallet.number}</span>
                  <div className="flex flex-col items-end">
                    <span className="text-right font-semibold">
                      RP {formatCurrency(wallet.balance)}
                    </span>
                    <span className="-mt-2 text-sm text-zinc-600">Balance</span>
                  </div>
                </div>
              </label>
            ))}
          </div>
        </div>

        <div className="flex w-full flex-col">
          <label htmlFor="">Description</label>
          <input
            className="w-full rounded-lg border border-gray-200 px-4 py-2 text-lg shadow-md"
            type="text"
            id="name"
            placeholder="Transaction Description"
            onChange={(e) =>
              setUserInput((prev) => ({
                ...prev,
                description: e.target.value,
              }))
            }
          />
        </div>

        <div className="flex w-full items-center justify-center gap-4">
          <CancelButton
            text="Clear Form"
            // onclick={() => {
            //   form.reset();
            //   setZeroBalance(false);
            //   setType("");
            //   form.reset({ balance: 0 });
            // }}
          />
          <SubmitButton text="Add Transaction" />
        </div>
      </form>
    </div>
  );
}

const CheckboxStyle = styled.div`
  .neon-checkbox {
    --primary: #00ffaa;
    --primary-dark: #00cc88;
    --primary-light: #88ffdd;
    --size: 30px;
    position: absolute;
    bottom: 10px;
    left: 10px;
    width: var(--size);
    height: var(--size);
    cursor: pointer;
    -webkit-tap-highlight-color: transparent;
  }

  .neon-checkbox input {
    display: none;
  }

  .neon-checkbox__frame {
    position: relative;
    width: 100%;
    height: 100%;
  }

  .neon-checkbox__box {
    position: absolute;
    inset: 0;
    background: white;
    border-radius: 4px;
    border: 2px solid var(--color-indigo-200);
    transition: all 0.4s ease;
  }

  .neon-checkbox__check-container {
    position: absolute;
    inset: 2px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .neon-checkbox__check {
    width: 80%;
    height: 80%;
    fill: none;
    stroke: black;
    stroke-width: 3;
    stroke-linecap: round;
    stroke-linejoin: round;
    stroke-dasharray: 40;
    stroke-dashoffset: 40;
    transform-origin: center;
    transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
  }

  .neon-checkbox__glow {
    position: absolute;
    inset: -2px;
    border-radius: 6px;
    background: var(--color-indigo-200);
    opacity: 0;
    filter: blur(8px);
    transform: scale(1.2);
    transition: all 0.4s ease;
  }

  .neon-checkbox__borders {
    position: absolute;
    inset: 0;
    border-radius: 4px;
    overflow: hidden;
  }

  .neon-checkbox__borders span {
    position: absolute;
    width: 40px;
    height: 1px;
    background: var(--primary);
    opacity: 0;
    transition: opacity 0.4s ease;
  }

  .neon-checkbox__borders span:nth-child(1) {
    top: 0;
    left: -100%;
    animation: borderFlow1 2s linear infinite;
  }

  .neon-checkbox__borders span:nth-child(2) {
    top: -100%;
    right: 0;
    width: 1px;
    height: 40px;
    animation: borderFlow2 2s linear infinite;
  }

  .neon-checkbox__borders span:nth-child(3) {
    bottom: 0;
    right: -100%;
    animation: borderFlow3 2s linear infinite;
  }

  .neon-checkbox__borders span:nth-child(4) {
    bottom: -100%;
    left: 0;
    width: 1px;
    height: 40px;
    animation: borderFlow4 2s linear infinite;
  }

  .neon-checkbox__particles span {
    position: absolute;
    width: 4px;
    height: 4px;
    background: var(--primary);
    border-radius: 50%;
    opacity: 0;
    pointer-events: none;
    top: 50%;
    left: 50%;
    box-shadow: 0 0 6px var(--primary);
  }

  .neon-checkbox__rings {
    position: absolute;
    inset: -20px;
    pointer-events: none;
  }

  .neon-checkbox__rings .ring {
    position: absolute;
    inset: 0;
    border-radius: 50%;
    border: 1px solid var(--primary);
    opacity: 0;
    transform: scale(0);
  }

  .neon-checkbox__sparks span {
    position: absolute;
    width: 20px;
    height: 1px;
    background: linear-gradient(90deg, var(--primary), transparent);
    opacity: 0;
  }

  /* Hover Effects */
  .neon-checkbox:hover .neon-checkbox__box {
    border-color: black;
    transform: scale(1.05);
  }

  /* Checked State */
  .neon-checkbox input:checked ~ .neon-checkbox__frame .neon-checkbox__box {
    border-color: black;
    background: rgba(0, 255, 170, 0.1);
  }

  .neon-checkbox input:checked ~ .neon-checkbox__frame .neon-checkbox__check {
    stroke-dashoffset: 0;
    transform: scale(1.1);
  }

  .neon-checkbox input:checked ~ .neon-checkbox__frame .neon-checkbox__glow {
    opacity: 0.2;
  }

  .neon-checkbox
    input:checked
    ~ .neon-checkbox__frame
    .neon-checkbox__borders
    span {
    opacity: 1;
  }

  /* Particle Animations */
  .neon-checkbox
    input:checked
    ~ .neon-checkbox__frame
    .neon-checkbox__particles
    span {
    animation: particleExplosion 0.6s ease-out forwards;
  }

  .neon-checkbox
    input:checked
    ~ .neon-checkbox__frame
    .neon-checkbox__rings
    .ring {
    animation: ringPulse 0.6s ease-out forwards;
  }

  .neon-checkbox
    input:checked
    ~ .neon-checkbox__frame
    .neon-checkbox__sparks
    span {
    animation: sparkFlash 0.6s ease-out forwards;
  }

  /* Animations */
  @keyframes borderFlow1 {
    0% {
      transform: translateX(0);
    }
    100% {
      transform: translateX(200%);
    }
  }

  @keyframes borderFlow2 {
    0% {
      transform: translateY(0);
    }
    100% {
      transform: translateY(200%);
    }
  }

  @keyframes borderFlow3 {
    0% {
      transform: translateX(0);
    }
    100% {
      transform: translateX(-200%);
    }
  }

  @keyframes borderFlow4 {
    0% {
      transform: translateY(0);
    }
    100% {
      transform: translateY(-200%);
    }
  }

  @keyframes particleExplosion {
    0% {
      transform: translate(-50%, -50%) scale(1);
      opacity: 0;
    }
    20% {
      opacity: 1;
    }
    100% {
      transform: translate(
          calc(-50% + var(--x, 20px)),
          calc(-50% + var(--y, 20px))
        )
        scale(0);
      opacity: 0;
    }
  }

  @keyframes ringPulse {
    0% {
      transform: scale(0);
      opacity: 1;
    }
    100% {
      transform: scale(2);
      opacity: 0;
    }
  }

  @keyframes sparkFlash {
    0% {
      transform: rotate(var(--r, 0deg)) translateX(0) scale(1);
      opacity: 1;
    }
    100% {
      transform: rotate(var(--r, 0deg)) translateX(30px) scale(0);
      opacity: 0;
    }
  }

  /* Particle Positions */
  .neon-checkbox__particles span:nth-child(1) {
    --x: 25px;
    --y: -25px;
  }
  .neon-checkbox__particles span:nth-child(2) {
    --x: -25px;
    --y: -25px;
  }
  .neon-checkbox__particles span:nth-child(3) {
    --x: 25px;
    --y: 25px;
  }
  .neon-checkbox__particles span:nth-child(4) {
    --x: -25px;
    --y: 25px;
  }
  .neon-checkbox__particles span:nth-child(5) {
    --x: 35px;
    --y: 0px;
  }
  .neon-checkbox__particles span:nth-child(6) {
    --x: -35px;
    --y: 0px;
  }
  .neon-checkbox__particles span:nth-child(7) {
    --x: 0px;
    --y: 35px;
  }
  .neon-checkbox__particles span:nth-child(8) {
    --x: 0px;
    --y: -35px;
  }
  .neon-checkbox__particles span:nth-child(9) {
    --x: 20px;
    --y: -30px;
  }
  .neon-checkbox__particles span:nth-child(10) {
    --x: -20px;
    --y: 30px;
  }
  .neon-checkbox__particles span:nth-child(11) {
    --x: 30px;
    --y: 20px;
  }
  .neon-checkbox__particles span:nth-child(12) {
    --x: -30px;
    --y: -20px;
  }

  /* Spark Rotations */
  .neon-checkbox__sparks span:nth-child(1) {
    --r: 0deg;
    top: 50%;
    left: 50%;
  }
  .neon-checkbox__sparks span:nth-child(2) {
    --r: 90deg;
    top: 50%;
    left: 50%;
  }
  .neon-checkbox__sparks span:nth-child(3) {
    --r: 180deg;
    top: 50%;
    left: 50%;
  }
  .neon-checkbox__sparks span:nth-child(4) {
    --r: 270deg;
    top: 50%;
    left: 50%;
  }

  /* Ring Delays */
  .neon-checkbox__rings .ring:nth-child(1) {
    animation-delay: 0s;
  }
  .neon-checkbox__rings .ring:nth-child(2) {
    animation-delay: 0.1s;
  }
  .neon-checkbox__rings .ring:nth-child(3) {
    animation-delay: 0.2s;
  }
`;
