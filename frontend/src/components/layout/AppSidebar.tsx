
import React, { useState, useEffect } from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { getNavItemsByRole } from '@/config/nav';
import { cn } from '@/lib/utils';
import { ChevronRight, LayoutDashboard, X } from 'lucide-react';

interface AppSidebarProps {
    className?: string;
    isMobileOpen: boolean;
    onMobileClose: () => void;
}

export const AppSidebar: React.FC<AppSidebarProps> = ({ className, isMobileOpen, onMobileClose }) => {
    const { user } = useAuth();
    const location = useLocation();
    const navItems = getNavItemsByRole(user?.role);
    const [collapsed, setCollapsed] = useState(false);

    // Close mobile menu on route change
    useEffect(() => {
        onMobileClose();
    }, [location.pathname]);

    return (
        <>
            {/* Mobile Overlay */}
            {isMobileOpen && (
                <div
                    className="fixed inset-0 z-40 bg-black/50 md:hidden transition-opacity"
                    onClick={onMobileClose}
                />
            )}

            {/* Sidebar Container */}
            <aside
                className={cn(
                    "fixed top-0 left-0 z-50 h-[100dvh] bg-white border-r border-gray-200 transition-all duration-300 ease-in-out md:static md:h-screen md:translate-x-0",
                    isMobileOpen ? "translate-x-0 w-64 shadow-xl" : "-translate-x-full md:translate-x-0",
                    collapsed ? "md:w-16" : "md:w-64",
                    className
                )}
            >
                {/* Sidebar Header */}
                <div className="flex h-16 items-center justify-between border-b px-4">
                    {!collapsed && (
                        <div className="flex items-center gap-2 font-bold text-xl text-primary-600">
                            {/* Replace with Logo if available */}
                            <LayoutDashboard className="h-6 w-6" />
                            <span className="truncate">EduCenter</span>
                        </div>
                    )}
                    {collapsed && (
                        <div className="mx-auto">
                            <LayoutDashboard className="h-6 w-6 text-primary-600" />
                        </div>
                    )}

                    {/* Mobile Close Button */}
                    <button onClick={onMobileClose} className="md:hidden text-gray-500 hover:text-gray-700">
                        <X size={24} />
                    </button>

                    {/* Desktop Collapse Toggle (Optional, can be moved to bottom or header) */}
                    {/* For now, simplified: keeping static width on desktop or using a toggle if required. 
                        User requirement: "collapse/expand trên desktop".
                        Let's add a toggle button at the bottom.
                    */}
                </div>

                {/* Navigation Items */}
                <div className="flex flex-col gap-1 p-2 overflow-y-auto h-[calc(100vh-4rem-3rem)]">
                    {navItems.map((item) => {
                        const Icon = item.icon;
                        const isActive = location.pathname.startsWith(item.path);

                        return (
                            <NavLink
                                key={item.path}
                                to={item.path}
                                className={({ isActive }) => cn(
                                    "flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-all hover:bg-gray-100",
                                    isActive ? "bg-blue-50 text-blue-700 hover:bg-blue-100" : "text-gray-700",
                                    collapsed ? "justify-center px-2" : ""
                                )}
                                title={collapsed ? item.label : undefined}
                            >
                                <Icon size={20} className={cn("shrink-0", isActive && "text-blue-600")} />
                                {!collapsed && <span>{item.label}</span>}
                                {!collapsed && isActive && <ChevronRight className="ml-auto h-4 w-4 text-blue-500 opacity-50" />}
                            </NavLink>
                        );
                    })}
                </div>

                {/* Sidebar Footer / Collapse Toggle */}
                <div className="border-t p-3 hidden md:flex">
                    <button
                        onClick={() => setCollapsed(!collapsed)}
                        className={cn(
                            "flex w-full items-center justify-center rounded-md p-2 text-gray-500 hover:bg-gray-100 transition-colors",
                            !collapsed && "justify-end"
                        )}
                    >
                        {collapsed ? <ChevronRight size={20} /> : <ChevronRight size={20} className="rotate-180" />}
                    </button>
                </div>
            </aside>
        </>
    );
};
