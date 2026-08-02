import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";
import Link from "next/link";
import type { ButtonHTMLAttributes, ReactNode } from "react";

import { cn } from "@/lib/utils";

const buttonVariants = cva(
  "inline-flex min-h-12 items-center justify-center gap-2 rounded-sm border px-6 text-sm font-medium transition-colors duration-150 disabled:pointer-events-none disabled:opacity-50",
  {
    variants: {
      variant: {
        primary:
          "border-foreground bg-foreground text-background hover:border-primary hover:bg-primary",
        secondary:
          "border-secondary bg-secondary text-secondary-foreground hover:border-foreground hover:bg-foreground",
        outline:
          "border-current bg-transparent hover:bg-foreground hover:text-background",
        ghost: "border-transparent bg-transparent hover:border-border",
        destructive:
          "border-destructive bg-destructive text-destructive-foreground hover:border-foreground hover:bg-foreground",
      },
    },
    defaultVariants: { variant: "primary" },
  },
);

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> &
  VariantProps<typeof buttonVariants> & {
    asChild?: boolean;
  };

export function Button({
  asChild = false,
  className,
  variant,
  ...props
}: ButtonProps) {
  const Component = asChild ? Slot : "button";
  return (
    <Component
      className={cn(buttonVariants({ variant }), className)}
      {...props}
    />
  );
}

export function ButtonLink({
  children,
  className,
  href,
  variant = "primary",
}: {
  children: ReactNode;
  className?: string;
  href: string;
  variant?: NonNullable<VariantProps<typeof buttonVariants>["variant"]>;
}) {
  return (
    <Button asChild className={className} variant={variant}>
      <Link href={href}>{children}</Link>
    </Button>
  );
}
