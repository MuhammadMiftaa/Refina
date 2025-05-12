import { z } from "zod";

const Category = z.object({
    group_name: z.string(),
    type: z.enum(["income", "expense"]),
    category: z.array(
        z.object({
            id: z.string(),
            name: z.string(),
        })
    )
})

export type CategoryType = z.infer<typeof Category>;