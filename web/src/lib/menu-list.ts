import { FileRoutesByPath } from "@tanstack/react-router";
import { LayoutGrid, LucideIcon, Settings } from "lucide-react";

type Submenu = {
  href: string;
  label: string;
  active: boolean;
};

type Menu = {
  href: FileRoutesByPath[keyof FileRoutesByPath]["fullPath"];
  label: string;
  active: boolean;
  icon: LucideIcon;
  submenus: Submenu[];
};

type Group = {
  groupLabel: string;
  menus: Menu[];
};

export function getMenuList(pathname: string): Group[] {
  return [
    {
      groupLabel: "",
      menus: [
        {
          href: "/dashboard/pipelines",
          label: "Pipelines",
          active: pathname.includes("/pipelines"),
          icon: LayoutGrid,
          submenus: [],
        },
        {
          href: "/dashboard/projects",
          label: "Projects",
          active: pathname.includes("/projects"),
          icon: LayoutGrid,
          submenus: [],
        },
        {
          href: "/dashboard/settings",
          label: "Settings",
          active: pathname.includes("/settings"),
          icon: Settings,
          submenus: [],
        },
      ],
    },
  ];
}
