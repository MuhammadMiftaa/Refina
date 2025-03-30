import { z } from "zod";

const Wallet = z.object({
  user_id: z.string(),
  name: z.string(),
  email: z.string(),
  wallets: z.array(
    z.object({
      id: z.string(),
      number: z.string(),
      balance: z.number(),
      name: z.string(),
      wallet_type: z.string(),
      type: z.string(),
    }),
  ),
});

const WalletType = z.object({
  id: z.string(),
  name: z.string(),
  type: z.string(),
  description: z.string(),
})

export type WalletType = z.infer<typeof Wallet>;
export type WalletTypeType = z.infer<typeof WalletType>;