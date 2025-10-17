import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useState, useCallback } from "react";
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

export default function Login({
  handleLogin,
  isAuthenticated,
}: {
  handleLogin: () => void;
  isAuthenticated: boolean;
}) {
  const backendURL = getBackendURL();
  const { setProfile } = useProfile(
    useShallow((state) => ({ setProfile: state.setProfile })),
  );
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const { register, handleSubmit, formState } = useForm<PostFormSchema>({
    resolver: zodResolver(postFormSchema),
  });

  // 🚀 gunakan useCallback untuk mencegah re-render berulang
  const onSubmit = useCallback(
    handleSubmit(async (data) => {
      try {
        setLoading(true);
        const res = await fetch(`${backendURL}/auth/login`, {
          method: "POST",
          credentials: "include",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(data),
        }).then((r) => r.json());

        if (res.status) {
          Cookies.set("token", res.data, createCookiesOpts());
          handleLogin();
          navigate("/");
        } else {
          setError(res.message);
        }
      } catch (err) {
        setError("Failed to connect to server.");
      } finally {
        setLoading(false);
      }
    }),
    [backendURL, handleLogin, navigate],
  );

  const handleOAuth = useCallback(
    async (server: string) => {
      try {
        const res = await fetch(`${backendURL}/auth/${server}/oauth`, {
          credentials: "include",
        });
        const data = await res.json();
        if (res.ok) window.location.href = data.url;
        else console.error(data.message || "Login failed");
      } catch (err) {
        console.error("Error during OAuth:", err);
      }
    },
    [backendURL],
  );

  // 🧠 Gunakan efek hanya 1x
  useEffect(() => {
    const token = searchParams.get("token");
    if (token) {
      Cookies.set("token", token, createCookiesOpts());
      const decoded = decodeJwt(token);
      setProfile({
        username: decoded.username as string,
        email: decoded.email as string,
      });
      handleLogin();
      navigate("/");
    }
  }, [searchParams, setProfile, handleLogin, navigate]);

  // 🧩 Lazy-load background agar tidak menghambat render awal
  const [bgLoaded, setBgLoaded] = useState(false);
  useEffect(() => {
    const img = new Image();
    img.src = "/background.jpeg";
    img.onload = () => setBgLoaded(true);
  }, []);

  return isAuthenticated ? (
    <Navigate to="/" />
  ) : (
    <div
      className={`grid min-h-screen w-full place-items-center text-[#645e74] transition-opacity duration-300 ${
        bgLoaded ? "opacity-100" : "opacity-0"
      }`}
      style={
        bgLoaded
          ? {
              backgroundImage: "url('/background.jpeg')",
              backgroundSize: "cover",
              backgroundPosition: "center",
              backgroundRepeat: "no-repeat",
            }
          : {}
      }
    >
      <div className="flex w-[clamp(300px,90vw,800px)] flex-col rounded-[22px] bg-white p-5 shadow-[0_50px_100px_rgba(0,0,0,0.08)] sm:flex-row sm:p-2">
        {/* 🖼️ gunakan <img> dengan lazy loading agar hero.svg tidak block render */}
        <div className="w-full overflow-hidden rounded-xl sm:w-1/2">
          <img
            src="/hero.svg"
            alt="Refina Hero"
            width={400}
            height={400}
            loading="lazy"
            className="object-cover h-full w-full"
          />
          <div className="absolute bottom-0 left-0 right-0 flex flex-col justify-end bg-gradient-to-b from-transparent to-[rgba(95,69,168,0.7)] px-6 py-6">
            <h2 className="text-[22px] font-medium text-white">
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
              {/* ✅ gunakan lazy load icon */}
              <img
                src="/google.svg"
                alt="Google"
                className="h-5"
                width={20}
                height={20}
                loading="lazy"
              />
              <p className="text-[#7e7c83]">
                <span className="hidden sm:inline">Login with</span> Google
              </p>
            </button>

            <button
              type="button"
              className="flex h-11 w-full items-center justify-center gap-2 rounded-md border border-transparent bg-[#f2f3f6] text-[15px]"
              onClick={() => handleOAuth("facebook")}
            >
              <img
                src="/facebook.svg"
                alt="Facebook"
                className="h-6"
                width={20}
                height={20}
                loading="lazy"
              />
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

          {/* Input email */}
          <div className="flex h-11 items-center gap-4 rounded-md border border-[#d0d0d6] px-4 text-[16px] outline-[#8864f0]">
            <IoMailOutline className="text-lg text-neutral-500" />
            <input
              type="email"
              placeholder="Email"
              {...register("email")}
              className="flex-1 focus:outline-none"
            />
          </div>
          {formState.errors.email && (
            <p className="text-sm text-red-500">
              <strong>Oops!</strong> {formState.errors.email.message}
            </p>
          )}

          {/* Input password */}
          <div className="relative flex h-11 items-center gap-4 rounded-md border border-[#d0d0d6] px-4 text-[16px] outline-[#8864f0]">
            <SlLock className="text-lg text-neutral-500" />
            <input
              type={showPassword ? "text" : "password"}
              placeholder="Password"
              {...register("password")}
              className="flex-1 focus:outline-none"
            />
            <button
              type="button"
              onClick={() => setShowPassword((s) => !s)}
              className="absolute right-4"
            >
              {showPassword ? <LuEye /> : <LuEyeOff />}
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
            {loading ? "Loading..." : "Login"}
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
