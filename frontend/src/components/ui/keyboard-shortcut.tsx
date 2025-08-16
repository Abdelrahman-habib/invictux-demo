import React from "react";
import { cn } from "@/lib/utils";
import { useKeyboardShortcut } from "@/hooks/use-keyboard-shortcut";

interface KeyboardShortcutProps {
  keys: string[];
  onKeyPressed: () => void;
  className?: string;
  displayText?: string;
}

export const KeyboardShortcut = React.memo(function KeyboardShortcut({
  keys,
  onKeyPressed,
  className,
  displayText,
}: KeyboardShortcutProps) {
  const { isMobile, displayTexts: formattedDisplayText } = useKeyboardShortcut({
    shortcut: keys,
    onShortcutPressed: onKeyPressed,
    displayText,
  });

  if (isMobile) return null;

  return (
    <kbd
      dir="ltr"
      className={cn(
        "select-none h-5 gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium text-muted-foreground inline-flex items-center justify-center",
        className
      )}
    >
      {formattedDisplayText.join(" + ")}
    </kbd>
  );
});
