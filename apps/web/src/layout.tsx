import {
  SidebarInset,
  SidebarProvider,
  // SidebarTrigger
} from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar";
import { SiteHeader } from "./components/app-header";

export default function Layout({
  children,
  pageName,
  error,
  isLoading,
  onCreate,
}: {
  children: React.ReactNode;
  pageName: string;
  onCreate?: () => void;
  error?: React.ReactNode;
  isLoading?: boolean;
}) {

  return (
    <SidebarProvider>
      <AppSidebar variant="inset" />
      <SidebarInset>
        <SiteHeader pageName={pageName} onCreate={onCreate} />
        {isLoading ? (
          <div className="p-4 w-full">Loading ...</div>
        ) : (
          error || <main className="p-4 w-full">{children}</main>
        )}
      </SidebarInset>
    </SidebarProvider>
  );
}
