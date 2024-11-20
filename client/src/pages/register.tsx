import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { Link, Navigate, useNavigate } from "react-router-dom";
import { z } from "zod";

const postFormSchema = z.object({
  name: z.string(),
  email: z.string().email(),
  password: z.string().min(8, "Password must be at least 8 characters"),
});

type PostFormSchema = z.infer<typeof postFormSchema>;

export default function SignUpPage(props: {
  isAuthenticated: boolean;
  handleLogin: () => void;
}) {
  const navigate = useNavigate();

  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");

  const { register, handleSubmit, formState, watch } = useForm<PostFormSchema>({
    resolver: zodResolver(postFormSchema),
  });

  const onSubmit = handleSubmit(async (data) => {
    if (watch("password") !== confirmPassword) {
      setError("Password and confirm password must be the same");
      return;
    }
    const res = await fetch("http://localhost:8080/v1/auth/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    }).then((res) => res.json());

    if (res.status) {
      props.handleLogin();
      navigate("/login");
    } else {
      setError(res.message);
    }
  });

  return props.isAuthenticated ? (
    <Navigate to={"/"} />
  ) : (
    <div
      className="min-h-screen flex after:content-[''] after:block after:absolute after:inset-0 after:bg-black after:bg-opacity-40"
      style={{
        backgroundImage: "url('/background.jpg')",
        backgroundSize: "cover",
        backgroundPosition: "center",
        backgroundRepeat: "no-repeat",
      }}
    >
      <div className="w-[45%] my-20 min-h-96 rounded-xl backdrop-blur bg-white/30 mx-auto z-10 font-inter">
        <div className="p-8">
          <h1 className="text-5xl pb-2 font-bold text-center bg-clip-text text-transparent from-white to-purple-500 bg-gradient-to-r from-20%">
            Sign Up
          </h1>
          <form onSubmit={onSubmit} className="mt-4 space-y-4">
            <div className="mb-5">
              <label
                htmlFor="name"
                className="block mb-2 text-sm font-medium text-transparent bg-clip-text bg-gradient-to-r from-white to-purple-400 from-60% w-fit "
              >
                Name
              </label>
              <input
                type="text"
                id="name"
                className="bg-black-50 border text-white  placeholder-zinc-400 text-sm rounded-lg focus:ring-purple-800 focus:border-purple-500 block w-full p-2.5 bg-gray-700 border-purple-500"
                placeholder="Abigail Rachel"
                {...register("name")}
              />
              {formState.errors.name && (
                <p className="mt-2 text-sm text-red-700 ">
                  <span className="font-medium">Oops!</span>{" "}
                  {formState.errors.name?.message?.toString()}
                </p>
              )}
            </div>
            <div className="mb-5">
              <label
                htmlFor="email"
                className="block mb-2 text-sm font-medium text-transparent bg-clip-text bg-gradient-to-r from-white to-purple-400 from-60% w-fit "
              >
                Email
              </label>
              <input
                type="email"
                id="email"
                className="bg-black-50 border text-white  placeholder-zinc-400 text-sm rounded-lg focus:ring-purple-800 focus:border-purple-500 block w-full p-2.5 bg-gray-700 border-purple-500"
                placeholder="aralie@mail.com"
                {...register("email")}
              />
              {formState.errors.email && (
                <p className="mt-2 text-sm text-red-700 ">
                  <span className="font-medium">Oops!</span>{" "}
                  {formState.errors.email?.message?.toString()}
                </p>
              )}
            </div>
            <div className="mb-5">
              <label
                htmlFor="password"
                className="block mb-2 text-sm font-medium text-transparent bg-clip-text bg-gradient-to-r from-white to-purple-400 from-60% w-fit "
              >
                Password
              </label>
              <input
                type="password"
                id="password"
                className="bg-black-50 border text-white  placeholder-zinc-400 text-sm rounded-lg focus:ring-purple-800 focus:border-purple-500 block w-full p-2.5 bg-gray-700 border-purple-500"
                placeholder="********"
                {...register("password")}
              />
              {formState.errors.password && (
                <p className="mt-2 text-sm text-red-700 ">
                  <span className="font-medium">Oops!</span>{" "}
                  {formState.errors.password?.message?.toString()}
                </p>
              )}
            </div>
            <div className="mb-5">
              <label
                htmlFor="password-confirm"
                className="block mb-2 text-sm font-medium text-transparent bg-clip-text bg-gradient-to-r from-white to-purple-400 from-60% w-fit "
              >
                Confirm Password
              </label>
              <input
                onChange={(e) => setConfirmPassword(e.target.value)}
                type="password"
                id="password-confirm"
                className="bg-black-50 border text-white  placeholder-zinc-400 text-sm rounded-lg focus:ring-purple-800 focus:border-purple-500 block w-full p-2.5 bg-gray-700 border-purple-500"
                placeholder="********"
              />
              {error && (
                <p className="mt-2 text-sm text-red-700 ">
                  <span className="font-medium">Oops! </span> {error}
                </p>
              )}
            </div>
            <div className="text-center text-white text-sm font-extralight">
              Already have an account?{" "}
              <Link to={"/login"} className="font-bold text-purple-950">
                Login
              </Link>{" "}
              here.
            </div>
            <button
              type="submit"
              className="w-full mt-5 text-white bg-gradient-to-r from-purple-500 via-purple-600 to-purple-700 hover:bg-gradient-to-br focus:ring-4 focus:outline-none focus:ring-purple-800 shadow-lg shadow-purple-800/80 font-medium rounded-lg text-sm px-5 py-2.5 text-center mb-2"
            >
              Sign Up
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
