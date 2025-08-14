import { Monitor, Shield, FileText, Settings, HardDrive } from "lucide-react";
import { DashboardTab } from "@/components/tabs/DashboardTab";
import { DevicesTab } from "@/components/tabs/DevicesTab";
import { SecurityTab } from "@/components/tabs/SecurityTab";
import { ReportsTab } from "@/components/tabs/ReportsTab";
import { SettingsTab } from "@/components/tabs/SettingsTab";

export const tabs = [
  {
    id: "dashboard",
    name: "Dashboard",
    icon: Monitor,
    content: DashboardTab,
  },
  {
    id: "devices",
    name: "Devices",
    icon: HardDrive,
    content: DevicesTab,
  },
  {
    id: "security",
    name: "Security",
    icon: Shield,
    content: SecurityTab,
  },
  {
    id: "reports",
    name: "Reports",
    icon: FileText,
    content: ReportsTab,
  },
  {
    id: "settings",
    name: "Settings",
    icon: Settings,
    content: SettingsTab,
  },
] as const;

export type TabId = (typeof tabs)[number]["id"];
