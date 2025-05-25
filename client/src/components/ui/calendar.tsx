import * as React from "react"
import { ChevronLeft, ChevronRight } from "lucide-react"
import { DayPicker } from "react-day-picker"

import { cn } from "@/lib/utils"
import { buttonVariants } from "@/components/ui/button"

function Calendar({
  className,
  classNames,
  showOutsideDays = true,
  ...props
}: React.ComponentProps<typeof DayPicker>) {
  return (
    <DayPicker
      showOutsideDays={showOutsideDays}
      className={cn("p-3", className)}
      classNames={{
        months: "",
        month: "",
        caption: "",
        caption_label: "absolute top-3 left-1/2 -translate-x-1/2",
        nav: "flex items-center justify-between",
        nav_button: cn(
          buttonVariants({ variant: "outline" }),
          ""
        ),
        nav_button_previous: "",
        nav_button_next: "",
        table: "",
        head_row: "",
        head_cell:
          "",
        row: "",
        cell: cn(
          "",
        ),
        day: cn(
          buttonVariants({ variant: "ghost" }),
          "size-8 p-0 font-normal aria-selected:opacity-100"
        ),
        day_range_start:
          "",
        day_range_end:
          "",
        day_selected:
          "",
        day_today: "",
        day_outside:
          "text-zinc-500",
        day_disabled: "",
        day_range_middle:
          "",
        day_hidden: "",
        ...classNames,
      }}
      components={{
        IconLeft: ({ className, ...props }) => (
          <ChevronLeft className={cn("size-4", className)} {...props} />
        ),
        IconRight: ({ className, ...props }) => (
          <ChevronRight className={cn("size-4", className)} {...props} />
        ),
      }}
      {...props}
    />
  )
}

export { Calendar }
