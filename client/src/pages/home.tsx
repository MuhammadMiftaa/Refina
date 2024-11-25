import { useEffect, useState } from "react";
import Template from "../components/layouts/template";
import moment from "moment";
import { FaRegArrowAltCircleDown, FaRegArrowAltCircleUp } from "react-icons/fa";

export default function HomePage() {
  const [time, setTime] = useState("");
  useEffect(() => {
    const updateTime = () => {
      const currentTime = moment().format("LT");
      setTime(currentTime);
    };

    updateTime();

    const interval = setInterval(updateTime, 1000);

    return () => clearInterval(interval);
  }, []);


  return (
    <Template>
      <div className="bg-secondary h-screen w-full text-white py-4">
        <div className="flex justify-between items-center border-b border-slate-600 pb-4">
          <div className="flex gap-4 items-center px-4">
            <img
              src="/ava.jpeg"
              alt=""
              className="h-12 w-12 rounded-full border border-purple-800"
            />
            <div>
              <h1 className="text-[0.65rem] text-slate-400 font-poppins">
                madmifta77@gmail.com{" "}
                <span className="py-0.5 px-1 rounded-md text-black font-bold bg-green-500 uppercase">
                  ✅Verified
                </span>
              </h1>
              <h2 className="text-md font-semibold font-poppins">
                Muhammad Mifta
              </h2>
            </div>
          </div>
          <div className="rounded-xl bg-transparent border border-purple-600 text-purple-600 py-1 px-3 mr-6">
            <h1 className="font-poppins font-bold text-xl">{time}</h1>
          </div>
        </div>
        <div className="pt-7 px-6">
          <h1 className="font-poppins font-semibold text-2xl">
            General Statistic
          </h1>
          <div className="flex gap-6 mt-7">
            <div className="w-72 rounded-xl border border-slate-500 flex flex-col p-6 font-poppins bg-gradient-to-t from-[rgba(168,85,247,0.1)] via-transparent to-transparent">
              <h1 className="text-xs text-slate-400">Monthly Income</h1>
              <h2 className="text-2xl mt-5 flex">
                Rp. 3.500.000{" "}
                <span className="text-red-500 flex items-center gap-1 text-xs ml-4">
                  <FaRegArrowAltCircleDown />
                  10%
                </span>
              </h2>
            </div>
            <div className="w-72 rounded-xl border border-slate-500 flex flex-col p-6 font-poppins bg-gradient-to-t from-[rgba(168,85,247,0.1)] via-transparent to-transparent">
              <h1 className="text-xs text-slate-400">Monthly Outcome</h1>
              <h2 className="text-2xl mt-5 flex">
                Rp. 1.000.000{" "}
                <span className="text-green-500 flex items-center gap-1 text-xs ml-4">
                  <FaRegArrowAltCircleUp />
                  10%
                </span>
              </h2>
            </div>
            <div className="w-72 rounded-xl border border-slate-500 flex flex-col p-6 font-poppins bg-gradient-to-t from-[rgba(168,85,247,0.1)] via-transparent to-transparent">
              <h1 className="text-xs text-slate-400">Monthly Savings</h1>
              <h2 className="text-2xl mt-5 flex">
                Rp. 2.500.000{" "}
                <span className="text-green-500 flex items-center gap-1 text-xs ml-4">
                  <FaRegArrowAltCircleUp />
                  10%
                </span>
              </h2>
            </div>
          </div>
        </div>
      </div>
    </Template>
  );
}
