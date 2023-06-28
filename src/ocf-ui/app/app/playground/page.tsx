'use client'

import { Redo2 } from "lucide-react"
import * as React from "react"
import { Button } from "@/registry/new-york/ui/button"
import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "@/registry/new-york/ui/hover-card"
import { Label } from "@/registry/new-york/ui/label"
import { Separator } from "@/registry/new-york/ui/separator"
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "@/registry/new-york/ui/tabs"
import { Textarea } from "@/registry/new-york/ui/textarea"
import { ElementType } from "react"
import { CodeViewer } from "@/components/code-viewer"
import { Icons } from "@/components/icons"
import { MaxLengthSelector } from "@/components/maxlength-selector"
import { ModelSelector } from "@/components/model-selector"
import { PresetActions } from "@/components/preset-actions"
import { PresetSave } from "@/components/preset-save"
import { PresetSelector } from "@/components/preset-selector"
import { PresetShare } from "@/components/preset-share"
import { presets } from "./data/presets"
import { TemperatureSelector } from "@/components/temperature-selector"
import { TopPSelector } from "@/components/top-p-selector"
import { TopKSelector } from "@/components/top-k-selector"
import { types } from "./data/models"
import "./styles.css"
import Image from "next/image"
import { public_relay } from "@/lib/api"

async function getData() {
  const res = await fetch(public_relay+'/api/v1/status/table')
  if (!res.ok) {
    throw new Error('Failed to fetch data')
  }
  return res.json()
}

function getModels(data:any) {
  let services:any = []
  for (let node of data.nodes) {
    if (!services.includes(node.service) && node.service.startsWith("inference:")) {
      services.push(node.service.replace("inference:", ""))
    }
  }
  services = services.map((service:any) => {
    return {
      id: service,
      name: service,
      description: service,
      type: 'Text-Completion',
      strengths: "N/A"
    }
  })
  return services
}

export default async function PlaygroundPage() {

  const data = await getData()
  const models = getModels(data)
  
  const completion_ref = React.createRef<HTMLTextAreaElement>();
  const temperature_ref = React.createRef<HTMLSpanElement>();

  async function submit() {
    let prompt = completion_ref.current?.value
    let temperature = temperature_ref.current?.innerText
    let resp = await fetch(public_relay+'/api/v1/request/inference', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        model_name:'togethercomputer/RedPajama-INCITE-Chat-3B-v1',
        params: {
          prompt: prompt,
          temperature: temperature
        }
      })
    })
    let response_json:any = await resp.json()
    let reply = JSON.parse(response_json['data'])['output']['text']
    completion_ref.current!.value += reply
  }
  return (
    <>
      <div className="md:hidden">
        <Image
          src="/examples/playground-light.png"
          width={1280}
          height={916}
          alt="Playground"
          className="block dark:hidden"
        />
        <Image
          src="/examples/playground-dark.png"
          width={1280}
          height={916}
          alt="Playground"
          className="hidden dark:block"
        />
      </div>
      <div className="hidden h-full flex-col md:flex">
        <div className="container flex flex-col items-start justify-between space-y-2 py-4 sm:flex-row sm:items-center sm:space-y-0 md:h-16">
          <h2 className="text-lg font-semibold">Playground</h2>
          <div className="ml-auto flex w-full space-x-2 sm:justify-end">
            {/* <PresetSelector presets={presets} />
            <PresetSave /> */}
            <div className="hidden space-x-2 md:flex">
              <CodeViewer />
              {/* <PresetShare /> */}
            </div>
            {/* <PresetActions /> */}
          </div>
        </div>
        <Separator />
        <Tabs defaultValue="complete" className="flex-1">
          <div className="container h-full py-6">
            <div className="grid h-full items-stretch gap-6 md:grid-cols-[1fr_200px]">
              <div className="hidden flex-col space-y-4 sm:flex md:order-2">
                <div className="grid gap-2">
                  <HoverCard openDelay={200}>
                    <HoverCardTrigger asChild>
                      <span className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                        Mode
                      </span>
                    </HoverCardTrigger>
                    <HoverCardContent className="w-[320px] text-sm" side="left">
                      Choose the interface that best suits your task. You can
                      provide: a simple prompt to complete, starting and ending
                      text to insert a completion within, or some text with
                      instructions to edit it.
                    </HoverCardContent>
                  </HoverCard>
                  <TabsList className="grid grid-cols-3">
                    <TabsTrigger value="complete">
                      <span className="sr-only">Complete</span>
                      <Icons.completeMode className="h-5 w-5" />
                    </TabsTrigger>
                    <TabsTrigger value="insert">
                      <span className="sr-only">Insert</span>
                      <Icons.insertMode className="h-5 w-5" />
                    </TabsTrigger>
                    <TabsTrigger value="edit">
                      <span className="sr-only">Edit</span>
                      <Icons.editMode className="h-5 w-5" />
                    </TabsTrigger>
                  </TabsList>
                </div>
                <ModelSelector types={types} models={models} />
                <TemperatureSelector defaultValue={[0.6]} ref={temperature_ref}/>
                <MaxLengthSelector defaultValue={[256]} />
                <TopPSelector defaultValue={[0.9]} />
                <TopKSelector defaultValue={[50]} />
              </div>
              <div className="md:order-1">
                <TabsContent value="complete" className="mt-0 border-0 p-0">
                  <div className="flex h-full flex-col space-y-4">
                    <Textarea
                      placeholder="Write a tagline for an ice cream shop"
                      className="min-h-[400px] flex-1 p-4 md:min-h-[700px] lg:min-h-[700px]" ref={completion_ref}
                    />
                    <div className="flex items-center space-x-2">
                      <Button onClick={submit}>Submit</Button>
                      <Button variant="secondary">
                        <span className="sr-only">Show history</span>
                        <Redo2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </TabsContent>
                <TabsContent value="insert" className="mt-0 border-0 p-0">
                  <div className="flex flex-col space-y-4">
                    <div className="grid h-full grid-rows-2 gap-6 lg:grid-cols-2 lg:grid-rows-1">
                      <Textarea
                        placeholder="We're writing to [inset]. Congrats!"
                        className="h-full min-h-[300px] lg:min-h-[700px] xl:min-h-[700px]"
                      />
                      <div className="rounded-md border bg-muted"></div>
                    </div>
                    <div className="flex items-center space-x-2">
                      <Button onClick={submit}>Submit</Button>
                      <Button variant="secondary">
                        <span className="sr-only">Clear</span>
                        <Redo2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </TabsContent>
                <TabsContent value="edit" className="mt-0 border-0 p-0">
                  <div className="flex flex-col space-y-4">
                    <div className="grid h-full gap-6 lg:grid-cols-2">
                      <div className="flex flex-col space-y-4">
                        <div className="flex flex-1 flex-col space-y-2">
                          <Label htmlFor="input">Input</Label>
                          <Textarea
                            id="input"
                            placeholder="We is going to the market."
                            className="flex-1 lg:min-h-[580px]"
                          />
                        </div>
                        <div className="flex flex-col space-y-2">
                          <Label htmlFor="instructions">Instructions</Label>
                          <Textarea
                            id="instructions"
                            placeholder="Fix the grammar."
                          />
                        </div>
                      </div>
                      <div className="mt-[21px] min-h-[400px] rounded-md border bg-muted lg:min-h-[700px]" />
                    </div>
                    <div className="flex items-center space-x-2">
                      <Button onClick={() => {submit}}>Submit</Button>
                      
                      <Button variant="secondary">
                        <span className="sr-only">Clear</span>
                        <Redo2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </TabsContent>
              </div>
            </div>
          </div>
        </Tabs>
      </div>
    </>
  )
}