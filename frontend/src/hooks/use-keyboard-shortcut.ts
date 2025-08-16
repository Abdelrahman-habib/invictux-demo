import { useEffect, useCallback, useMemo } from "react";
import { useIsMobile } from "@/hooks/use-mobile";

type Shortcut = string[];

interface ShortcutConfig {
  shortcut: Shortcut;
  onShortcutPressed?: () => void;
  onShortcutDown?: () => void;
  onShortcutUp?: () => void;
  displayText?: string;
}

type ShortcutConfigs = ShortcutConfig | ShortcutConfig[];

export function useKeyboardShortcut(configs: ShortcutConfigs) {
  const isMobile = useIsMobile();
  const configsArray = useMemo(
    () => (Array.isArray(configs) ? configs : [configs]),
    [configs]
  );

  useEffect(() => {
    if (isMobile) return;

    const mouseButtonMap: Record<string, number> = {
      leftclick: 0,
      middleclick: 1,
      rightclick: 2,
    };
    const mouseShortcutKeys = Object.keys(mouseButtonMap);
    const hasMouseShortcut = configsArray.some((config) =>
      config.shortcut.some((key) =>
        mouseShortcutKeys.includes(key?.toLowerCase())
      )
    );

    const handleKeyPress = (event: KeyboardEvent) => {
      for (const config of configsArray) {
        const { shortcut: keys } = config;
        // Only handle if not a keyboard shortcut
        if (keys.some((key) => mouseShortcutKeys.includes(key.toLowerCase())))
          continue;
        const isShiftRequired = keys.some(
          (key) => key.toLowerCase() === "shift"
        );
        const isCtrlRequired =
          keys.some((key) => key.toLowerCase() === "ctrl") ||
          keys.some((key) => key.toLowerCase() === "meta");
        const isAltRequired = keys.some((key) => key.toLowerCase() === "alt");
        const targetKey = keys[keys.length - 1];
        const shiftMatches = isShiftRequired === event.shiftKey;
        const ctrlMatches = isCtrlRequired === (event.metaKey || event.ctrlKey);
        const altMatches = isAltRequired === event.altKey;
        const keyMatches = event.key?.toLowerCase() === targetKey.toLowerCase();
        if (shiftMatches && ctrlMatches && altMatches && keyMatches) {
          event.preventDefault();
          if (config.onShortcutPressed) config.onShortcutPressed();
          if (config.onShortcutDown) config.onShortcutDown();
          break;
        }
      }
    };

    const handleKeyUp = (event: KeyboardEvent) => {
      for (const config of configsArray) {
        const { shortcut: keys } = config;
        if (keys.some((key) => mouseShortcutKeys.includes(key.toLowerCase())))
          continue;
        const isShiftRequired = keys.some(
          (key) => key.toLowerCase() === "shift"
        );
        const isCtrlRequired =
          keys.some((key) => key.toLowerCase() === "ctrl") ||
          keys.some((key) => key.toLowerCase() === "meta");
        const isAltRequired = keys.some((key) => key.toLowerCase() === "alt");
        const targetKey = keys[keys.length - 1];
        const shiftMatches = isShiftRequired === event.shiftKey;
        const ctrlMatches = isCtrlRequired === (event.metaKey || event.ctrlKey);
        const altMatches = isAltRequired === event.altKey;
        const keyMatches = event.key?.toLowerCase() === targetKey.toLowerCase();
        if (shiftMatches && ctrlMatches && altMatches && keyMatches) {
          if (config.onShortcutUp) config.onShortcutUp();
          break;
        }
      }
    };

    const handleMouseDown = (event: MouseEvent) => {
      for (const config of configsArray) {
        const { shortcut: keys } = config;
        const mouseKey = keys.find((key) =>
          mouseShortcutKeys.includes(key.toLowerCase())
        );
        if (!mouseKey) continue;
        const isShiftRequired = keys.some(
          (key) => key.toLowerCase() === "shift"
        );
        const isCtrlRequired =
          keys.some((key) => key.toLowerCase() === "ctrl") ||
          keys.some((key) => key.toLowerCase() === "meta");
        const isAltRequired = keys.some((key) => key.toLowerCase() === "alt");
        const buttonMatches =
          event.button === mouseButtonMap[mouseKey.toLowerCase()];
        const shiftMatches = isShiftRequired === event.shiftKey;
        const ctrlMatches = isCtrlRequired === (event.metaKey || event.ctrlKey);
        const altMatches = isAltRequired === event.altKey;
        if (buttonMatches && shiftMatches && ctrlMatches && altMatches) {
          event.preventDefault();
          if (config.onShortcutPressed) config.onShortcutPressed();
          if (config.onShortcutDown) config.onShortcutDown();
          break;
        }
      }
    };

    const handleMouseUp = (event: MouseEvent) => {
      for (const config of configsArray) {
        const { shortcut: keys } = config;
        const mouseKey = keys.find((key) =>
          mouseShortcutKeys.includes(key.toLowerCase())
        );
        if (!mouseKey) continue;
        const isShiftRequired = keys.some(
          (key) => key.toLowerCase() === "shift"
        );
        const isCtrlRequired =
          keys.some((key) => key.toLowerCase() === "ctrl") ||
          keys.some((key) => key.toLowerCase() === "meta");
        const isAltRequired = keys.some((key) => key.toLowerCase() === "alt");
        const buttonMatches =
          event.button === mouseButtonMap[mouseKey.toLowerCase()];
        const shiftMatches = isShiftRequired === event.shiftKey;
        const ctrlMatches = isCtrlRequired === (event.metaKey || event.ctrlKey);
        const altMatches = isAltRequired === event.altKey;
        if (buttonMatches && shiftMatches && ctrlMatches && altMatches) {
          if (config.onShortcutUp) config.onShortcutUp();
          break;
        }
      }
    };

    window.addEventListener("keydown", handleKeyPress);
    window.addEventListener("keyup", handleKeyUp);
    if (hasMouseShortcut) {
      window.addEventListener("mousedown", handleMouseDown);
      window.addEventListener("mouseup", handleMouseUp);
    }
    return () => {
      window.removeEventListener("keydown", handleKeyPress);
      window.removeEventListener("keyup", handleKeyUp);
      if (hasMouseShortcut) {
        window.removeEventListener("mousedown", handleMouseDown);
        window.removeEventListener("mouseup", handleMouseUp);
      }
    };
  }, [configsArray, isMobile]);

  const getDisplayTexts = useCallback(() => {
    return configsArray.map((config) => {
      if (config.displayText) return config.displayText;

      return config.shortcut
        .map((key) => {
          switch (key) {
            case "Shift":
              return "⇧";
            case "Ctrl":
              return "⌘";
            case "Alt":
              return "⌥";
            case "Meta":
              return "⌘";
            default:
              return key;
          }
        })
        .join(" ");
    });
  }, [configsArray]);

  return { isMobile, displayTexts: getDisplayTexts() };
}
