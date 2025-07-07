import { BsArrowDownLeftCircle, BsArrowUpRightCircle } from "react-icons/bs";
import { PiArrowsLeftRightLight } from "react-icons/pi";

export default function CategoryType({ type }: { type: string }) {
  const commonClass =
    "flex justify-center items-center gap-2 relative z-10 overflow-hidden rounded-lg border-2 px-7 py-2 text-[15px] leading-[1.4em] tracking-wide capitalize transition-all duration-300";

  const label = type.replace(/_/g, " ");

  if (type === "expense") {
    return (
      <div className="group relative flex items-center justify-center">
        <button
          className={`${commonClass} border-rose-500 bg-gradient-to-r from-rose-500/10 via-transparent to-rose-500/10 text-rose-500 shadow-[inset_0_0_10px_rgba(244,63,94,0.4),_0_0_9px_3px_rgba(244,63,94,0.1)] group-hover:shadow-[inset_0_0_10px_rgba(244,63,94,0.6),_0_0_9px_3px_rgba(244,63,94,0.2)]`}
        >
          {label}
          <BsArrowUpRightCircle />
        </button>
      </div>
    );
  }

  if (type === "income") {
    return (
      <div className="group relative flex items-center justify-center">
        <button
          className={`${commonClass} border-blue-500 bg-gradient-to-r from-blue-500/10 via-transparent to-blue-500/10 text-blue-500 shadow-[inset_0_0_10px_rgba(59,130,246,0.4),_0_0_9px_3px_rgba(59,130,246,0.1)] group-hover:shadow-[inset_0_0_10px_rgba(59,130,246,0.6),_0_0_9px_3px_rgba(59,130,246,0.2)]`}
        >
          {label}
          <BsArrowDownLeftCircle />
        </button>
      </div>
    );
  }

  if (type === "fund_transfer") {
    return (
      <div className="group relative flex items-center justify-center">
        <button
          className={`${commonClass} border-indigo-500 bg-gradient-to-r from-indigo-500/10 via-transparent to-indigo-500/10 text-indigo-500 shadow-[inset_0_0_10px_rgba(99,102,241,0.4),_0_0_9px_3px_rgba(99,102,241,0.1)] group-hover:shadow-[inset_0_0_10px_rgba(99,102,241,0.6),_0_0_9px_3px_rgba(99,102,241,0.2)]`}
        >
          {label}
          <PiArrowsLeftRightLight />
        </button>
      </div>
    );
  }

  return (
    <div className="group relative flex items-center justify-center">
      <button
        className={`${commonClass} border-gray-500 bg-gradient-to-r from-gray-500/10 via-transparent to-gray-500/10 text-gray-500 shadow-[inset_0_0_10px_rgba(107,114,128,0.4),_0_0_9px_3px_rgba(107,114,128,0.1)] group-hover:shadow-[inset_0_0_10px_rgba(107,114,128,0.6),_0_0_9px_3px_rgba(107,114,128,0.2)]`}
      >
        {label}
      </button>
    </div>
  );
}
