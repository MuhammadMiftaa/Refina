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
})

export type WalletType = z.infer<typeof Wallet>;
export type WalletTypeType = z.infer<typeof WalletType>;