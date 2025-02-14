import type { Meta, StoryObj } from '@storybook/react'
import * as Lucide from 'lucide-react'
import { Breadcrumb, BreadcrumbPage, BreadcrumbSeparator } from '../breadcrumb/breadcrumb'
import { BreadcrumbItem, BreadcrumbLink, BreadcrumbList } from '../breadcrumb/breadcrumb'
import { Button } from '../button/button'
import { Separator } from '../separator/separator'
import { Sidebar, SidebarTrigger } from './sidebar'
import { SidebarContent, SidebarInput, SidebarInset } from './sidebar-content'
import { SidebarFooter, SidebarHeader } from './sidebar-content'
import { SidebarGroup, SidebarGroupContent, SidebarGroupLabel } from './sidebar-group'
import { SidebarMenu, SidebarMenuButton, SidebarMenuItem } from './sidebar-menu'
import { SidebarProvider } from './sidebar-provider'

function LayoutWrapper({ children }: { children: React.ReactNode }) {
  return (
    <SidebarProvider>
      {children}
      <SidebarInset>
        <header className="flex h-16 shrink-0 items-center gap-2 border-b px-4">
          <SidebarTrigger className="-ml-1" />
          <Separator orientation="vertical" className="mr-2 h-4" />
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbItem className="hidden md:block">
                <BreadcrumbLink href="#">Building Your Application</BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator className="hidden md:block" />
              <BreadcrumbItem>
                <BreadcrumbPage>Data Fetching</BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
        </header>
        <div className="flex flex-1 flex-col gap-4 p-4">
          <div className="grid auto-rows-min gap-4 md:grid-cols-3">
            <div className="aspect-video rounded-xl bg-muted/50" />
            <div className="aspect-video rounded-xl bg-muted/50" />
            <div className="aspect-video rounded-xl bg-muted/50" />
          </div>
          <div className="min-h-[100vh] flex-1 rounded-xl bg-muted/50 md:min-h-min" />
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}

const meta: Meta = {
  title: 'Layout Components/Sidebar',
  component: Sidebar,
  parameters: {
    layout: 'fullscreen',
  },
  decorators: [
    (Story) => (
      <LayoutWrapper>
        {/* @ts-ignore */}
        <Story />
      </LayoutWrapper>
    ),
  ],
} satisfies Meta<typeof Sidebar>

export default meta
type Story = StoryObj<typeof meta>

const items = [
  {
    title: 'Home',
    url: '#',
    icon: Lucide.Home,
  },
  {
    title: 'Inbox',
    url: '#',
    icon: Lucide.Inbox,
  },
  {
    title: 'Calendar',
    url: '#',
    icon: Lucide.Calendar,
  },
  {
    title: 'Search',
    url: '#',
    icon: Lucide.Search,
  },
  {
    title: 'Settings',
    url: '#',
    icon: Lucide.Settings,
  },
]

export const Default: Story = {
  render: () => (
    <Sidebar>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Application</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {items.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild>
                    <a href={item.url}>
                      <item.icon className="size-4" />
                      <span>{item.title}</span>
                    </a>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
  ),
}

export const CollapsibleIcon: Story = {
  render: () => (
    <Sidebar variant="sidebar" collapsible="icon">
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Application</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {items.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild>
                    <a href={item.url}>
                      <item.icon className="size-4" />
                      <span>{item.title}</span>
                    </a>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
  ),
}

export const WithSearch: Story = {
  render: () => (
    <Sidebar>
      <SidebarHeader>
        <SidebarInput placeholder="Search..." />
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Application</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {items.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild>
                    <a href={item.url}>
                      <item.icon className="size-4" />
                      <span>{item.title}</span>
                    </a>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
  ),
}

export const WithFooter: Story = {
  render: () => (
    <Sidebar>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Application</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {items.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild>
                    <a href={item.url}>
                      <item.icon className="size-4" />
                      <span>{item.title}</span>
                    </a>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <Button variant="ghost" className="w-full justify-start">
          <Lucide.LogOut className="mr-2 size-4" />
          Logout
        </Button>
      </SidebarFooter>
    </Sidebar>
  ),
}
