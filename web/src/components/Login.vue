<script setup lang="ts">
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogHeader,
} from "@/components/ui/alert-dialog";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ref } from "vue";
import {Tabs, TabsContent, TabsList, TabsTrigger} from "@/components/ui/tabs";

const flag = ref(true);

const toggleFlag = (e:CustomEvent) => {
  const success = e.detail.successful as boolean;

  const status = e.detail.xhr.status as number;

  status == 200 && success ? flag.value = false : flag.value = true;

};

document.addEventListener("htmx:afterRequest", (toggleFlag as EventListener));

const tabsChange = () => {
  //@ts-ignore
  htmx.process(document.querySelector(".description"))
}

</script>

<template>
  <AlertDialog :open="flag">
    <AlertDialogContent>
      <Tabs default-value="login" class="w-full">
        <TabsList class="grid w-full grid-cols-2">
          <TabsTrigger value="login" @click="tabsChange">
            Login
          </TabsTrigger>
          <TabsTrigger value="register" @click="tabsChange">
            Register
          </TabsTrigger>
        </TabsList>
        <TabsContent value="login" class="text-2xl">
          <AlertDialogHeader>
            <AlertDialogDescription class="description">
              <form
                  class="grid grid-cols-2 grid-rows-1 gap-7"
                  hx-post="/auth/login"
                  hx-target="#container-websocket"
                  hx-swap="outerHTML"
              >
                <div class="grid w-full max-w-sm items-center gap-1.5">
                  <Label for="username" class="text-xl">Login</Label>
                  <Input
                      name="username"
                      id="username"
                      type="text"
                      placeholder="dorian"
                      class="text-xl italic"
                  />
                </div>
                <div class="grid w-full max-w-sm items-center gap-1.5">
                  <Label for="password" class="text-xl">Password</Label>
                  <Input
                      name="password"
                      id="password"
                      type="password"
                      value="dorian"
                      placeholder="password"
                      class="text-xl italic"
                  />
                </div>
                <AlertDialogAction type="submit" class="text-2xl">Connection</AlertDialogAction>
              </form>
            </AlertDialogDescription>
          </AlertDialogHeader>
        </TabsContent>
        <TabsContent value="register">
          <AlertDialogHeader>
            <AlertDialogDescription class="description">
              <form
                  class="grid grid-cols-2 grid-rows-1 gap-7"
                  hx-post="/auth/register"
                  hx-target="#container-websocket"
                  hx-swap="outerHTML"
              >
                <div class="grid w-full max-w-sm items-center gap-1.5">
                  <Label for="username" class="text-xl">Register</Label>
                  <Input
                      name="username"
                      id="username"
                      type="text"
                      placeholder="username"
                      class="text-xl italic"
                  />
                </div>
                <div class="grid w-full max-w-sm items-center gap-1.5">
                  <Label for="password" class="text-xl">Password</Label>
                  <Input
                      name="password"
                      id="password"
                      type="password"
                      placeholder="password"
                      class="text-xl italic"
                  />
                </div>
                <AlertDialogAction type="submit" class="text-2xl">Sign in</AlertDialogAction>
              </form>
            </AlertDialogDescription>
          </AlertDialogHeader>
        </TabsContent>
      </Tabs>
      <p id="error-form"></p>
    </AlertDialogContent>
  </AlertDialog>
</template>
