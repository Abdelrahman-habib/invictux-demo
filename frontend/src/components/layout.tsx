import { ApptHeader } from "./ApptHeader";
import { TitleProvider } from "@/contexts/TitleProvider";

export function RootLayout({ children }: { children?: React.ReactNode }) {
  return (
    <TitleProvider>
      <div className="w-full h-full bg-transparent rounded-lg border text-foreground overflow-hidden">
        <ApptHeader />
        {children}
      </div>
    </TitleProvider>
  );
}
