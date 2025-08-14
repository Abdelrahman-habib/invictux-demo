import React from "react";
import { Copy, Square, X } from "lucide-react";
import {
  WindowMinimise,
  Quit,
  WindowToggleMaximise,
} from "@/../wailsjs/runtime/runtime";
import { Button } from "@/components/ui/button";
import { useTitle } from "@/hooks/use-title";

export function ApptHeader() {
  const { title } = useTitle();
  const [isMaximised, setIsMaximised] = React.useState(true);

  const handleClose = () => {
    Quit();
  };

  const handleMinimize = () => {
    WindowMinimise();
  };

  const handleToggleMaximise = () => {
    setIsMaximised((prev) => !prev);
    WindowToggleMaximise();
  };

  return (
    <div
      className="flex items-center justify-between px-3 py-2 border-b"
      style={{ "--wails-draggable": "drag" } as React.CSSProperties}
    >
      <div className="flex-1 flex items-center gap-2 select-none cursor-move">
        <span className="text-sm text-muted-foreground line-clamp-1">
          {title}
        </span>
      </div>
      <div
        className="flex items-center gap-1"
        style={{ "--wails-draggable": "no-drag" } as React.CSSProperties}
      >
        <Button
          variant="ghost"
          onClick={handleMinimize}
          className="transition-colors"
        >
          <div className="w-3 h-0.5 bg-muted-foreground" />
        </Button>
        <Button
          variant="ghost"
          onClick={handleToggleMaximise}
          className="transition-colors"
        >
          {isMaximised ? (
            <Square className="size-3" />
          ) : (
            <>
              <Copy className="size-3 rotate-90" />
            </>
          )}
        </Button>
        <Button
          variant="ghost"
          onClick={handleClose}
          className="hover:bg-destructive/50 transition-colors"
        >
          <X className="w-3 h-3 hover:text-destructive-foreground" />
        </Button>
      </div>
    </div>
  );
}
