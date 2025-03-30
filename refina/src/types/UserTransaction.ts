import { z } from "zod";

const Transaction = z.object({
  user_id: z.string(),
  name: z.string(),
  email: z.string(),
  wallets: z.array(
    z.object({
      id: z.string(),
      number: z.string(),
      balance: z.number(),
      name: z.string(),
      type: z.string(),
      transactions: z.array(
        z.object({
          id: z.string(),
          name: z.string(),
          type: z.string(),
          amount: z.number(),
          date: z.string(),
          description: z.string(),
          image: z.string().nullable(),
        }),
      ),
    }),
  ),
});

export type TransactionType = z.infer<typeof Transaction>;