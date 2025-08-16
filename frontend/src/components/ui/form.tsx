"use client";

import * as React from "react";
import * as LabelPrimitive from "@radix-ui/react-label";
import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";
import {
  Controller,
  ControllerProps,
  FieldPath,
  FieldValues,
  FormProvider,
  useFormContext,
} from "react-hook-form";

import { cn } from "@/lib/utils";
import { Label } from "@/components/ui/label";

const formLabelVariants = cva("transition-colors duration-300", {
  variants: {
    disabled: {
      true: "cursor-not-allowed opacity-50",
      false: "",
    },
    variant: {
      default: "",
      outline:
        "border-t border-l border-r border-input rounded-t-md px-2 pb-1.5 text-foreground/80",
    },
    focusVariant: {
      default: "",
      outline: "bg-primary/20 border-primary/50 text-foreground",
    },
  },
  defaultVariants: {
    variant: "default",
  },
});

const Form = FormProvider;

type FormFieldContextValue<
  TFieldValues extends FieldValues = FieldValues,
  TName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>
> = {
  name: TName;
};

const FormFieldContext = React.createContext<FormFieldContextValue>(
  {} as FormFieldContextValue
);

const FormField = <
  TFieldValues extends FieldValues = FieldValues,
  TName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>
>({
  ...props
}: ControllerProps<TFieldValues, TName>) => {
  return (
    <FormFieldContext.Provider value={{ name: props.name }}>
      <Controller {...props} />
    </FormFieldContext.Provider>
  );
};

const useFormField = () => {
  const [isFocused, setIsFocused] = React.useState(false);
  const fieldContext = React.useContext(FormFieldContext);
  const itemContext = React.useContext(FormItemContext);
  // Removed isSelectOpen from FormControlContext as per generic solution
  const { getFieldState, formState } = useFormContext();

  const fieldState = getFieldState(fieldContext.name, formState);

  if (!fieldContext) {
    throw new Error("useFormField should be used within <FormField>");
  }

  const { id, formItemRef } = itemContext;

  React.useEffect(() => {
    const checkFocus = () => {
      const isElementFocused =
        formItemRef.current &&
        formItemRef.current.contains(document.activeElement);

      // Check if any descendant element within the FormItem has data-state="open"
      const isPopoverOpen =
        formItemRef.current &&
        formItemRef.current.querySelector('[data-state="open"]');

      setIsFocused(isElementFocused || !!isPopoverOpen);
    };

    // Listen for focus changes on the whole window
    window.addEventListener("focusin", checkFocus);
    window.addEventListener("focusout", checkFocus);

    checkFocus();

    return () => {
      window.removeEventListener("focusin", checkFocus);
      window.removeEventListener("focusout", checkFocus);
    };
  }, [formItemRef]); // Removed isSelectOpen from dependencies
  return {
    id,
    name: fieldContext.name,
    formItemId: `${id}-form-item`,
    formDescriptionId: `${id}-form-item-description`,
    formMessageId: `${id}-form-item-message`,
    ...fieldState,
    isFocused,
  };
};

type FormItemContextValue = {
  id: string;
  labelVariant: VariantProps<typeof formLabelVariants>["variant"];
  setLabelVariant: (
    variant: VariantProps<typeof formLabelVariants>["variant"]
  ) => void;
  formItemRef: React.MutableRefObject<HTMLDivElement | null>;
};

const FormItemContext = React.createContext<FormItemContextValue>(
  {} as FormItemContextValue
);

const FormItem = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => {
  const id = React.useId();
  const [labelVariant, setLabelVariant] =
    React.useState<VariantProps<typeof formLabelVariants>["variant"]>(
      "default"
    );
  const formItemRef = React.useRef<HTMLDivElement | null>(null);

  return (
    <FormItemContext.Provider
      value={{ id, labelVariant, setLabelVariant, formItemRef }}
    >
      <div
        ref={(node) => {
          // Forward the ref
          if (typeof ref === "function") {
            ref(node);
          } else if (ref) {
            (ref as React.MutableRefObject<HTMLDivElement | null>).current =
              node;
          }
          // Set the internal ref
          formItemRef.current = node;
        }}
        className={cn("space-y-2 group", className)}
        {...props}
      />
    </FormItemContext.Provider>
  );
});
FormItem.displayName = "FormItem";

interface FormLabelProps
  extends React.ComponentPropsWithoutRef<typeof LabelPrimitive.Root>,
    VariantProps<typeof formLabelVariants> {}

const FormLabel = React.forwardRef<
  React.ElementRef<typeof LabelPrimitive.Root>,
  FormLabelProps & { disabled?: boolean }
>(({ className, variant, focusVariant, disabled, ...props }, ref) => {
  const { error, formItemId, isFocused } = useFormField();
  const { setLabelVariant } = React.useContext(FormItemContext);

  React.useEffect(() => {
    setLabelVariant(variant);
  }, [variant, setLabelVariant]);

  return (
    <Label
      ref={ref}
      className={cn(
        formLabelVariants({ variant, disabled }),
        error && "text-destructive",
        isFocused && formLabelVariants({ focusVariant }),
        className
      )}
      htmlFor={formItemId}
      {...props}
    />
  );
});
FormLabel.displayName = "FormLabel";

const FormControl = React.forwardRef<
  React.ElementRef<typeof Slot>,
  React.ComponentPropsWithoutRef<typeof Slot>
>(({ className, ...props }, ref) => {
  const { error, formItemId, formDescriptionId, formMessageId } =
    useFormField();
  const { labelVariant } = React.useContext(FormItemContext);

  return (
    <Slot
      ref={ref}
      id={formItemId}
      aria-describedby={
        !error
          ? `${formDescriptionId}`
          : `${formDescriptionId} ${formMessageId}`
      }
      aria-invalid={!!error}
      className={cn(labelVariant === "outline" && "rounded-tr-none", className)}
      {...props}
    />
  );
});
FormControl.displayName = "FormControl";

const FormDescription = React.forwardRef<
  HTMLParagraphElement,
  React.HTMLAttributes<HTMLParagraphElement> & { asChild?: boolean }
>(({ asChild, className, ...props }, ref) => {
  const { formDescriptionId } = useFormField();
  const Comp = asChild ? Slot : "p";
  return (
    <Comp
      ref={ref}
      id={formDescriptionId}
      className={cn("text-sm text-muted-foreground", className)}
      {...props}
    />
  );
});
FormDescription.displayName = "FormDescription";

const FormMessage = React.forwardRef<
  HTMLParagraphElement,
  React.HTMLAttributes<HTMLParagraphElement>
>(({ className, children, ...props }, ref) => {
  const { error, formMessageId } = useFormField();
  const body = error ? String(error?.message) : children;

  if (!body) {
    return null;
  }

  return (
    <p
      ref={ref}
      id={formMessageId}
      className={cn("text-sm font-medium text-destructive", className)}
      {...props}
    >
      {body}
    </p>
  );
});
FormMessage.displayName = "FormMessage";

export {
  useFormField,
  Form,
  FormItem,
  FormLabel,
  FormControl,
  FormDescription,
  FormMessage,
  FormField,
};
