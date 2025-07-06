import { z } from "zod";

const Attachment = z.object({
    id: z.string(),
    transaction_id: z.string(),
    image: z.string(),
    format: z.string(),
    size: z.number(),
})

export type AttachmentType = z.infer<typeof Attachment>;