
  import {
    Table,
    TableBody,
    TableCaption,
    TableCell,
    TableFooter,
    TableHead,
    TableHeader,
    TableRow,
  } from "../../components/ui/table";
  import { FiPlus } from "react-icons/fi";
  import Template from "../../components/layouts/template";
  import { LuHistory } from "react-icons/lu";
  import { BsArrowDownLeftCircle, BsArrowUpRightCircle } from "react-icons/bs";
  
  export default function TransactionsPage() {
    return (
      <Template>
        <div className="bg-secondary w-full text-white p-5">
          <div className="flex justify-between items-center bg-primary text-purple-500 rounded-2xl p-4">
            <div className="flex items-center gap-2">
              <LuHistory />
              <h1 className="text-2xl font-bold font-poppins">
                Transactions History
              </h1>
            </div>
            <button className="bg-[#00D47E] text-primary py-2 px-4 rounded-lg flex items-center font-poppins font-light">
              <FiPlus />
              <span className="ml-2">Add</span>
            </button>
          </div>
          <div className="bg-primary rounded-2xl p-4 mt-4">
            <Table className="font-poppins">
              <TableCaption>A list of your recent invoices.</TableCaption>
              <TableHeader>
                <TableRow className="border-none text-xs">
                  {/* <TableHead className="w-[60px]">ID</TableHead> */}
                  <TableHead>Desc</TableHead>
                  <TableHead>Date</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Category</TableHead>
                  <TableHead className="text-right">Amount</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {transactions.map((transaction) => (
                  <TableRow key={transaction.id} className="border-zinc-500">
                    {/* <TableCell className="py-5">{transaction.id}</TableCell> */}
                    <TableCell className="py-5 ">
                      {transaction.description}
                    </TableCell>
                    <TableCell className="py-5">{transaction.date}</TableCell>
                    <TableCell className={`py-5`}>
                      <h1
                        className={`flex items-center gap-2 font-light uppercase px-3 py-1 w-fit rounded-2xl border ${
                          transaction.transaction_type === "expense"
                            ? "border-red-500 "
                            : "border-green-500 "
                        }`}
                      >
                        {transaction.transaction_type}{" "}
                        <span
                          className={` ${
                            transaction.transaction_type === "expense"
                              ? "text-red-500 "
                              : "text-green-500 "
                          }`}
                        >
                          {transaction.transaction_type === "expense" ? (
                            <BsArrowUpRightCircle />
                          ) : (
                            <BsArrowDownLeftCircle />
                          )}
                        </span>
                      </h1>
                    </TableCell>
                    <TableCell className="py-5">{transaction.category}</TableCell>
                    <TableCell className="text-right">
                      {transaction.amount}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
              <TableFooter>
                <TableRow>
                  <TableCell colSpan={3}>Total</TableCell>
                  <TableCell className="text-right">$2,500.00</TableCell>
                </TableRow>
              </TableFooter>
            </Table>
          </div>
        </div>
      </Template>
    );
  }
  
  type transactionsType = {
    id: string;
    amount: number;
    transaction_type: "expense" | "income";
    date: string;
    description: string;
    category: string;
    user_id: string;
  };
  
  export const transactions: transactionsType[] = [
    {
      id: "1",
      amount: 100000,
      transaction_type: "expense",
      date: "2024-11-20",
      description: "Belanja di supermarket",
      category: "Groceries",
      user_id: "101",
    },
    {
      id: "2",
      amount: 50000,
      transaction_type: "expense",
      date: "2024-11-21",
      description: "Transportasi online",
      category: "Transport",
      user_id: "101",
    },
    {
      id: "3",
      amount: 200000,
      transaction_type: "income",
      date: "2024-11-21",
      description: "Gaji mingguan",
      category: "Salary",
      user_id: "101",
    },
    {
      id: "4",
      amount: 80000,
      transaction_type: "expense",
      date: "2024-11-22",
      description: "Makan malam",
      category: "Food & Drinks",
      user_id: "102",
    },
    {
      id: "5",
      amount: 150000,
      transaction_type: "income",
      date: "2024-11-22",
      description: "Bonus proyek",
      category: "Bonus",
      user_id: "101",
    },
    {
      id: "6",
      amount: 75000,
      transaction_type: "expense",
      date: "2024-11-23",
      description: "Beli buku",
      category: "Education",
      user_id: "102",
    },
    {
      id: "7",
      amount: 120000,
      transaction_type: "expense",
      date: "2024-11-24",
      description: "Tagihan listrik",
      category: "Utilities",
      user_id: "101",
    },
    {
      id: "8",
      amount: 200000,
      transaction_type: "income",
      date: "2024-11-24",
      description: "Penjualan barang bekas",
      category: "Other Income",
      user_id: "102",
    },
    {
      id: "9",
      amount: 50000,
      transaction_type: "expense",
      date: "2024-11-25",
      description: "Ngopi di cafe",
      category: "Food & Drinks",
      user_id: "101",
    },
    {
      id: "10",
      amount: 500000,
      transaction_type: "income",
      date: "2024-11-26",
      description: "Transfer keluarga",
      category: "Family",
      user_id: "103",
    },
    {
      id: "11",
      amount: 300000,
      transaction_type: "expense",
      date: "2024-11-26",
      description: "Servis motor",
      category: "Transport",
      user_id: "102",
    },
    {
      id: "12",
      amount: 100000,
      transaction_type: "income",
      date: "2024-11-27",
      description: "Pemasukan tambahan",
      category: "Other Income",
      user_id: "101",
    },
    {
      id: "13",
      amount: 120000,
      transaction_type: "expense",
      date: "2024-11-28",
      description: "Pulsa internet",
      category: "Utilities",
      user_id: "102",
    },
    {
      id: "14",
      amount: 45000,
      transaction_type: "expense",
      date: "2024-11-28",
      description: "Cemilan malam",
      category: "Food & Drinks",
      user_id: "103",
    },
    {
      id: "15",
      amount: 600000,
      transaction_type: "income",
      date: "2024-11-29",
      description: "Pendapatan investasi",
      category: "Investment",
      user_id: "101",
    },
    {
      id: "16",
      amount: 100000,
      transaction_type: "expense",
      date: "2024-11-29",
      description: "Langganan streaming",
      category: "Entertainment",
      user_id: "102",
    },
    {
      id: "17",
      amount: 85000,
      transaction_type: "expense",
      date: "2024-11-30",
      description: "Makan siang",
      category: "Food & Drinks",
      user_id: "103",
    },
    {
      id: "18",
      amount: 250000,
      transaction_type: "income",
      date: "2024-11-30",
      description: "Hadiah teman",
      category: "Gift",
      user_id: "101",
    },
    {
      id: "19",
      amount: 90000,
      transaction_type: "expense",
      date: "2024-12-01",
      description: "Beli pakaian",
      category: "Shopping",
      user_id: "102",
    },
    {
      id: "20",
      amount: 150000,
      transaction_type: "income",
      date: "2024-12-01",
      description: "Jual buku bekas",
      category: "Other Income",
      user_id: "103",
    },
  ];
  