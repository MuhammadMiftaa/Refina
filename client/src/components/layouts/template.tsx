import { ReactNode } from "react";
import { BsCreditCard2BackFill } from "react-icons/bs";
import { FaChevronRight } from "react-icons/fa";
import { FaChartSimple } from "react-icons/fa6";
import { GoArrowSwitch } from "react-icons/go";
import { ImExit } from "react-icons/im";
import { IoHome } from "react-icons/io5";

export default function Template(props: { children: ReactNode }) {
  return (
    <div className="">
      <div className="bg-secondary text-white fixed h-screen w-[20%] py-6 flex flex-col border-r border-slate-600">
        <div className="flex gap-3 items-center border-b border-slate-600 pb-4 mb-4 px-5">
          <img src="/logo.png" alt="logo" className="h-10 w-10 rounded-lg" />
          <div className="flex flex-col -mt-1">
            <h1 className="font-bold uppercase">
              Refina
            </h1>
            <h2 className="font-light text-zinc-300 text-xs">
              Financial Management
            </h2>
          </div>
        </div>
        <div className="flex flex-col px-5">
          <h1 className="uppercase text-zinc-400 text-xs">Menu</h1>
          <div className="flex flex-col gap-1 mt-2">
            <div className="flex p-2 rounded-lg items-center justify-between text-purple-500 bg-primary">
              <div className="flex gap-2 items-center">
                <div className="">
                  <IoHome />
                </div>
                <h1 className="text-sm">Dashboard</h1>
              </div>
              <div className="text-xs text-purple-500 scale-100">
                <FaChevronRight />
              </div>
            </div>
            <div className="flex p-2 rounded-lg items-center justify-between text-zinc-400 hover:bg-primary hover:text-purple-500 duration-500 cursor-pointer group">
              <div className="flex gap-2 items-center">
                <div className="">
                  <GoArrowSwitch />
                </div>
                <h1 className="font-light text-sm">Transactions</h1>
              </div>
              <div className="text-xs text-purple-500 scale-0 group-hover:scale-100 duration-500">
                <FaChevronRight />
              </div>
            </div>
            <div className="flex p-2 rounded-lg items-center justify-between text-zinc-400 hover:bg-primary hover:text-purple-500 duration-500 cursor-pointer group">
              <div className="flex gap-2 items-center">
                <div className="">
                  <FaChartSimple />
                </div>
                <h1 className="font-light text-sm">Analytics</h1>
              </div>
              <div className="text-xs text-purple-500 scale-0 group-hover:scale-100 duration-500">
                <FaChevronRight />
              </div>
            </div>
            <div className="flex p-2 rounded-lg items-center justify-between text-zinc-400 hover:bg-primary hover:text-purple-500 duration-500 cursor-pointer group">
              <div className="flex gap-2 items-center">
                <div className="">
                  <BsCreditCard2BackFill />
                </div>
                <h1 className="font-light text-sm">Wallet</h1>
              </div>
              <div className="text-xs text-purple-500 scale-0 group-hover:scale-100 duration-500">
                <FaChevronRight />
              </div>
            </div>
          </div>
        </div>
        <div className="flex items-center absolute left-7 bottom-5 gap-2 text-zinc-400 hover:text-red-600 duration-500 cursor-pointer">
          <div>
            <ImExit />
          </div>
          <h1 className="-mt-1 uppercase text-sm">Logout</h1>
        </div>
      </div>
      <div className="pl-[20%]">{props.children}</div>
    </div>
  );
}