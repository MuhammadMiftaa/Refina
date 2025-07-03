import { useQuery } from "@tanstack/react-query";
import Cookies from "js-cookie";
import { FormEvent, useEffect, useState } from "react";
import { DataTable } from "./data-table";
import { ArrowUpDown, MoreHorizontal } from "lucide-react";
import { ColumnDef } from "@tanstack/react-table";
import { BsArrowDownLeftCircle, BsArrowUpRightCircle } from "react-icons/bs";
import {
  formatRawDate,
  generateColorByType,
  toLocalISOString,
} from "@/helper/Helper";
import { PiArrowsLeftRightLight } from "react-icons/pi";
import { getBackendURL } from "@/lib/readenv";
import { TransactionType } from "@/types/UserTransaction";
import GlobalLoading from "@/components/ui/global-loading";
import { IoCloseOutline } from "react-icons/io5";
import { useLocation, useNavigate } from "react-router";
import Autocomplete from "@mui/material/Autocomplete";
import { CategoryType } from "@/types/Category";
import { WalletsByTypeType } from "@/types/UserWallet";
import { TextField } from "@mui/material";
import dayjs from "dayjs";
import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs";
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider";
import { MobileDatePicker } from "@mui/x-date-pickers/MobileDatePicker";
import { NumericFormat } from "react-number-format";
import { FileUpload } from "@/components/ui/file-upload";
import { SubmitButton } from "@/components/ui/submit-button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { IoCopyOutline } from "react-icons/io5";
import { FaRegTrashAlt } from "react-icons/fa";

async function fetchTransactions() {
  const backendURL = getBackendURL();

  const token = Cookies.get("token");

  const res = await fetch(`${backendURL}/transactions/user`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });
  if (!res.ok) {
    throw new Error("Failed to fetch transactions");
  }

  return res.json();
}

async function fetchDetailTransaction(id: string) {
  const backendURL = getBackendURL();

  const token = Cookies.get("token");

  const res = await fetch(`${backendURL}/transactions/${id}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });
  if (!res.ok) {
    throw new Error("Failed to fetch transactions");
  }

  return res.json();
}

async function fetchCategories(type: string) {
  const backendURL = getBackendURL();

  const token = Cookies.get("token");

  const res = await fetch(`${backendURL}/categories/type/` + type, {
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
  const backendURL = getBackendURL();

  const token = Cookies.get("token");

  const res = await fetch(`${backendURL}/wallets/user-by-type`, {
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
  const backendURL = getBackendURL();
  const navigate = useNavigate();
  const location = useLocation();
  const searchParams = new URLSearchParams(location.search);
  const id = searchParams.get("id") || "";
  const type = searchParams.get("type") || "";

  const { data: transactionsData, isLoading: transactionsLoading } = useQuery({
    queryKey: ["transactions"],
    queryFn: fetchTransactions,
  });

  const {
    data: transactionData,
    isLoading: transactionLoading,
    refetch: refetchTransaction,
  } = useQuery({
    queryKey: ["detail transaction", id],
    queryFn: () => fetchDetailTransaction(id as string),
    enabled: false,
  });

  const {
    data: categoriesData,
    isLoading: categoriesLoading,
    refetch: refetchCategories,
  } = useQuery({
    queryKey: ["categories", type],
    queryFn: () => fetchCategories(type as string),
    enabled: false,
  });

  const {
    data: walletsData,
    isLoading: walletsLoading,
    refetch: refetchWallets,
  } = useQuery({
    queryKey: ["wallets"],
    queryFn: fetchWallets,
    enabled: false,
  });

  const Transactions: TransactionType[] = transactionsData?.data ?? [];
  const TransactionDetail: TransactionType = transactionData?.data ?? {};
  const Categories: CategoryType[] = categoriesData?.data ?? [];
  const Wallets: WalletsByTypeType[] = walletsData?.data ?? [];

  const [isOpen, setIsOpen] = useState(false);
  const [categories, setCategories] = useState([
    {
      id: "",
      name: "",
      group_name: "",
    },
  ]);
  const [wallets, setWallets] = useState([
    {
      id: "",
      name: "",
      number: "",
      balance: 0,
      type: "",
    },
  ]);
  type FileItem = File | string;
  const [files, setFiles] = useState<FileItem[]>([]);
  const [clearFiles, setClearFiles] = useState(false);
  const [userInput, setUserInput] = useState({
    id: "",
    amount: 0,
    wallet_id: "",
    category_id: "",
    date: toLocalISOString(new Date()),
    description: "",
    attachments: [""],
    // FUND TRANSFER
    from_wallet_id: "",
    to_wallet_id: "",
    admin_fee: 0,
  });
  const [isAutocompleteReady, setIsAutocompleteReady] = useState(false);

  useEffect(() => {
    if (Categories.length > 0) {
      const flatMap = Categories.flatMap((group) =>
        group.category?.map((item) => ({
          ...item,
          group_name: group.group_name,
        })),
      );
      setCategories(flatMap);
    }
  }, [categoriesData]);

  useEffect(() => {
    if (Wallets.length > 0) {
      const flatMap = Wallets.flatMap((group) =>
        group.wallets?.map((item) => ({
          ...item,
          type: group.type,
        })),
      );
      setWallets(flatMap);
    }
  }, [walletsData]);

  useEffect(() => {
    if (!transactionLoading && TransactionDetail.id) {
      setUserInput({
        id: TransactionDetail.id,
        amount: TransactionDetail.amount,
        wallet_id: TransactionDetail.wallet_id,
        category_id: TransactionDetail.category_id,
        date: toLocalISOString(new Date(TransactionDetail.transaction_date)),
        description: TransactionDetail.description,
        from_wallet_id: "",
        to_wallet_id: "",
        admin_fee: 0,
        attachments: TransactionDetail.attachments,
      });
    }
  }, [TransactionDetail.id, transactionLoading]);

  useEffect(() => {
    const hasValidCategory =
      categories.length > 0 &&
      categories.some((cat) => cat.id === userInput.category_id);
    if (
      !categoriesLoading &&
      !transactionLoading &&
      !walletsLoading
      // hasValidCategory
    ) {
      setIsAutocompleteReady(true);
    }
    console.log("Categories loading:", categoriesLoading);
    console.log("Transaction loading:", transactionLoading);
    console.log("Wallets loading:", walletsLoading);
    console.log("Categories:", hasValidCategory);
    console.log("User input category ID:", userInput);
    console.log("User input category ID cok:", TransactionDetail.category_id);
    console.log("User input category ID hancok:", userInput.category_id);
    console.log("Categories:", categories);
  }, [
    categoriesLoading,
    transactionLoading,
    walletsLoading,
    categories,
    userInput.category_id,
  ]);

  useEffect(() => {
    if (isOpen && id) {
      refetchTransaction();
    }
    if (isOpen && type) {
      refetchCategories();
    }
    if (isOpen) {
      refetchWallets();
    }
  }, [isOpen, id, type]);

  useEffect(() => {
    setIsAutocompleteReady(false);
  }, [id, type]);

  useEffect(() => {
    if (!transactionLoading && TransactionDetail.id) {
      if (
        TransactionDetail.attachments &&
        TransactionDetail.attachments.length > 0
      ) {
        setFiles(TransactionDetail.attachments);
      } else {
        setFiles([]);
      }
      setClearFiles(false);
    }
  }, [TransactionDetail.attachments, transactionLoading]);

  const handleIsOpen = () => {
    if (!isOpen) {
      // Reset state dulu agar tidak ada stale state
      setIsAutocompleteReady(false);
      setUserInput({
        id: "",
        amount: 0,
        wallet_id: "",
        category_id: "",
        date: toLocalISOString(new Date()),
        description: "",
        from_wallet_id: "",
        to_wallet_id: "",
        admin_fee: 0,
        attachments: [],
      });

      navigate(""); // Bersihkan query param
    }

    setIsOpen((prev) => !prev);
  };

  const handleFileUpload = (file: File[]) => {
    const Files = [...files, ...file];
    setFiles(Files);
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const token = Cookies.get("token") || "";

    const data =
      type === "fund_transfer"
        ? {
            ...userInput,
            cash_in_category_id: categories.find((c) => c.name === "Cash In")
              ?.id,
            cash_out_category_id: categories.find((c) => c.name === "Cash Out")
              ?.id,
          }
        : userInput;

    try {
      const res = await fetch(`${backendURL}/transactions/${userInput.id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(data),
      });

      if (!res.ok) {
        throw new Error("Failed to add transaction");
      }

      // const resData = await res.json();

      // if (type === "fund_transfer") {
      //   await uploadAttachment(
      //     resData.data.cash_in_transaction_id,
      //     files,
      //     token,
      //   );
      //   await uploadAttachment(
      //     resData.data.cash_out_transaction_id,
      //     files,
      //     token,
      //   );
      // } else {
      //   await uploadAttachment(resData.data.id, files, token);
      // }

      setIsOpen(false);
      navigate("/transactions");
    } catch (error) {
      console.error("Error creating transaction:", error);
    }
  };

  if (transactionsLoading) {
    return <GlobalLoading />;
  }

  return (
    <div className="font-poppins min-h-screen w-screen p-4 md:w-full md:p-6">
      <div className="flex flex-col items-start justify-between gap-4 md:flex-row md:items-center">
        <h1 className="text-3xl font-semibold md:text-4xl">Your Transaction</h1>
      </div>

      <div className="mt-6 rounded-2xl">
        <DataTable
          columns={columns}
          data={Transactions}
          setIsOpen={handleIsOpen}
        />
      </div>

      <div
        className={`inset-0 bg-gray-700 duration-300 md:fixed ${isOpen ? "z-[98] opacity-50" : "-z-[97] opacity-0"}`}
        onClick={handleIsOpen}
      />
      <div
        className={`fixed top-0 bottom-0 w-screen bg-white duration-500 md:w-[30vw] ${isOpen ? "right-0 z-[98]" : "-right-[100vw] md:-right-[30vw]"}`}
      >
        <div className="flex w-full items-center justify-between p-4">
          <h1 className="text-xl">Edit your transaction</h1>
          <button className="cursor-pointer text-3xl" onClick={handleIsOpen}>
            <IoCloseOutline className="" />
          </button>
        </div>
        <div className="h-1 w-full bg-zinc-800" />
        <div className="h-full w-full">
          {isAutocompleteReady ? (
            <form
              className="relative flex flex-col gap-4 p-6"
              onSubmit={(e) => handleSubmit(e)}
            >
              {/* Transaction Category */}
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
                        "&::before": {
                          display: "none",
                        },
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
                  value={
                    categories.find(
                      (cat) => cat.id === userInput.category_id,
                    ) || null
                  }
                  renderInput={(params) => (
                    <TextField
                      className="font-poppins"
                      {...params}
                      label="Transaction type"
                    />
                  )}
                  renderGroup={(params) => (
                    <li className="z-[100]" key={params.key}>
                      <h1 className="font-poppins pt-2 pl-2 text-sm font-semibold text-indigo-600 capitalize">
                        {params.group.replace(/-/g, " ")}
                      </h1>
                      <h2 className="font-poppins">{params.children}</h2>
                    </li>
                  )}
                />
              </div>

              {/* Transaction Wallet */}
              <div className="flex w-full flex-col">
                <label className="mb-2" htmlFor="type">
                  Wallets
                </label>
                <Autocomplete
                  className="rounded-lg border-gray-200 shadow-md"
                  options={wallets.sort(
                    (a, b) => -b.type.localeCompare(a.type),
                  )}
                  groupBy={(option) => option.type}
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
                        "&::before": {
                          display: "none",
                        },
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
                      wallet_id: newValue?.id || "",
                    }));
                  }}
                  value={
                    wallets?.find(
                      (wallet) => wallet?.id === userInput.wallet_id,
                    ) || null
                  }
                  renderInput={(params) => (
                    <TextField
                      className="font-poppins"
                      {...params}
                      label="Transaction type"
                    />
                  )}
                  renderGroup={(params) => (
                    <li className="z-[100]" key={params.key}>
                      <h1 className="font-poppins pt-2 pl-2 text-sm font-semibold text-indigo-600">
                        {params.group}
                      </h1>
                      <h2 className="font-poppins">{params.children}</h2>
                    </li>
                  )}
                />
              </div>

              {/* Transaction Date */}
              <div className="flex w-full flex-col">
                <label className="mb-2" htmlFor="date">
                  Date
                </label>
                <LocalizationProvider dateAdapter={AdapterDayjs}>
                  <MobileDatePicker
                    defaultValue={dayjs(userInput.date)}
                    onChange={(value) => {
                      const date = value?.toDate() || new Date();
                      setUserInput((prev) => ({
                        ...prev,
                        date: toLocalISOString(date),
                      }));
                    }}
                    slotProps={{
                      textField: {
                        className: "font-poppins shadow-md bg-transparent",
                        sx: {
                          "& .MuiFormControl-root": {
                            backgroundColor: "transparent",
                          },
                          "& .MuiPickersTextField-root": {
                            backgroundColor: "transparent",
                          },
                          "& .MuiPickersInputBase-root": {
                            borderRadius: "8px !important", // <== penting
                          },
                          "& .MuiOutlinedInput-root": {
                            borderRadius: "8px !important", // <== penting
                            fontFamily: "Poppins, sans-serif",
                          },
                          "& .MuiOutlinedInput-notchedOutline": {
                            borderRadius: "8px !important", // <== ini bagian yang menampilkan border-nya
                            borderColor: "#d1d5db", // gray-300 default
                          },
                          "& .MuiOutlinedInput-root:hover .MuiOutlinedInput-notchedOutline":
                            {
                              borderColor: "#4f46e5", // hover indigo-600
                            },
                          "& .MuiOutlinedInput-root.Mui-focused .MuiOutlinedInput-notchedOutline":
                            {
                              borderColor: "#4f46e5", // focus indigo-600
                              borderWidth: "2px",
                            },
                          "& .MuiInputLabel-root": {
                            fontFamily: "Poppins, sans-serif",
                            color: "#6b7280",
                            "&.Mui-focused": {
                              color: "#4f46e5",
                            },
                          },
                        },
                      },
                    }}
                  />
                </LocalizationProvider>
              </div>

              {/* Transaction Amount */}
              <div className="flex w-full flex-col">
                <label className="mb-2">Amount (IDR)</label>
                <NumericFormat
                  className="shadow-md"
                  value={userInput.amount}
                  onValueChange={(values) => {
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
                      fontSize: {
                        xs: "1.25rem", // Mobile (0px and up)
                        md: "1.25rem", // Desktop (900px and up)
                      },
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

              {/* Transaction Description */}
              <div className="flex w-full flex-col">
                <label className="mb-2">Description</label>
                <input
                  className="w-full rounded-lg border border-gray-200 px-4 py-3.5 text-lg shadow-md focus:border-indigo-600 focus:ring-1 focus:ring-indigo-600 focus:outline-none"
                  type="text"
                  id="name"
                  placeholder="Transaction Description"
                  value={userInput.description}
                  onChange={(e) =>
                    setUserInput((prev) => ({
                      ...prev,
                      description: e.target.value,
                    }))
                  }
                />
              </div>

              {/* File Upload Section */}
              <div className="flex min-h-96 w-full flex-col items-center justify-center">
                <div className="font-poppins flex w-full flex-col items-center justify-center gap-2">
                  <p className="relative z-20 text-center text-base font-bold text-neutral-700 dark:text-neutral-300">
                    Upload Receipt/Invoice (optional)
                  </p>
                  <p className="relative z-20 text-center text-base font-normal text-neutral-400 dark:text-neutral-400">
                    Drag or drop your files here or click to upload
                  </p>
                </div>
                <FileUpload
                  onChange={handleFileUpload}
                  clearFiles={clearFiles}
                />
              </div>

              <div className="absolute right-10 bottom-10 left-10">
                <SubmitButton text="Save Transaction" />
              </div>
            </form>
          ) : (
            <div className="flex h-full w-full items-start justify-center">
              <div className="">
                <GlobalLoading />
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

const columns: ColumnDef<TransactionType>[] = [
  {
    accessorKey: "description",
    header: ({ column }) => (
      <button
        className="flex items-center"
        onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
      >
        Description
        <ArrowUpDown className="ml-2 h-4 w-4 cursor-pointer" />
      </button>
    ),
  },
  {
    accessorKey: "wallet_type",
    header: ({ column }) => (
      <button
        className="flex items-center"
        onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
      >
        Wallet Type
        <ArrowUpDown className="ml-2 h-4 w-4 cursor-pointer" />
      </button>
    ),
  },
  {
    accessorKey: "transaction_date",
    header: ({ column }) => (
      <button
        onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        className="mx-auto flex items-center justify-center text-center"
      >
        Date
        <ArrowUpDown className="ml-2 h-4 w-4 cursor-pointer" />
      </button>
    ),
    cell: ({ row }: { row: any }) => {
      const date: string = row.getValue("transaction_date");
      const formattedDate = formatRawDate(date);

      return (
        <div className="flex flex-col items-center">
          {/* <h1 className="font-light">{formattedDate[1]}</h1> */}
          <p className="text-sm text-nowrap text-zinc-500">
            {/* {formattedDate[0]}, {formattedDate[2]} */}
            {formattedDate}
          </p>
        </div>
      );
    },
  },
  {
    accessorKey: "category_type",
    header: ({ column }) => (
      <button
        className="mx-auto flex items-center"
        onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
      >
        Category
        <ArrowUpDown className="ml-2 h-4 w-4 cursor-pointer" />
      </button>
    ),
    cell: ({ row }: { row: any }) => {
      const type: string = row.getValue("category_type");
      return (
        <h1
          className={`mx-auto flex w-fit items-center gap-2 rounded-2xl border px-3 py-1 font-light text-nowrap uppercase text-${generateColorByType(type)} border-${generateColorByType(type)}`}
        >
          {type.replace(/_/g, " ")}{" "}
          <span className={`text-${generateColorByType(type)}`}>
            {type === "expense" ? (
              <BsArrowUpRightCircle />
            ) : type === "income" ? (
              <BsArrowDownLeftCircle />
            ) : (
              <PiArrowsLeftRightLight />
            )}
          </span>
        </h1>
      );
    },
  },
  {
    accessorKey: "category_name",
    header: ({ column }) => (
      <button
        className="flex items-center"
        onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
      >
        Category
        <ArrowUpDown className="ml-2 h-4 w-4 cursor-pointer" />
      </button>
    ),
  },
  {
    accessorKey: "amount",
    header: ({ column }) => (
      <button
        className="flex w-full items-center justify-end"
        onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
      >
        Amount
        <ArrowUpDown className="ml-2 h-4 w-4 cursor-pointer" />
      </button>
    ),
    cell: ({ row }: { row: any }) => {
      const amount = parseFloat(row.getValue("amount"));
      const formatted = new Intl.NumberFormat("id-ID", {
        style: "currency",
        currency: "IDR",
        minimumFractionDigits: 0,
      }).format(amount);

      return <div className="text-right font-medium">{formatted}</div>;
    },
  },
  {
    id: "actions",
    cell: ({ row }: { row: any }) => {
      const transaction = row.original;
      const backendURL = getBackendURL();
      const [deleteConfirm, setDeleteConfirm] = useState(false);

      const deleteTransaction = async (id: string) => {
        const token = Cookies.get("token") || "";

        try {
          const res = await fetch(`${backendURL}/transactions/${id}`, {
            method: "DELETE",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
          });

          if (!res.ok) {
            throw new Error("Failed to add transaction");
          }

          setDeleteConfirm(false);
        } catch (error) {
          console.error("Error deleting transaction:", error);
          setDeleteConfirm(false);
        }
      };

      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button className="mx-auto flex h-8 w-8 cursor-pointer items-center justify-center p-0">
              <MoreHorizontal className="h-4 w-4" />
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="font-poppins border">
            {/* <DropdownMenuLabel>Actions</DropdownMenuLabel> */}
            <DropdownMenuItem
              className="flex items-center justify-between gap-5 focus:bg-sky-500"
              onClick={(e) => {
                e.stopPropagation();
                navigator.clipboard.writeText(transaction.id);
              }}
            >
              Copy transaction ID
              <IoCopyOutline />
            </DropdownMenuItem>
            <DropdownMenuItem
              onClick={(e) => {
                e.stopPropagation();
                setDeleteConfirm(true);
              }}
              className="flex items-center justify-between gap-5 focus:bg-rose-500"
            >
              Delete Transaction
              <FaRegTrashAlt />
            </DropdownMenuItem>
          </DropdownMenuContent>

          <div
            className={`fixed inset-0 flex items-center justify-center bg-zinc-800/70 duration-100 ${deleteConfirm ? "z-[9999] opacity-100" : "-z-[9999] opacity-0"}`}
            onClick={(e) => {
              e.stopPropagation();
              setDeleteConfirm(false);
            }}
          >
            <div className="flex flex-col rounded-xl bg-gray-100 p-7 text-black">
              <h1 className="text-lg font-bold">
                Are you sure to delete this transaction?
              </h1>
              <hr />
              <p>This action cannot be canceled.</p>
              <div className="mt-4 flex justify-end gap-3">
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    setDeleteConfirm(false);
                  }}
                  className="cursor-pointer rounded-lg bg-zinc-300 px-6 py-3 duration-300 hover:bg-zinc-400"
                >
                  Cancel
                </button>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    deleteTransaction(transaction.id);
                  }}
                  className="flex cursor-pointer items-center justify-between gap-4 rounded-lg bg-rose-500 px-6 py-3 text-white duration-300 hover:bg-rose-600"
                >
                  Delete
                  <FaRegTrashAlt />
                </button>
              </div>
            </div>
          </div>
        </DropdownMenu>
      );
    },
  },
];
