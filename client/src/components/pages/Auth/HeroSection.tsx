// HeroSection.tsx
export default function HeroSection() {
  return (
    <div className="relative w-full overflow-hidden rounded-xl font-[Space_Grotesk] sm:w-1/2">
      <img
        src="https://res.cloudinary.com/dblibr1t2/image/upload/f_avif,q_1,w_200/hero_uc1obl.svg"
        alt="Hero Illustration"
        fetchPriority="high"
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
  );
}
