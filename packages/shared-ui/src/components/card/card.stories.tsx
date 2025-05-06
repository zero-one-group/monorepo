import { zodResolver } from '@hookform/resolvers/zod'
import type { Meta, StoryObj } from '@storybook/react'
import { consola } from 'consola'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { Button } from '../button/button'
import { Form, FormControl, FormLabel, FormMessage } from '../form/form'
import { FormDescription, FormField, FormItem } from '../form/form'
import { Input } from '../input/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../select/select'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from './card'

const FormSchema = z.object({
  projectName: z.string({
    required_error: 'Project name is required.',
  }),
  framework: z.string({
    required_error: 'Please select a framework.',
  }),
})

const meta: Meta<typeof Card> = {
  title: 'Basic Components/Card',
  component: Card,
}

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  render: () => {
    const form = useForm<z.infer<typeof FormSchema>>({
      resolver: zodResolver(FormSchema),
    })

    function onSubmit(data: z.infer<typeof FormSchema>) {
      consola.debug(data)
    }

    return (
      <Card className="w-sm">
        <CardHeader>
          <CardTitle>Create project</CardTitle>
          <CardDescription>Deploy your new project in one-click.</CardDescription>
        </CardHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <CardContent className="space-y-4">
              <FormField
                control={form.control}
                name="projectName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Project Name</FormLabel>
                    <FormControl>
                      <Input placeholder="Enter project name" {...field} />
                    </FormControl>
                    <FormDescription>This is your project's display name.</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="framework"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Framework</FormLabel>
                    <Select onValueChange={field.onChange} defaultValue={field.value}>
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Select a framework" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value="next">Next.js</SelectItem>
                        <SelectItem value="sveltekit">SvelteKit</SelectItem>
                        <SelectItem value="astro">Astro</SelectItem>
                        <SelectItem value="nuxt">Nuxt.js</SelectItem>
                      </SelectContent>
                    </Select>
                    <FormDescription>Select your preferred framework.</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
            <CardFooter className="flex justify-between">
              <Button type="button" variant="outline">
                Cancel
              </Button>
              <Button type="submit">Deploy</Button>
            </CardFooter>
          </form>
        </Form>
      </Card>
    )
  },
}
