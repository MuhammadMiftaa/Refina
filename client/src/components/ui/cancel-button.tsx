import styled from "styled-components";

export const CancelButton = ({
  text,
  loading,
  onclick,
}: {
  text: string;
  loading?: boolean;
  onclick?: () => void;
}) => {
  return (
    <StyledWrapper onClick={onclick} className="w-full">
      <button
        type="button"
        className="btn-donate disabled:cursor-not-allowed disabled:grayscale-50"
        disabled={loading}
      >
        {loading ? (
          <div role="status">
            <svg
              aria-hidden="true"
              className="inline h-6 w-6 animate-spin fill-gray-800 text-gray-200"
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
          text
        )}
      </button>
    </StyledWrapper>
  );
};

const StyledWrapper = styled.div`
  .btn-donate {
    --clr-font-main: hsla(0, 0%, 20%, 1);
    --btn-bg-1: hsla(0, 0%, 70%, 1);
    --btn-bg-2: hsla(0, 0%, 50%, 1);
    --btn-bg-color: hsla(0, 0%, 100%, 1);
    --radii: 0.5em;
    cursor: pointer;
    padding: 0.9em 1em;
    width: 100%;
    min-height: 44px;
    font-size: var(--size, 1rem);
    font-weight: 500;
    transition: 0.8s;
    background-size: 280% auto;
    background-image: linear-gradient(
      325deg,
      var(--btn-bg-2) 0%,
      var(--btn-bg-1) 55%,
      var(--btn-bg-2) 90%
    );
    border: none;
    border-radius: var(--radii);
    color: var(--btn-bg-color);
    box-shadow:
      0px 0px 10px rgba(128, 128, 128, 0.5),
      0px 3px 3px -1px rgba(90, 90, 90, 0.25),
      inset 4px 4px 8px rgba(180, 180, 180, 0.5),
      inset -4px -4px 8px rgba(90, 90, 90, 0.35);
  }

  .btn-donate:hover {
    background-position: right top;
  }

  .btn-donate:is(:focus, :focus-visible, :active) {
    outline: none;
    box-shadow:
      0 0 0 3px var(--btn-bg-color),
      0 0 0 6px var(--btn-bg-2);
  }

  @media (prefers-reduced-motion: reduce) {
    .btn-donate {
      transition: linear;
    }
  }
`;
