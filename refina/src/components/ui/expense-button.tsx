import styled from "styled-components";
import { GiPayMoney } from "react-icons/gi";

const IncomeButton = ({ onclick }: { onclick: () => void }) => {
  return (
    <StyledWrapper onClick={onclick}>
      <button className="Btn">
        Expense
        <GiPayMoney className="icon" />
      </button>
    </StyledWrapper>
  );
};

const StyledWrapper = styled.div`
  .Btn {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: flex-start;
    width: 130px;
    height: 40px;
    border: none;
    padding: 0px 20px;
    background-color: #EF4444;
    color: white;
    font-weight: 500;
    cursor: pointer;
    border-radius: 10px;
    box-shadow: 5px 5px 0px #DC2626;
    transition-duration: 0.3s;
  }

  .icon {
    width: 13px;
    position: absolute;
    right: 0;
    margin-right: 20px;
    fill: white;
    transition-duration: 0.3s;
  }

  .Btn:hover {
    color: transparent;
  }

  .Btn:hover .icon {
    right: 43%;
    margin: 0;
    padding: 0;
    border: none;
    transition-duration: 0.3s;
  }

  .Btn:active {
    transform: translate(3px, 3px);
    transition-duration: 0.3s;
    box-shadow: 2px 2px 0px rgb(140, 32, 212);
  }
`;

export default IncomeButton;
