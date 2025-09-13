import { Bot } from "lucide-react";

import { NavLink } from "./nav-link";
import { ThemeToggle } from "./theme/theme-toggle";
import { Separator } from "./ui/separator";

export function Header() {
  return (
    <div className="border-b">
      <div className="flex h-16 items-center gap-6 px-6">
        <Bot className="h-7 w-7" />
        <Separator orientation="vertical" className="h-6" />
        <nav className="flex items-center space-x-4 lg:space-x-6">
          <NavLink to="/nginx">
            <img
              className="h-4 w-4 grayscale"
              src="/logos/nginx.svg"
              alt="Nginx Logo"
            />
            Nginx
          </NavLink>
        </nav>

        <div className="ml-auto flex items-center gap-2">
          <ThemeToggle key="@jarvis-devops-theme" />
        </div>
      </div>
    </div>
  );
}
