import { cn } from "@/lib/utils";
import { NumericFormat, type NumericFormatProps } from "react-number-format";
import { LucideIcon } from "lucide-react";

export interface NumberInputProps extends NumericFormatProps {
  Icon?: LucideIcon | string;
  iconPosition?: "start" | "end";
}

export function NumberInput({
  thousandSeparator = true,
  customInput,
  className,
  Icon,
  iconPosition = "end",
  ...props
}: NumberInputProps) {
  return (
    <div className="relative">
      <NumericFormat
        className={cn(
          "peer",
          iconPosition === "start" && "ps-12",
          iconPosition === "end" && "pe-12",
          className
        )}
        thousandSeparator={thousandSeparator}
        customInput={customInput}
        {...props}
      />
      {Icon && (
        <span
          className={cn(
            "text-muted-foreground pointer-events-none absolute inset-y-0 flex items-center justify-center text-sm peer-disabled:opacity-50 bg-gradient-to-r from-background to-transparent to-90% z-10 rounded-full px-2 my-px",
            iconPosition === "start" && "start-2 ps-2 pe-8",
            iconPosition === "end" && "end-2 ps-8 pe-2"
          )}
        >
          {typeof Icon === "string" ? Icon : <Icon className="h-4 w-4" />}
        </span>
      )}
    </div>
  );
}
