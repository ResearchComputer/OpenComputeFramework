'use client'
import { Metadata } from "next"
import Image from "next/image"
import { Activity, CreditCard, DollarSign, Users } from "lucide-react"
import {useState, useEffect} from 'react'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/registry/new-york/ui/card"
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "@/registry/new-york/ui/tabs"
import { ServiceOverview } from "@/components/service-overview"
import { NodeOverview } from "@/components/node-overview"
import { public_relay } from "@/lib/api"

async function getData() {
  const res = await fetch(public_relay+'/api/v1/status/table')
  if (!res.ok) {
    throw new Error('Failed to fetch data')
  }
  return res.json()
}

function countServices(data:any) {
  let services:any = []
  for (let node of data.nodes) {
    if (!services.includes(node.service)) {
      services.push(node.service)
    }
  }
  return services.length
}

export default function DashboardPage() {
  const [data, setData] = useState({nodes:[]})
  const [isLoading, setLoading] = useState(true)

  useEffect(() => {
    fetch(public_relay+'/api/v1/status/table')
      .then((res) => {
        return res.json()
      })
      .then((data:any) => {
        setData(data)
        setLoading(false)
      }).catch((err) => {
      })
  }, [])
  if (isLoading) return <p>Loading...</p>

  return (
    <>
      <div className="md:hidden">
        <Image
          src="/examples/dashboard-light.png"
          width={1280}
          height={866}
          alt="Dashboard"
          className="block dark:hidden"
        />
        <Image
          src="/examples/dashboard-dark.png"
          width={1280}
          height={866}
          alt="Dashboard"
          className="hidden dark:block"
        />
      </div>
      <div className="hidden full flex-col md:flex">
        <div className="container flex-1 space-y-4 p-8 pt-6">
          <div className="flex items-center justify-between space-y-2">
            <h2 className="text-lg font-semibold">Dashboard</h2>
          </div>
          
          <Tabs defaultValue="overview" className="space-y-4">
            <TabsList>
              <TabsTrigger value="overview">Overview</TabsTrigger>
              <TabsTrigger value="analytics" disabled>
                Analytics
              </TabsTrigger>
              <TabsTrigger value="reports" disabled>
                Reports
              </TabsTrigger>
              <TabsTrigger value="notifications" disabled>
                Notifications
              </TabsTrigger>
            </TabsList>
            <TabsContent value="overview" className="space-y-4">
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">
                      Connected Nodes
                    </CardTitle>
                    <DollarSign className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{data.nodes.length}</div>
                    <p className="text-xs text-muted-foreground">
                      ??? from last month
                    </p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">
                      Provided Service
                    </CardTitle>
                    <Users className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{countServices(data)}</div>
                    <p className="text-xs text-muted-foreground">
                      ??? from last month
                    </p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">Served Requests</CardTitle>
                    <CreditCard className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">?</div>
                    <p className="text-xs text-muted-foreground">
                      ? from last month
                    </p>
                  </CardContent>
                </Card>
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">
                      Authorized Users
                    </CardTitle>
                    <Activity className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">?</div>
                    <p className="text-xs text-muted-foreground">
                      ? since last hour
                    </p>
                  </CardContent>
                </Card>
              </div>
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
                <Card className="col-span-4">
                  <CardHeader>
                    <CardTitle>Services</CardTitle>
                  </CardHeader>
                  <CardContent className="pl-2">
                    <ServiceOverview nodes={data.nodes} />
                  </CardContent>
                </Card>
                <Card className="col-span-3">
                  <CardHeader>
                    <CardTitle>Connected Nodes</CardTitle>
                    <CardDescription>
                      Total {data.nodes.length} nodes connected
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <NodeOverview nodes={data.nodes} />
                  </CardContent>
                </Card>
              </div>
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </>
  )
}
