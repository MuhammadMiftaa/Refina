export type UserTransactionsType = {
  id: string;
  user_name: string;
  wallet_name: string;
  wallet_type: string;
  category_name: string;
  category_type: "expense" | "income";
  amount: number;
  date: string;
  description: string;
  category: string;
  user_id: string;
};
