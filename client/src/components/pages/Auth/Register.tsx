import { zodResolver } from "@hookform/resolvers/zod";
import { FormEventHandler, useState } from "react";
import { useForm } from "react-hook-form";
import { Link, Navigate, useNavigate } from "react-router-dom";
import { z } from "zod";
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
} from "@/components/ui/input-otp";
import { REGEXP_ONLY_DIGITS_AND_CHARS } from "input-otp";
import { getBackendURL } from "../../../lib/readenv";

const postFormSchema = z.object({
  name: z.string(),
  email: z.string().email(),
  password: z.string().min(8, "Password must be at least 8 characters"),
});

type PostFormSchema = z.infer<typeof postFormSchema>;

export default function Register(props: {
  isAuthenticated: boolean;
  handleLogin: () => void;
}) {
  const backendURL = getBackendURL();

  const navigate = useNavigate();

  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [submited, setSubmited] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(false);
  const [errorOTP, setErrorOTP] = useState("");

  const { register, handleSubmit, formState, watch } = useForm<PostFormSchema>({
    resolver: zodResolver(postFormSchema),
  });

  const onSubmit = handleSubmit(async (data) => {
    if (watch("password") !== confirmPassword) {
      setError("Password and confirm password must be the same");
      return;
    }

    setLoading(true);
    const res = await fetch(`${backendURL}/auth/register`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    }).then((res) => res.json());

    if (res.status) {
      localStorage.setItem("email", res.data.email);
      sendOTP();
    } else {
      setError(res.message);
    }
  });

  const sendOTP = async () => {
    const email = localStorage.getItem("email") || "";
    if (email === "") {
      return;
    }

    const res = await fetch(`${backendURL}/auth/send/otp`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email: email }),
    }).then((res) => res.json());

    if (res.status) {
      setLoading(false);
      setSubmited(true);
    } else {
      setError(res.message);
    }
  };

  const verifyOTP: FormEventHandler<HTMLFormElement> = async (e) => {
    e.preventDefault();

    setLoading(true);
    const email = localStorage.getItem("email") || "";
    const code = Array.from(document.querySelectorAll("input")).reduce(
      (acc, input) => acc + input.value,
      "",
    );

    if (email === "") {
      return;
    }

    const res = await fetch(`${backendURL}/auth/verify/otp`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email: email, otp: code }),
    }).then((res) => res.json());

    if (res.status) {
      setLoading(false);
      navigate("/login");
    } else {
      setLoading(false);
      setErrorOTP(res.message);
      setTimeout(() => {
        setErrorOTP("");
      }, 3000);
    }
  };

  return props.isAuthenticated ? (
    <Navigate to={"/"} />
  ) : (
    <div
      className="relative flex min-h-screen after:absolute after:inset-0 after:bg-black/40 after:content-['']"
      style={{
        backgroundImage: "url('/background.jpg')",
        backgroundSize: "cover",
        backgroundPosition: "center",
        backgroundRepeat: "no-repeat",
      }}
    >
      <div className="font-inter z-10 mx-auto my-20 flex min-h-96 w-[45%] flex-col justify-center rounded-xl bg-white/30 backdrop-blur">
        <div className="p-8">
          {submited ? (
            <>
              <h1 className="bg-gradient-to-r from-white from-20% to-purple-500 bg-clip-text pb-2 text-center text-4xl font-bold text-transparent">
                Verify Your Email Address
              </h1>
              <h2 className="text-center text-xl font-light text-zinc-300">
                Please enter the 6-digit code we sent to{" "}
              </h2>
              <h3 className="-mt-1 text-center text-xl font-bold text-white">
                {localStorage.getItem("email")}
              </h3>
              <form
                onSubmit={verifyOTP}
                className="mx-auto mt-9 flex flex-col items-center"
              >
                <InputOTP
                  className="flex justify-center"
                  maxLength={6}
                  pattern={REGEXP_ONLY_DIGITS_AND_CHARS}
                >
                  <InputOTPGroup className="flex justify-center">
                    <InputOTPSlot
                      className="h-16 w-16 text-3xl text-white shadow-none"
                      index={0}
                    />
                    <InputOTPSlot
                      className="h-16 w-16 text-3xl text-white shadow-none"
                      index={1}
                    />
                    <InputOTPSlot
                      className="h-16 w-16 text-3xl text-white shadow-none"
                      index={2}
                    />
                    <InputOTPSlot
                      className="h-16 w-16 text-3xl text-white shadow-none"
                      index={3}
                    />
                    <InputOTPSlot
                      className="h-16 w-16 text-3xl text-white shadow-none"
                      index={4}
                    />
                    <InputOTPSlot
                      className="h-16 w-16 text-3xl text-white shadow-none"
                      index={5}
                    />
                  </InputOTPGroup>
                </InputOTP>
                <p
                  id="helper-text-explanation"
                  className="mt-2 text-center text-sm font-light text-white"
                >
                  Did't get OTP Code?{"  "}
                  <button
                    onClick={sendOTP}
                    type="button"
                    className="mx-auto font-bold text-purple-400"
                  >
                    Send Again
                  </button>
                </p>
                <button
                  type="submit"
                  className="mt-5 w-full rounded-lg bg-gradient-to-r from-purple-500 via-purple-600 to-purple-700 px-5 py-2.5 text-center text-sm font-medium text-white shadow-lg shadow-purple-800/80 hover:bg-gradient-to-br focus:ring-4 focus:ring-purple-800 focus:outline-none"
                >
                  {loading ? (
                    <div role="status">
                      <svg
                        aria-hidden="true"
                        className="inline h-6 w-6 animate-spin fill-purple-600 text-gray-200"
                        viewBox="0 0 100 101"
                        fill="none"
                        xmlns="http://www.w3.org/2000/svg"
                      >
                        <path
                          d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z"
                          fill="currentColor"
                        />
                        <path
                          d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z"
                          fill="currentFill"
                        />
                      </svg>
                      <span className="sr-only">Loading...</span>
                    </div>
                  ) : (
                    "Verify"
                  )}
                </button>
                {/* <p className="font-light text-zinc-300 text-sm mt-16 text-center">
                  Want to Change Your Email?
                  <button
                    onClick={() => setSubmited(!submited)}
                    className="font-bold text-white"
                  >
                    {" "}
                    Back to Register
                  </button>
                </p> */}
              </form>
            </>
          ) : (
            <>
              <h1 className="bg-gradient-to-r from-white from-20% to-purple-500 bg-clip-text pb-2 text-center text-5xl font-bold text-transparent">
                Sign Up
              </h1>
              <form onSubmit={onSubmit} className="mt-4 space-y-4 py-4">
                <div className="mb-5">
                  <label
                    htmlFor="name"
                    className="mb-2 block w-fit bg-gradient-to-r from-white from-60% to-purple-400 bg-clip-text text-sm font-medium text-transparent"
                  >
                    Name
                  </label>
                  <input
                    type="text"
                    id="name"
                    className="bg-black-50 block w-full rounded-lg border border-purple-500 bg-gray-700 p-2.5 text-sm text-white placeholder-zinc-400 focus:border-purple-500 focus:ring-purple-800"
                    placeholder="Abigail Rachel"
                    {...register("name")}
                  />
                  {formState.errors.name && (
                    <p className="mt-2 text-sm text-red-700">
                      <span className="font-medium">Oops!</span>{" "}
                      {formState.errors.name?.message?.toString()}
                    </p>
                  )}
                </div>
                <div className="mb-5">
                  <label
                    htmlFor="email"
                    className="mb-2 block w-fit bg-gradient-to-r from-white from-60% to-purple-400 bg-clip-text text-sm font-medium text-transparent"
                  >
                    Email
                  </label>
                  <input
                    type="email"
                    id="email"
                    className="bg-black-50 block w-full rounded-lg border border-purple-500 bg-gray-700 p-2.5 text-sm text-white placeholder-zinc-400 focus:border-purple-500 focus:ring-purple-800"
                    placeholder="aralie@mail.com"
                    {...register("email")}
                  />
                  {formState.errors.email && (
                    <p className="mt-2 text-sm text-red-700">
                      <span className="font-medium">Oops!</span>{" "}
                      {formState.errors.email?.message?.toString()}
                    </p>
                  )}
                </div>
                <div className="mb-5">
                  <label
                    htmlFor="password"
                    className="mb-2 block w-fit bg-gradient-to-r from-white from-60% to-purple-400 bg-clip-text text-sm font-medium text-transparent"
                  >
                    Password
                  </label>
                  <input
                    type="password"
                    id="password"
                    className="bg-black-50 block w-full rounded-lg border border-purple-500 bg-gray-700 p-2.5 text-sm text-white placeholder-zinc-400 focus:border-purple-500 focus:ring-purple-800"
                    placeholder="********"
                    {...register("password")}
                  />
                  {formState.errors.password && (
                    <p className="mt-2 text-sm text-red-700">
                      <span className="font-medium">Oops!</span>{" "}
                      {formState.errors.password?.message?.toString()}
                    </p>
                  )}
                </div>
                <div className="mb-5">
                  <label
                    htmlFor="password-confirm"
                    className="mb-2 block w-fit bg-gradient-to-r from-white from-60% to-purple-400 bg-clip-text text-sm font-medium text-transparent"
                  >
                    Confirm Password
                  </label>
                  <input
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    type="password"
                    id="password-confirm"
                    className="bg-black-50 block w-full rounded-lg border border-purple-500 bg-gray-700 p-2.5 text-sm text-white placeholder-zinc-400 focus:border-purple-500 focus:ring-purple-800"
                    placeholder="********"
                  />
                  {error && (
                    <p className="mt-2 text-sm text-red-700">
                      <span className="font-medium">Oops! </span> {error}
                    </p>
                  )}
                </div>
                <div className="text-center text-sm font-extralight text-white">
                  Already have an account?{" "}
                  <Link to={"/login"} className="font-bold text-purple-950">
                    Login
                  </Link>{" "}
                  here.
                </div>
                <button
                  disabled={loading}
                  type="submit"
                  className="mt-5 mb-2 w-full rounded-lg bg-gradient-to-r from-purple-500 via-purple-600 to-purple-700 px-5 py-2.5 text-center text-sm font-medium text-white shadow-lg shadow-purple-800/80 hover:bg-gradient-to-br focus:ring-4 focus:ring-purple-800 focus:outline-none"
                >
                  {loading ? (
                    <div role="status">
                      <svg
                        aria-hidden="true"
                        className="inline h-6 w-6 animate-spin fill-purple-600 text-gray-200"
                        viewBox="0 0 100 101"
                        fill="none"
                        xmlns="http://www.w3.org/2000/svg"
                      >
                        <path
                          d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z"
                          fill="currentColor"
                        />
                        <path
                          d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z"
                          fill="currentFill"
                        />
                      </svg>
                      <span className="sr-only">Loading...</span>
                    </div>
                  ) : (
                    "Sign Up"
                  )}
                </button>
              </form>
            </>
          )}
        </div>
      </div>
      {errorOTP && (
        <div
          id="toast-danger"
          className="fixed right-5 bottom-5 z-10 mb-4 flex w-full max-w-xs items-center rounded-lg bg-white p-4 text-gray-500 shadow"
          role="alert"
        >
          <div className="inline-flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg bg-red-100 text-red-500">
            <svg
              className="h-5 w-5"
              aria-hidden="true"
              xmlns="http://www.w3.org/2000/svg"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path d="M10 .5a9.5 9.5 0 1 0 9.5 9.5A9.51 9.51 0 0 0 10 .5Zm3.707 11.793a1 1 0 1 1-1.414 1.414L10 11.414l-2.293 2.293a1 1 0 0 1-1.414-1.414L8.586 10 6.293 7.707a1 1 0 0 1 1.414-1.414L10 8.586l2.293-2.293a1 1 0 0 1 1.414 1.414L11.414 10l2.293 2.293Z" />
            </svg>
            <span className="sr-only">Error icon</span>
          </div>
          <div className="ms-3 text-sm font-normal">{errorOTP}</div>
          <button
            disabled={loading}
            onClick={() => setErrorOTP("")}
            type="button"
            className="-mx-1.5 -my-1.5 ms-auto inline-flex h-8 w-8 items-center justify-center rounded-lg bg-white p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-900 focus:ring-2 focus:ring-gray-300"
            data-dismiss-target="#toast-danger"
            aria-label="Close"
          >
            <span className="sr-only">Close</span>
            <svg
              className="h-3 w-3"
              aria-hidden="true"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 14 14"
            >
              <path
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"
              />
            </svg>
          </button>
        </div>
      )}
    </div>
  );
}
