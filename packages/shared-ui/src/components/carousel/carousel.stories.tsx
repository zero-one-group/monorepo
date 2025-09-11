import type { Meta, StoryObj } from "@storybook/react";
import { Card, CardContent } from "../card/card";
import {
	Carousel,
	CarouselContent,
	CarouselItem,
	CarouselNext,
	CarouselPrevious,
} from "./carousel";

const meta: Meta<typeof Carousel> = {
	title: "Basic Components/Carousel",
	component: Carousel,
};

export default meta;
type Story = StoryObj<typeof meta>;

const slides = [
	{ id: "slide-1", number: 1 },
	{ id: "slide-2", number: 2 },
	{ id: "slide-3", number: 3 },
	{ id: "slide-4", number: 4 },
	{ id: "slide-5", number: 5 },
	{ id: "slide-5", number: 5 },
];

export const Default: Story = {
	render: () => (
		<Carousel className="w-full max-w-xs">
			<CarouselContent>
				{slides.map((slide) => (
					<CarouselItem key={slide.id}>
						<Card>
							<CardContent className="flex aspect-square items-center justify-center p-6">
								<span className="font-semibold text-4xl">{slide.number}</span>
							</CardContent>
						</Card>
					</CarouselItem>
				))}
			</CarouselContent>
			<CarouselPrevious />
			<CarouselNext />
		</Carousel>
	),
};

export const Vertical: Story = {
	render: () => (
		<Carousel
			opts={{ align: "start" }}
			orientation="vertical"
			className="w-full max-w-xs"
		>
			<CarouselContent className="-mt-1 h-[200px]">
				{slides.map((slide) => (
					<CarouselItem key={slide.id} className="pt-1 md:basis-1/2">
						<div className="p-1">
							<Card>
								<CardContent className="flex items-center justify-center p-6">
									<span className="font-semibold text-3xl">{slide.number}</span>
								</CardContent>
							</Card>
						</div>
					</CarouselItem>
				))}
			</CarouselContent>
			<CarouselPrevious />
			<CarouselNext />
		</Carousel>
	),
};

export const Spacing: Story = {
	render: () => (
		<Carousel className="w-full max-w-sm">
			<CarouselContent className="-ml-1">
				{slides.map((slide) => (
					<CarouselItem
						key={slide.id}
						className="pl-1 md:basis-1/2 lg:basis-1/3"
					>
						<div className="p-1">
							<Card>
								<CardContent className="flex aspect-square items-center justify-center p-6">
									<span className="font-semibold text-2xl">{slide.number}</span>
								</CardContent>
							</Card>
						</div>
					</CarouselItem>
				))}
			</CarouselContent>
			<CarouselPrevious />
			<CarouselNext />
		</Carousel>
	),
};
