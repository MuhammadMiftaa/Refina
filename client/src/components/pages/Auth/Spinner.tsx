// Spinner.tsx
export default function Spinner() {
  return (
    <div role="status" className="flex justify-center">
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
  );
}
