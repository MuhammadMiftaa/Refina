import styled from "styled-components";

export const SubmitButton = ({ text }: { text: string }) => {
  return (
    <StyledWrapper className="w-full md:w-70">
      <button type="submit" className="btn-donate">{text}</button>
    </StyledWrapper>
  );
};

const StyledWrapper = styled.div`
  .btn-donate {
    --clr-font-main: hsla(0 0% 20% / 100);
    --btn-bg-1: hsla(194 100% 69% / 1);
    --btn-bg-2: hsla(217 100% 56% / 1);
    --btn-bg-color: hsla(360 100% 100% / 1);
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
      0px 0px 20px rgba(71, 184, 255, 0.5),
      0px 5px 5px -1px rgba(58, 125, 233, 0.25),
      inset 4px 4px 8px rgba(175, 230, 255, 0.5),
      inset -4px -4px 8px rgba(19, 95, 216, 0.35);
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
