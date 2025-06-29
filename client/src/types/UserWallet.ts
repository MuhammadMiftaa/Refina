import { z } from "zod";

const Wallet = z.object({
  id: z.string(),
  user_id: z.string(),
  wallet_number: z.string(),
  wallet_balance: z.number(),
  wallet_name: z.string(),
  wallet_type_name: z.string(),
  wallet_type: z.string(),
});

const WalletType = z.object({
  id: z.string(),
  name: z.string(),
  type: z.string(),
  description: z.string(),
});

const WalletsByType = z.object({
  user_id: z.string(),
  type: z.string(),
  wallets: z.array(
    z.object({
      id: z.string(),
      name: z.string(),
      number: z.string(),
      balance: z.number(),
    }),
  ),
});

export type WalletType = z.infer<typeof Wallet>;
export type WalletTypeType = z.infer<typeof WalletType>;
export type WalletsByTypeType = z.infer<typeof WalletsByType>;
