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
import { MdOutlineAlternateEmail } from "react-icons/md";
import { IoMailOutline } from "react-icons/io5";
import { SlLock } from "react-icons/sl";
import { LuEye, LuEyeOff } from "react-icons/lu";

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
  const [loading, setLoading] = useState<boolean>(false);
  const [showPassword, setShowPassword] = useState<boolean>(false);
  const [showConfirmPassword, setShowConfirmPassword] =
    useState<boolean>(false);

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
      navigate("/register/verification");
    } else {
      setError(res.message);
    }
  };

  return props.isAuthenticated ? (
    <Navigate to={"/"} />
  ) : (
    <div
      className="grid min-h-screen w-full place-items-center text-[#645e74]"
      style={{
        backgroundImage:
          "url('https://images.unsplash.com/photo-1567201864585-6baec9110dac?q=80&w=687&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D')",
        backgroundSize: "cover",
        backgroundPosition: "center",
        backgroundRepeat: "no-repeat",
      }}
    >
      <div className="flex w-[clamp(300px,90vw,800px)] flex-col-reverse rounded-[22px] bg-white p-5 shadow-[0_50px_100px_rgba(0,0,0,0.08)] sm:flex-row sm:p-2">
        <form
          onSubmit={onSubmit}
          className="flex w-full flex-col gap-3 px-5 py-7 sm:w-1/2 sm:px-12 sm:py-8"
        >
          <h2 className="text-center text-[24px] font-semibold tracking-[0.5px] text-[#8864f0] sm:text-left">
            Refina
          </h2>
          <h3 className="mb-3 text-center text-[14px] sm:text-left">
            Create your free account to get started
          </h3>

          <div className="flex h-11 items-center gap-4 rounded-md border border-[#d0d0d6] px-4 text-[16px] outline-[#8864f0]">
            <div>
              <MdOutlineAlternateEmail className="text-lg text-neutral-500" />
            </div>
            <input
              type="text"
              placeholder="Name"
              {...register("name")}
              className="focus:outline-none"
            />
          </div>
          {formState.errors.name && (
            <p className="text-sm text-red-500">
              <strong>Oops!</strong> {formState.errors.name.message}
            </p>
          )}

          <div className="flex h-11 items-center gap-4 rounded-md border border-[#d0d0d6] px-4 text-[16px] outline-[#8864f0]">
            <div>
              <IoMailOutline className="text-lg text-neutral-500" />
            </div>
            <input
              type="email"
              placeholder="Email"
              {...register("email")}
              className="focus:outline-none"
            />
          </div>
          {formState.errors.email && (
            <p className="text-sm text-red-500">
              <strong>Oops!</strong> {formState.errors.email.message}
            </p>
          )}

          <div className="relative flex h-11 items-center gap-4 rounded-md border border-[#d0d0d6] px-4 text-[16px] outline-[#8864f0]">
            <div>
              <SlLock className="text-lg text-neutral-500" />
            </div>
            <input
              type={showPassword ? "text" : "password"}
              placeholder="Password"
              {...register("password")}
              className="focus:outline-none"
            />
            <button
              className="absolute right-5 flex cursor-pointer"
              onClick={() => setShowPassword(!showPassword)}
              type="button"
            >
              {showPassword ? (
                <LuEye className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2" />
              ) : (
                <LuEyeOff className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2" />
              )}
            </button>
          </div>
          {formState.errors.password && (
            <p className="text-sm text-red-500">
              <strong>Oops!</strong> {formState.errors.password.message}
            </p>
          )}

          <div className="flex h-11 items-center gap-4 rounded-md border border-[#d0d0d6] px-4 text-[16px] outline-[#8864f0] relative">
            <div>
              <SlLock className="text-lg text-neutral-500" />
            </div>
            <input
              type={showConfirmPassword ? "text" : "password"}
              placeholder="Confirm Password"
              onChange={(e) => setConfirmPassword(e.target.value)}
              className="focus:outline-none"
            />
            <button
              className="absolute right-5 flex cursor-pointer"
              onClick={() => setShowConfirmPassword(!showConfirmPassword)}
              type="button"
            >
              {showConfirmPassword ? (
                <LuEye className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2" />
              ) : (
                <LuEyeOff className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2" />
              )}
            </button>
          </div>
          {error && (
            <p className="text-sm text-red-500">
              <strong>Oops!</strong> {error}
            </p>
          )}

          <button
            type="submit"
            disabled={loading}
            className="h-11 cursor-pointer rounded-md bg-[#8864f0] text-[17px] text-white transition-colors duration-200 hover:bg-[#7a5dcf] disabled:cursor-not-allowed disabled:bg-[#8864f0]/50 disabled:text-gray-300"
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
              "Register"
            )}
          </button>

          <p className="mt-4 text-center text-sm">
            Already have an account?{" "}
            <Link to="/login" className="font-bold text-purple-800">
              Login
            </Link>{" "}
            here.
          </p>
        </form>

        <div className="w-full overflow-hidden rounded-xl bg-[url('/hero.svg')] bg-cover bg-no-repeat font-[Space_Grotesk] sm:w-1/2">
          <div className="rounded-inherit flex h-full min-h-[170px] flex-col justify-center bg-gradient-to-b from-[rgba(95,69,168,0)] to-[rgba(95,69,168,0.7)] px-10 py-10 sm:items-start sm:justify-center sm:px-9">
            <h2 className="text-[22px] leading-tight font-medium text-white">
              Join Refina and take control of your financial journey
            </h2>
            <h3 className="mt-3 hidden text-[18px] text-[#c7c2d6] sm:block">
              Set goals, manage your budget, and make smarter financial
              decisions.
            </h3>
          </div>
        </div>
      </div>
    </div>
  );
}

export function RegisterOTP(props: { isAuthenticated: boolean }) {
  const navigate = useNavigate();
  const backendURL = getBackendURL();

  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState("");
  const [errorOTP, setErrorOTP] = useState("");

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
      className="grid min-h-screen w-full place-items-center text-[#645e74]"
      style={{
        backgroundImage:
          "url('https://images.unsplash.com/photo-1567201864585-6baec9110dac?q=80&w=687&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D')",
        backgroundSize: "cover",
        backgroundPosition: "center",
        backgroundRepeat: "no-repeat",
      }}
    >
      <div className="flex w-[clamp(300px,90vw,800px)] flex-col-reverse rounded-[22px] bg-white p-5 shadow-[0_50px_100px_rgba(0,0,0,0.08)] sm:flex-row sm:p-2">
        <form
          onSubmit={verifyOTP}
          className="flex w-full flex-col gap-3 px-5 py-7 sm:w-1/2 sm:px-12 sm:py-8"
        >
          <h2 className="text-center text-[24px] font-semibold tracking-[0.5px] text-[#8864f0] sm:text-left">
            Refina
          </h2>
          <h3 className="mb-3 text-center text-[14px] sm:text-left">
            Create your free account to get started
          </h3>

          <div className="mx-auto my-8 flex flex-col items-center">
            <InputOTP
              className="flex justify-center"
              maxLength={6}
              pattern={REGEXP_ONLY_DIGITS_AND_CHARS}
            >
              <InputOTPGroup className="flex justify-center">
                <InputOTPSlot
                  className="h-12 w-12 text-3xl shadow-none"
                  index={0}
                />
                <InputOTPSlot
                  className="h-12 w-12 text-3xl shadow-none"
                  index={1}
                />
                <InputOTPSlot
                  className="h-12 w-12 text-3xl shadow-none"
                  index={2}
                />
                <InputOTPSlot
                  className="h-12 w-12 text-3xl shadow-none"
                  index={3}
                />
                <InputOTPSlot
                  className="h-12 w-12 text-3xl shadow-none"
                  index={4}
                />
                <InputOTPSlot
                  className="h-12 w-12 text-3xl shadow-none"
                  index={5}
                />
              </InputOTPGroup>
            </InputOTP>
            <p
              id="helper-text-explanation"
              className="mt-2 text-center text-sm font-light text-neutral-500"
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
          </div>

          {errorOTP && (
            <p className="text-sm text-red-500">
              <strong>Oops!</strong> {errorOTP}
            </p>
          )}
          {error && (
            <p className="text-sm text-red-500">
              <strong>Oops!</strong> {error}
            </p>
          )}

          <button
            type="submit"
            disabled={loading}
            className="h-11 cursor-pointer rounded-md bg-[#8864f0] text-[17px] text-white transition-colors duration-200 hover:bg-[#7a5dcf] disabled:cursor-not-allowed disabled:bg-[#8864f0]/50 disabled:text-gray-300"
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
        </form>

        <div className="w-full overflow-hidden rounded-xl bg-[url('/hero.svg')] bg-cover bg-no-repeat font-[Space_Grotesk] sm:w-1/2">
          <div className="rounded-inherit flex h-full min-h-[170px] flex-col justify-center bg-gradient-to-b from-[rgba(95,69,168,0)] to-[rgba(95,69,168,0.7)] px-10 py-10 sm:items-start sm:justify-center sm:px-9">
            <h2 className="text-[22px] leading-tight font-medium text-white">
              Join Refina and take control of your financial journey
            </h2>
            <h3 className="mt-3 hidden text-[18px] text-[#c7c2d6] sm:block">
              Set goals, manage your budget, and make smarter financial
              decisions.
            </h3>
          </div>
        </div>
      </div>
    </div>
  );
}
