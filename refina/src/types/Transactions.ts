export type transactionsType = {
  id: string;
  amount: number;
  transaction_type: "expense" | "income";
  date: string;
  description: string;
  category: string;
  user_id: string;
};
