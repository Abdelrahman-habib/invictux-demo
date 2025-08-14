import { LucideIcon } from "lucide-react";

export interface Tab {
  id: string;
  name: string;
  icon: LucideIcon;
  content: React.ComponentType;
}

export type TabId =
  | "dashboard"
  | "devices"
  | "security"
  | "reports"
  | "settings";
