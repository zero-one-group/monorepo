import type { Meta, StoryObj } from "@storybook/react";
import * as Lucide from "lucide-react";
import { getInitials } from "../../utils";
import { Avatar, AvatarFallback, AvatarImage } from "../avatar/avatar";
import { Button } from "../button/button";
import { HoverCard, HoverCardContent, HoverCardTrigger } from "./hover-card";

const meta: Meta<typeof HoverCard> = {
	title: "Basic Components/HoverCard",
	component: HoverCard,
	parameters: {
		layout: "centered",
	},
	argTypes: {},
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: () => (
		<HoverCard>
			<HoverCardTrigger asChild>
				<Button variant="link">@riipandi</Button>
			</HoverCardTrigger>
			<HoverCardContent className="w-80" align="center">
				<div className="flex justify-between space-x-4">
					<Avatar>
						<AvatarImage src="https://avatars.githubusercontent.com/u/921834?v=4" />
						<AvatarFallback>{getInitials("Aris Ripandi")}</AvatarFallback>
					</Avatar>
					<div className="space-y-1">
						<h4 className="font-semibold text-sm">@riipandi</h4>
						<p className="text-sm">
							Polyglot & enthusiastic developer, web artisan, lecturer, and Open
							Source enthusiast.
						</p>
						<div className="flex items-center pt-2">
							<Lucide.Calendar className="mr-2 size-4 opacity-70" />{" "}
							<span className="text-muted-foreground text-xs">
								Joined December 2021
							</span>
						</div>
					</div>
				</div>
			</HoverCardContent>
		</HoverCard>
	),
};

export const AlignStart: Story = {
	render: () => (
		<HoverCard>
			<HoverCardTrigger asChild>
				<Button variant="link">@riipandi</Button>
			</HoverCardTrigger>
			<HoverCardContent className="w-80" align="start">
				<div className="flex justify-between space-x-4">
					<Avatar>
						<AvatarImage src="https://avatars.githubusercontent.com/u/921834?v=4" />
						<AvatarFallback>{getInitials("Aris Ripandi")}</AvatarFallback>
					</Avatar>
					<div className="space-y-1">
						<h4 className="font-semibold text-sm">@riipandi</h4>
						<p className="text-sm">
							Polyglot & enthusiastic developer, web artisan, lecturer, and Open
							Source enthusiast.
						</p>
						<div className="flex items-center pt-2">
							<Lucide.Calendar className="mr-2 size-4 opacity-70" />{" "}
							<span className="text-muted-foreground text-xs">
								Joined December 2021
							</span>
						</div>
					</div>
				</div>
			</HoverCardContent>
		</HoverCard>
	),
};

export const AlignEnd: Story = {
	render: () => (
		<HoverCard>
			<HoverCardTrigger asChild>
				<Button variant="link">@riipandi</Button>
			</HoverCardTrigger>
			<HoverCardContent className="w-80" align="end">
				<div className="flex justify-between space-x-4">
					<Avatar>
						<AvatarImage src="https://avatars.githubusercontent.com/u/921834?v=4" />
						<AvatarFallback>{getInitials("Aris Ripandi")}</AvatarFallback>
					</Avatar>
					<div className="space-y-1">
						<h4 className="font-semibold text-sm">@riipandi</h4>
						<p className="text-sm">
							Polyglot & enthusiastic developer, web artisan, lecturer, and Open
							Source enthusiast.
						</p>
						<div className="flex items-center pt-2">
							<Lucide.Calendar className="mr-2 size-4 opacity-70" />{" "}
							<span className="text-muted-foreground text-xs">
								Joined December 2021
							</span>
						</div>
					</div>
				</div>
			</HoverCardContent>
		</HoverCard>
	),
};
