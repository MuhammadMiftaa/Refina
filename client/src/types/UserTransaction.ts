import { z } from "zod";

const Transaction = z.object({
  id: z.string(),
  user_id: z.string(),
  wallet_id: z.string(),
  wallet_number: z.string(),
  wallet_type: z.string(),
  wallet_type_name: z.string(),
  wallet_balance: z.number(),
  category_name: z.string(),
  category_type: z.string(),
  amount: z.number(),
  transaction_date: z.string(),
  description: z.string(),
  image: z.string().nullable(),
});

export type TransactionType = z.infer<typeof Transaction>;
