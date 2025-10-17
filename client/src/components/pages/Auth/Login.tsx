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

export default function LoginOptimized(props: {
  handleLogin: () => void;
  isAuthenticated: boolean;
}) {
  const backendURL = getBackendURL();
  const { setProfile } = useProfile(
    useShallow((state) => ({ setProfile: state.setProfile })),
  );
  const navigate = useNavigate();
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const { register, handleSubmit, formState } = useForm<PostFormSchema>({
    resolver: zodResolver(postFormSchema),
  });

  // 🧠 Optimized Fetch with try/catch
  const onSubmit = handleSubmit(async (data) => {
    try {
      setLoading(true);
      const res = await fetch(`${backendURL}/auth/login`, {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      });
      const result = await res.json();

      if (result.status) {
        Cookies.set("token", result.data, createCookiesOpts());
        props.handleLogin();
        navigate("/");
      } else {
        setError(result.message || "Login failed");
      }
    } catch (err) {
      console.error("Error during login:", err);
      setError("Network error. Please try again.");
    } finally {
      setLoading(false);
    }
  });

  // OAuth handler
  const handleOAuth = async (server: string) => {
    try {
      const res = await fetch(`${backendURL}/auth/${server}/oauth`, {
        method: "GET",
        credentials: "include",
      });
      const data = await res.json();
      if (res.ok) window.location.href = data.url;
      else console.error(data.message || "Login failed");
    } catch (err) {
      console.error("Error during OAuth:", err);
    }
  };

  // Handle token from redirect
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
  }, []);

  if (props.isAuthenticated) return <Navigate to={"/"} />;

  return (
    <div
      className="grid min-h-screen w-full place-items-center text-[#645e74]"
      style={{
        backgroundImage: "url('/background.jpeg')",
        backgroundSize: "cover",
        backgroundPosition: "center",
        backgroundRepeat: "no-repeat",
      }}
    >
      {/* 🔥 Preload background untuk LCP */}
      <link rel="preload" as="image" href="/background.jpeg" />

      <div className="flex w-[clamp(300px,90vw,800px)] flex-col rounded-[22px] bg-white p-5 shadow-[0_50px_100px_rgba(0,0,0,0.08)] sm:flex-row sm:p-2">
        {/* 🖼️ Optimized Hero Section */}
        <div className="relative w-full overflow-hidden rounded-xl sm:w-1/2 font-[Space_Grotesk]">
          <img
            src="/hero.svg"
            alt="Hero Illustration"
            loading="lazy"
            decoding="async"
            className="absolute inset-0 h-full w-full object-cover"
          />
          <div className="relative flex h-full min-h-[170px] flex-col justify-center bg-gradient-to-b from-[rgba(95,69,168,0)] to-[rgba(95,69,168,0.7)] px-10 py-10 sm:items-start sm:justify-center sm:px-9">
            <h2 className="text-[22px] leading-tight font-medium text-white">
              Manage your finances easily and securely with Refina
            </h2>
            <h3 className="mt-3 hidden text-[18px] text-[#c7c2d6] sm:block">
              Track spending, manage budgets, and set your financial goals.
            </h3>
          </div>
        </div>

        {/* 🧾 Form Section */}
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

          {/* OAuth Buttons */}
          <div className="flex gap-2 sm:flex-col">
            <button
              type="button"
              aria-label="Login with Google"
              className="flex h-11 w-full items-center justify-center gap-2 rounded-md border border-transparent bg-[#f2f3f6] text-[15px]"
              onClick={() => handleOAuth("google")}
            >
              <img src="/google.svg" alt="Google" className="h-5" loading="lazy" />
              <p className="text-[#7e7c83]">
                <span className="hidden sm:inline">Login with</span> Google
              </p>
            </button>
            <button
              type="button"
              aria-label="Login with Facebook"
              className="flex h-11 w-full items-center justify-center gap-2 rounded-md border border-transparent bg-[#f2f3f6] text-[15px]"
              onClick={() => handleOAuth("facebook")}
            >
              <img src="/facebook.svg" alt="Facebook" className="h-6" loading="lazy" />
              <p className="text-[#7e7c83]">
                <span className="hidden sm:inline">Login with</span> Facebook
              </p>
            </button>
          </div>

          {/* Divider */}
          <span className="relative h-6 text-center">
            <span className="absolute top-1/2 left-0 h-px w-full -translate-y-1/2 bg-[#d0d0d6] opacity-60"></span>
            <span className="absolute top-1/2 left-1/2 z-10 -translate-x-1/2 -translate-y-1/2 bg-white px-3 text-xs">
              Or
            </span>
          </span>

          {/* Email Input */}
          <div className="flex h-11 items-center gap-4 rounded-md border border-[#d0d0d6] px-4 text-[16px] outline-[#8864f0]">
            <IoMailOutline className="text-lg text-neutral-500" />
            <input
              type="email"
              placeholder="Email"
              {...register("email")}
              aria-label="Email"
              className="w-full focus:outline-none"
            />
          </div>
          {formState.errors.email && (
            <p className="text-sm text-red-500">
              <strong>Oops!</strong> {formState.errors.email.message}
            </p>
          )}

          {/* Password Input */}
          <div className="relative flex h-11 items-center gap-4 rounded-md border border-[#d0d0d6] px-4 text-[16px] outline-[#8864f0]">
            <SlLock className="text-lg text-neutral-500" />
            <input
              type={showPassword ? "text" : "password"}
              placeholder="Password"
              {...register("password")}
              aria-label="Password"
              className="w-full focus:outline-none"
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-5 cursor-pointer"
              aria-label={showPassword ? "Hide password" : "Show password"}
            >
              {showPassword ? <LuEye /> : <LuEyeOff />}
            </button>
          </div>
          {formState.errors.password && (
            <p className="text-sm text-red-500">
              <strong>Oops!</strong> {formState.errors.password.message}
            </p>
          )}

          {/* Error Message */}
          {error && (
            <p className="text-sm text-red-500">
              <strong>Oops!</strong> {error}
            </p>
          )}

          {/* Submit Button */}
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
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908Z"
                    fill="currentColor"
                  />
                  <path
                    d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539..."
                    fill="currentFill"
                  />
                </svg>
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
