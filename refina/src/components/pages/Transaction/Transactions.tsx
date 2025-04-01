import ExpenseButton from "@/components/ui/expense-button";
import IncomeButton from "@/components/ui/income-button";
import { useNavigate } from "react-router";

export default function Transactions() {
  const navigate = useNavigate();
  return (
    <div className="font-poppins min-h-screen w-full md:px-4">
      <div className="flex items-start justify-between gap-4 md:items-center">
        <h1 className="text-3xl font-semibold md:text-4xl">Your Transaction</h1>
        <div className="flex gap-5 items-center">
          <ExpenseButton onclick={() => navigate("/transactions/add/expense")} />
          <IncomeButton onclick={() => navigate("/transactions/add/income")} />
        </div>
      </div>
    </div>
  );
}
