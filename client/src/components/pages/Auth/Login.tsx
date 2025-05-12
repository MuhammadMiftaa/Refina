import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { VscGithub } from "react-icons/vsc";
import { Link, Navigate, useNavigate, useSearchParams } from "react-router-dom";
import { z } from "zod";
import { getBackendURL } from "../../../lib/readenv";
import Cookies from "js-cookie";
import { decodeJwt } from "jose";
import { useProfile } from "@/store/useProfile";
import { useShallow } from "zustand/shallow";

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

  const [error, setError] = useState("");

  const { register, handleSubmit, formState } = useForm<PostFormSchema>({
    resolver: zodResolver(postFormSchema),
  });

  const onSubmit = handleSubmit(async (data) => {
    const res = await fetch(`${backendURL}/auth/login`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    }).then((res) => res.json());

    if (res.status) {
      Cookies.set("token", res.data);
      props.handleLogin();
      navigate("/");
    } else {
      setError(res.message);
    }
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
      Cookies.set("token", token);
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
      className="relative flex min-h-screen after:absolute after:inset-0 after:bg-black/40 after:content-['']"
      style={{
        backgroundImage: "url('/background.jpg')",
        backgroundSize: "cover",
        backgroundPosition: "center",
        backgroundRepeat: "no-repeat",
      }}
    >
      <div className="font-inter z-10 mx-auto my-20 min-h-96 w-[40%] rounded-xl bg-white/30 backdrop-blur">
        <div className="p-8">
          <h1 className="bg-gradient-to-r from-white from-20% to-purple-500 bg-clip-text pb-2 text-center text-5xl font-bold text-transparent">
            Sign In
          </h1>
          <form onSubmit={onSubmit} className="mt-4 space-y-4">
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
              {error && (
                <p className="mt-2 text-sm text-red-700">
                  <span className="font-medium">Oops! </span> {error}
                </p>
              )}
            </div>
            <div className="text-center text-sm font-extralight text-white">
              Don't have an account?{" "}
              <Link to={"/register"} className="font-bold text-purple-950">
                Register
              </Link>{" "}
              here.
            </div>
            <button
              type="submit"
              className="mt-5 mb-2 w-full rounded-lg bg-gradient-to-r from-purple-500 via-purple-600 to-purple-700 px-5 py-2.5 text-center text-sm font-medium text-white shadow-lg shadow-purple-800/80 hover:bg-gradient-to-br focus:ring-4 focus:ring-purple-800 focus:outline-none"
            >
              Sign In
            </button>
          </form>
          <p className="font-poppins my-5 text-center text-[0.8rem] text-white before:mr-5 before:tracking-[-0.15em] before:content-['—————'] after:ml-5 after:tracking-[-0.15em] after:content-['—————']">
            Or sign in with it.
          </p>
          <div className="flex justify-stretch gap-6">
            <button
              onClick={() => handleOAuth("google")}
              type="button"
              className="flex w-full cursor-pointer justify-center rounded-xl border border-zinc-200 bg-gradient-to-br from-purple-400 via-purple-200 to-purple-400 py-2 text-2xl active:translate-y-0.5 active:shadow-none"
              style={{ boxShadow: "0 3px 3px #c084fc " }}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                x="0px"
                y="0px"
                width="26"
                height="26"
                viewBox="0 0 48 48"
              >
                <path
                  fill="#FFC107"
                  d="M43.611,20.083H42V20H24v8h11.303c-1.649,4.657-6.08,8-11.303,8c-6.627,0-12-5.373-12-12c0-6.627,5.373-12,12-12c3.059,0,5.842,1.154,7.961,3.039l5.657-5.657C34.046,6.053,29.268,4,24,4C12.955,4,4,12.955,4,24c0,11.045,8.955,20,20,20c11.045,0,20-8.955,20-20C44,22.659,43.862,21.35,43.611,20.083z"
                ></path>
                <path
                  fill="#FF3D00"
                  d="M6.306,14.691l6.571,4.819C14.655,15.108,18.961,12,24,12c3.059,0,5.842,1.154,7.961,3.039l5.657-5.657C34.046,6.053,29.268,4,24,4C16.318,4,9.656,8.337,6.306,14.691z"
                ></path>
                <path
                  fill="#4CAF50"
                  d="M24,44c5.166,0,9.86-1.977,13.409-5.192l-6.19-5.238C29.211,35.091,26.715,36,24,36c-5.202,0-9.619-3.317-11.283-7.946l-6.522,5.025C9.505,39.556,16.227,44,24,44z"
                ></path>
                <path
                  fill="#1976D2"
                  d="M43.611,20.083H42V20H24v8h11.303c-0.792,2.237-2.231,4.166-4.087,5.571c0.001-0.001,0.002-0.001,0.003-0.002l6.19,5.238C36.971,39.205,44,34,44,24C44,22.659,43.862,21.35,43.611,20.083z"
                ></path>
              </svg>
            </button>

            <button
              onClick={() => handleOAuth("github")}
              className="flex w-full cursor-pointer justify-center rounded-xl border border-purple-200 bg-gradient-to-br from-purple-400 via-purple-200 to-purple-400 py-2 text-2xl active:translate-y-0.5 active:shadow-none"
              style={{ boxShadow: "0 3px 3px #c084fc " }}
            >
              <VscGithub />
            </button>
            <button
              onClick={() => handleOAuth("microsoft")}
              className="flex w-full cursor-pointer justify-center rounded-xl border border-zinc-200 bg-gradient-to-br from-purple-400 via-purple-200 to-purple-400 py-2 text-2xl active:translate-y-0.5 active:shadow-none"
              style={{ boxShadow: "0 3px 3px #c084fc " }}
            >
              <svg
                x="0px"
                y="0px"
                width="23"
                height="23"
                // viewBox="0 0 48 48"
                enableBackground="new 0 0 2499.6 2500"
                viewBox="0 0 2499.6 2500"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="m1187.9 1187.9h-1187.9v-1187.9h1187.9z"
                  fill="#f1511b"
                />
                <path
                  d="m2499.6 1187.9h-1188v-1187.9h1187.9v1187.9z"
                  fill="#80cc28"
                />
                <path d="m1187.9 2500h-1187.9v-1187.9h1187.9z" fill="#00adef" />
                <path
                  d="m2499.6 2500h-1188v-1187.9h1187.9v1187.9z"
                  fill="#fbbc09"
                />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
