import type { Meta, StoryObj } from "@storybook/react";
import * as React from "react";
import { useMediaQuery } from "../../hooks/use-media-query";
import { Button } from "../button/button";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "../dialog/dialog";
import { Input } from "../input/input";
import { Label } from "../label/label";
import {
	Drawer,
	DrawerClose,
	DrawerContent,
	DrawerDescription,
	DrawerFooter,
	DrawerHeader,
	DrawerTitle,
	DrawerTrigger,
} from "./drawer";

const meta: Meta<typeof Drawer> = {
	title: "Basic Components/Drawer",
	component: Drawer,
};

export default meta;

export const Default: StoryObj = {
	render: () => (
		<Drawer>
			<DrawerTrigger asChild>
				<Button variant="outline">Open Drawer</Button>
			</DrawerTrigger>
			<DrawerContent>
				<DrawerHeader>
					<DrawerTitle>Edit Profile</DrawerTitle>
					<DrawerDescription>
						Make changes to your profile here. Click save when you're done.
					</DrawerDescription>
				</DrawerHeader>
				<div className="p-4">
					<p>Drawer content goes here</p>
				</div>
				<DrawerFooter>
					<Button>Save changes</Button>
				</DrawerFooter>
			</DrawerContent>
		</Drawer>
	),
};

export const WithForm: StoryObj = {
	render: () => (
		<Drawer>
			<DrawerTrigger asChild>
				<Button>Edit Settings</Button>
			</DrawerTrigger>
			<DrawerContent>
				<DrawerHeader>
					<DrawerTitle>Account Settings</DrawerTitle>
					<DrawerDescription>Update your account preferences</DrawerDescription>
				</DrawerHeader>
				<div className="space-y-4 p-4">
					<div className="space-y-2">
						<label htmlFor="name" className="font-medium text-sm">
							Name
						</label>
						<input
							id="name"
							className="w-full rounded-md border p-2"
							placeholder="Enter your name"
						/>
					</div>
					<div className="space-y-2">
						<label htmlFor="email" className="font-medium text-sm">
							Email
						</label>
						<input
							id="email"
							type="email"
							className="w-full rounded-md border p-2"
							placeholder="Enter your email"
						/>
					</div>
				</div>
				<DrawerFooter>
					<DrawerClose asChild>
						<Button variant="outline">Cancel</Button>
					</DrawerClose>
					<Button>Save Changes</Button>
				</DrawerFooter>
			</DrawerContent>
		</Drawer>
	),
};

export const WithCustomContent: StoryObj = {
	render: () => (
		<Drawer>
			<DrawerTrigger asChild>
				<Button variant="secondary">View Details</Button>
			</DrawerTrigger>
			<DrawerContent>
				<DrawerHeader>
					<DrawerTitle>Product Details</DrawerTitle>
					<DrawerDescription>
						View complete product information
					</DrawerDescription>
				</DrawerHeader>
				<div className="p-4">
					<div className="mb-4 aspect-video rounded-lg bg-muted" />
					<h3 className="mb-2 font-semibold text-lg">Product Name</h3>
					<p className="mb-4 text-muted-foreground">
						Lorem ipsum dolor sit amet consectetur adipisicing elit. Quisquam
						voluptates, quod, voluptatum, quae voluptatibus quas quidem
						voluptatem quibusdam quos quia nesciunt.
					</p>
					<div className="flex items-center justify-between">
						<span className="font-bold text-xl">$99.99</span>
						<Button size="sm">Add to Cart</Button>
					</div>
				</div>
				<DrawerFooter>
					<DrawerClose asChild>
						<Button variant="outline" className="w-full">
							Close
						</Button>
					</DrawerClose>
				</DrawerFooter>
			</DrawerContent>
		</Drawer>
	),
};

export const ResponsiveDialog: StoryObj = {
	render: () => <ResponsiveDialogDemo />,
};

function ResponsiveDialogDemo() {
	const [open, setOpen] = React.useState(false);
	const isDesktop = useMediaQuery("(min-width: 768px)");

	if (isDesktop) {
		return (
			<Dialog open={open} onOpenChange={setOpen}>
				<DialogTrigger asChild>
					<Button variant="outline">Edit Profile</Button>
				</DialogTrigger>
				<DialogContent className="sm:max-w-[425px]">
					<DialogHeader>
						<DialogTitle>Edit profile</DialogTitle>
						<DialogDescription>
							Make changes to your profile here. Click save when you're done.
						</DialogDescription>
					</DialogHeader>
					<ProfileForm />
				</DialogContent>
			</Dialog>
		);
	}

	return (
		<Drawer open={open} onOpenChange={setOpen}>
			<DrawerTrigger asChild>
				<Button variant="outline">Edit Profile</Button>
			</DrawerTrigger>
			<DrawerContent>
				<DrawerHeader className="text-left">
					<DrawerTitle>Edit profile</DrawerTitle>
					<DrawerDescription>
						Make changes to your profile here. Click save when you're done.
					</DrawerDescription>
				</DrawerHeader>
				<ProfileForm className="px-4" />
				<DrawerFooter className="pt-2">
					<DrawerClose asChild>
						<Button variant="outline">Cancel</Button>
					</DrawerClose>
				</DrawerFooter>
			</DrawerContent>
		</Drawer>
	);
}

function ProfileForm({ className }: React.ComponentProps<"form">) {
	return (
		<form className={`grid items-start gap-4 ${className}`}>
			<div className="grid gap-2">
				<Label htmlFor="email">Email</Label>
				<Input type="email" id="email" defaultValue="johndoe@example.com" />
			</div>
			<div className="grid gap-2">
				<Label htmlFor="username">Username</Label>
				<Input id="username" defaultValue="@riipandi" />
			</div>
			<Button type="submit">Save changes</Button>
		</form>
	);
}
