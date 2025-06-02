import {
  Home,
  Network,
  // Vibrate,
  ArrowUpCircleIcon,
  HelpCircleIcon,
  SettingsIcon,
  Vibrate,
} from "lucide-react";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { NavUser } from "./nav-user";
import { NavMain } from "./nav-main";
import { NavSecondary } from "./nav-secondary";
import { useAuthStore } from "@/store/auth";

const data = {
  user: {
    name: "shadcn",
    email: "m@example.com",
    avatar: "/avatars/shadcn.jpg",
  },
  navMain: [
    {
      title: "Monitors",
      url: "/monitors",
      icon: Home,
    },
    // {
    //   title: "Status pages",
    //   url: "/status-pages",
    //   icon: ListCheckIcon,
    // },
    {
      title: "Proxies",
      url: "/proxies",
      icon: Network,
    },
    {
      title: "Notifiers",
      url: "/notifiers",
      icon: Vibrate,
    },
  ],
  navSecondary: [
    {
      title: "Settings",
      url: "/settings",
      icon: SettingsIcon,
    },
    {
      title: "Get Help",
      url: "",
      icon: HelpCircleIcon,
    },
  ],
};

export function AppSidebar(props: React.ComponentProps<typeof Sidebar>) {
  const user = useAuthStore((state) => state.user);
  return (
    <Sidebar collapsible="offcanvas" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              asChild
              className="data-[slot=sidebar-menu-button]:!p-1.5"
            >
              <a href="/">
                <ArrowUpCircleIcon className="h-5 w-5" />
                <span className="text-base font-semibold">Peekaping</span>
              </a>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarContent>
        <NavMain items={data.navMain} />
        <NavSecondary items={data.navSecondary} className="mt-auto" />
      </SidebarContent>

      <SidebarFooter>
        {user && (
          <NavUser
            user={{
              name: user.email!,
              email: user.email!,
            }}
          />
        )}
      </SidebarFooter>
    </Sidebar>
  );
}
