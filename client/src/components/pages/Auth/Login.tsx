import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Link, Navigate, useNavigate, useSearchParams } from "react-router-dom";
import { z } from "zod";
import { getBackendURL } from "../../../lib/readenv";
import Cookies from "js-cookie";
import { decodeJwt } from "jose";
import { useProfile } from "@/store/useProfile";
import { useShallow } from "zustand/shallow";
import { SlLock } from "react-icons/sl";
import { LuEye, LuEyeOff } from "react-icons/lu";
import { IoMailOutline } from "react-icons/io5";
import { createCookiesOpts } from "@/helper/Helper";

const postFormSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8, "Password must be at least 8 characters"),
});

type PostFormSchema = z.infer<typeof postFormSchema>;

export default function Login(props: {
  handleLogin: () => void;
  isAuthenticated: boolean;
}) {
  const backendURL = getBackendURL();

  const { setProfile } = useProfile(
    useShallow((state) => ({ setProfile: state.setProfile })),
  );

  const navigate = useNavigate();

  const [showPassword, setShowPassword] = useState<boolean>(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const { register, handleSubmit, formState } = useForm<PostFormSchema>({
    resolver: zodResolver(postFormSchema),
  });

  const onSubmit = handleSubmit(async (data) => {
    setLoading(true);

    const res = await fetch(`${backendURL}/auth/login`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    }).then((res) => res.json());

    if (res.status) {
      Cookies.set("token", res.data, createCookiesOpts());
      props.handleLogin();
      navigate("/");
    } else {
      setError(res.message);
    }

    setLoading(false);
  });

  const handleOAuth = async (server: string) => {
    try {
      const res = await fetch(`${backendURL}/auth/${server}/oauth`, {
        method: "GET",
        credentials: "include",
      });
      const data = await res.json();

      if (res.ok) {
        window.location.href = data.url;
      } else {
        console.error(data.message || "Login failed");
      }
    } catch (err) {
      console.error("Error during OAuth:", err);
    }
  };

  const [searchParams] = useSearchParams();
  useEffect(() => {
    const token = searchParams.get("token");
    if (token) {
      Cookies.set("token", token, createCookiesOpts());
      const decoded = decodeJwt(token);
      setProfile({
        username: decoded.username as string,
        email: decoded.email as string,
      });
      props.handleLogin();
      navigate("/");
    }
  });

  return props.isAuthenticated ? (
    <Navigate to={"/"} />
  ) : (
    <div
      className="grid min-h-screen w-full place-items-center bg-[#e8dfff] text-[#645e74]"
      style={{
        backgroundImage: "url('/background.jpeg')",
        backgroundSize: "cover",
        backgroundPosition: "center",
        backgroundRepeat: "no-repeat",
      }}
    >
      <div className="flex w-[clamp(300px,90vw,800px)] flex-col rounded-[22px] bg-white p-5 shadow-[0_50px_100px_rgba(0,0,0,0.08)] sm:flex-row sm:p-2">
        <div className="w-full overflow-hidden rounded-xl bg-[url('/hero.svg')] bg-cover bg-no-repeat font-[Space_Grotesk] sm:w-1/2">
          <div className="rounded-inherit flex h-full min-h-[170px] flex-col justify-center bg-gradient-to-b from-[rgba(95,69,168,0)] to-[rgba(95,69,168,0.7)] px-10 py-10 sm:items-start sm:justify-center sm:px-9">
            <h2 className="text-[22px] leading-tight font-medium text-white">
              Manage your finances easily and securely with Refina
            </h2>
            <h3 className="mt-3 hidden text-[18px] text-[#c7c2d6] sm:block">
              Track spending, manage budgets, and set your financial goals.
            </h3>
          </div>
        </div>

        <form
          onSubmit={onSubmit}
          className="flex w-full flex-col gap-3 px-5 py-7 sm:w-1/2 sm:px-12 sm:py-8"
        >
          <h2 className="text-center text-[24px] font-semibold tracking-[0.5px] text-[#8864f0] sm:text-left">
            Refina
          </h2>
          <h3 className="mb-3 text-center text-[14px] sm:text-left">
            Login to your account
          </h3>

          <div className="flex gap-2 sm:flex-col">
            <button
              type="button"
              className="flex h-11 w-full items-center justify-center gap-2 rounded-md border border-transparent bg-[#f2f3f6] text-[15px]"
              onClick={() => handleOAuth("google")}
            >
              <img src="/google.svg" alt="Google" className="h-5" />
              <p className="text-[#7e7c83]">
                <span className="hidden sm:inline">Login with</span> Google
              </p>
            </button>
            <button
              type="button"
              className="flex h-11 w-full items-center justify-center gap-2 rounded-md border border-transparent bg-[#f2f3f6] text-[15px]"
              onClick={() => handleOAuth("facebook")}
            >
              <img src="/facebook.svg" alt="Facebook" className="h-6" />
              <p className="text-[#7e7c83]">
                <span className="hidden sm:inline">Login with</span> Facebook
              </p>
            </button>
          </div>

          <span className="relative h-6 text-center">
            <span className="absolute top-1/2 left-0 h-px w-full -translate-y-1/2 bg-[#d0d0d6] opacity-60"></span>
            <span className="absolute top-1/2 left-1/2 z-10 -translate-x-1/2 -translate-y-1/2 bg-white px-3 text-xs">
              Or
            </span>
          </span>

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
              "Login"
            )}
          </button>

          <p className="mt-4 text-center text-sm">
            Don't have an account?{" "}
            <Link to="/register" className="font-bold text-purple-800">
              Register
            </Link>{" "}
            here.
          </p>
        </form>
      </div>
    </div>
  );
}
