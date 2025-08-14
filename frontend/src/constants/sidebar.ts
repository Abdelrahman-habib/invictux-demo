import { Monitor, Shield, FileText, Settings, HardDrive } from "lucide-react";

export const sidebarList = [
  {
    id: "dashboard",
    name: "Dashboard",
    icon: Monitor,
    path: "/",
  },
  {
    id: "devices",
    name: "Devices",
    icon: HardDrive,
    path: "/devices",
  },
  {
    id: "security",
    name: "Security",
    icon: Shield,
    path: "/security",
  },
  {
    id: "reports",
    name: "Reports",
    icon: FileText,
    path: "/reports",
  },
  {
    id: "settings",
    name: "Settings",
    icon: Settings,
    path: "/settings",
  },
] as const;
