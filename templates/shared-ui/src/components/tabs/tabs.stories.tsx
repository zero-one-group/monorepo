import type { Meta, StoryObj } from '@storybook/react'
import * as Lucide from 'lucide-react'
import { Button } from '../button/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '../card/card'
import { Input } from '../input/input'
import { Label } from '../label/label'
import { ScrollArea, ScrollBar } from '../scroll-area/scroll-area'
import { Text } from '../text/text'
import { Tabs, TabsContent, TabsList, TabsTrigger } from './tabs'

const meta: Meta<typeof Tabs> = {
  title: 'Basic Components/Tabs',
  component: Tabs,
}

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  render: (_args) => (
    <Tabs defaultValue="account" className="w-[400px]">
      <TabsList className="grid w-full grid-cols-2">
        <TabsTrigger value="account">Account</TabsTrigger>
        <TabsTrigger value="password">Password</TabsTrigger>
      </TabsList>
      <TabsContent value="account">
        <Card>
          <CardHeader>
            <CardTitle>Account</CardTitle>
            <CardDescription>
              Make changes to your account here. Click save when you're done.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-2">
            <div className="space-y-1">
              <Label htmlFor="name">Name</Label>
              <Input id="name" defaultValue="Aris Ripandi" />
            </div>
            <div className="space-y-1">
              <Label htmlFor="username">Username</Label>
              <Input id="username" defaultValue="@riipandi" />
            </div>
          </CardContent>
          <CardFooter>
            <Button>Save changes</Button>
          </CardFooter>
        </Card>
      </TabsContent>
      <TabsContent value="password">
        <Card>
          <CardHeader>
            <CardTitle>Password</CardTitle>
            <CardDescription>
              Change your password here. After saving, you'll be logged out.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-2">
            <div className="space-y-1">
              <Label htmlFor="current">Current password</Label>
              <Input id="current" type="password" />
            </div>
            <div className="space-y-1">
              <Label htmlFor="new">New password</Label>
              <Input id="new" type="password" />
            </div>
          </CardContent>
          <CardFooter>
            <Button>Save password</Button>
          </CardFooter>
        </Card>
      </TabsContent>
    </Tabs>
  ),
}

// @ref: https://21st.dev/originui/tabs-like-bookmark-for-navigation/default
export const BookmarkStyle: Story = {
  render: (_args) => (
    <Tabs defaultValue="tab-1">
      <ScrollArea>
        <TabsList className="relative h-auto w-full justify-start gap-0.5 bg-transparent p-0 before:absolute before:inset-x-0 before:bottom-0 before:h-px before:bg-border">
          <TabsTrigger
            value="tab-1"
            className="overflow-hidden rounded-t-sm rounded-b-none border-border border-x border-t bg-muted py-2 text-xs data-[state=active]:z-10 data-[state=active]:shadow-none"
          >
            <Lucide.House
              className="-ms-0.5 me-1.5 size-4 opacity-60"
              strokeWidth={2}
              aria-hidden="true"
            />
            Overview
          </TabsTrigger>
          <TabsTrigger
            value="tab-2"
            className="overflow-hidden rounded-t-sm rounded-b-none border-border border-x border-t bg-muted py-2 text-xs data-[state=active]:z-10 data-[state=active]:shadow-none"
          >
            <Lucide.PanelsTopLeft
              className="-ms-0.5 me-1.5 size-4 opacity-60"
              strokeWidth={2}
              aria-hidden="true"
            />
            Projects
          </TabsTrigger>
          <TabsTrigger
            value="tab-3"
            className="overflow-hidden rounded-t-sm rounded-b-none border-border border-x border-t bg-muted py-2 text-xs data-[state=active]:z-10 data-[state=active]:shadow-none"
          >
            <Lucide.Box
              className="-ms-0.5 me-1.5 size-4 opacity-60"
              strokeWidth={2}
              aria-hidden="true"
            />
            Packages
          </TabsTrigger>
        </TabsList>
        <ScrollBar orientation="horizontal" />
      </ScrollArea>
      <TabsContent value="tab-1" className="p-2">
        <Text size="sm">Content for Tab 1</Text>
      </TabsContent>
      <TabsContent value="tab-2" className="p-2">
        <Text size="sm">Content for Tab 2</Text>
      </TabsContent>
      <TabsContent value="tab-3" className="p-2">
        <Text size="sm">Content for Tab 3</Text>
      </TabsContent>
    </Tabs>
  ),
}

// @ref: https://21st.dev/k3menn/tabs-in-cell-for-navigation/default
export const TabsInCell: Story = {
  render: (_args) => (
    <Tabs defaultValue="tab-1">
      <ScrollArea>
        <TabsList className="-space-x-px mb-3 h-auto justify-start bg-background p-0 shadow-black/5 shadow-sm rtl:space-x-reverse">
          <TabsTrigger
            value="tab-1"
            className="relative overflow-hidden rounded-none border border-border py-2 after:pointer-events-none after:absolute after:inset-x-0 after:bottom-0 after:h-0.5 first:rounded-s last:rounded-e data-[state=active]:bg-muted data-[state=active]:after:bg-primary"
          >
            <Lucide.House
              className="-ms-0.5 me-1.5 opacity-60"
              size={16}
              strokeWidth={2}
              aria-hidden="true"
            />
            Overview
          </TabsTrigger>
          <TabsTrigger
            value="tab-2"
            className="relative overflow-hidden rounded-none border border-border py-2 after:pointer-events-none after:absolute after:inset-x-0 after:bottom-0 after:h-0.5 first:rounded-s last:rounded-e data-[state=active]:bg-muted data-[state=active]:after:bg-primary"
          >
            <Lucide.PanelsTopLeft
              className="-ms-0.5 me-1.5 opacity-60"
              size={16}
              strokeWidth={2}
              aria-hidden="true"
            />
            Projects
          </TabsTrigger>
          <TabsTrigger
            value="tab-3"
            className="relative overflow-hidden rounded-none border border-border py-2 after:pointer-events-none after:absolute after:inset-x-0 after:bottom-0 after:h-0.5 first:rounded-s last:rounded-e data-[state=active]:bg-muted data-[state=active]:after:bg-primary"
          >
            <Lucide.Box
              className="-ms-0.5 me-1.5 opacity-60"
              size={16}
              strokeWidth={2}
              aria-hidden="true"
            />
            Packages
          </TabsTrigger>
        </TabsList>
        <ScrollBar orientation="horizontal" />
      </ScrollArea>
      <TabsContent value="tab-1" className="p-2">
        <Text size="sm">Content for Tab 1</Text>
      </TabsContent>
      <TabsContent value="tab-2" className="p-2">
        <Text size="sm">Content for Tab 2</Text>
      </TabsContent>
      <TabsContent value="tab-3" className="p-2">
        <Text size="sm">Content for Tab 3</Text>
      </TabsContent>
    </Tabs>
  ),
}
