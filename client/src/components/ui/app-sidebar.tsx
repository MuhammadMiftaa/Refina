"use client";

import { NavMain } from "@/components/ui/nav-main";
import { NavProjects } from "@/components/ui/nav-projects";
import { NavUser } from "@/components/ui/nav-user";
import { TeamSwitcher } from "@/components/ui/team-switcher";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";
import { useProfile } from "@/store/useProfile";
import { useShallow } from "zustand/shallow";
import { data } from "@/helper/Data";

// This is sample data.
export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const { username, email } = useProfile(
    useShallow((state) => ({
      username: state.username,
      email: state.email,
    })),
  );

  return (
    <Sidebar className="font-poppins" collapsible="icon" {...props}>
      <SidebarHeader>
        <TeamSwitcher />
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} />
        <NavProjects projects={data.projects} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={{ ...data.user, name: username, email: email }} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
